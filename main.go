package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/api"
	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/commands"
	"github.com/oppzippy/BoostRequestBot/boost_request/middleware"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	db_repository "github.com/oppzippy/BoostRequestBot/boost_request/repository/database"
)

func main() {
	log.Println("Starting bot")
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	db, err := sql.Open("mysql", dataSourceName+"?parseTime=true")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	err = MigrateUp("mysql://" + dataSourceName + "?multiStatements=true")
	if err != nil {
		log.Fatalf("Error running database migrations: %v", err)
	}

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating discord connection: %v", err)
	}
	defer discord.Close()
	discord.Identify.Intents = discordgo.IntentsNone

	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Connect) {
		log.Println("Connected to discord")
	})
	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Disconnect) {
		log.Println("Disconnected from discord")
	})

	var repo repository.Repository = db_repository.NewRepository(db)
	brm := boost_request.NewBoostRequestManager(discord, repo)
	brm.LoadBoostRequests()
	registerCommandRouter(discord, repo)

	defer brm.Destroy()

	server := api.NewWebAPI(repo)

	err = discord.Open()
	if err != nil {
		log.Fatalf("Error connecting to discord: %v", err)
	}

	sc := make(chan os.Signal, 1)

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Error starting http server: %v", err)
		} else {
			sc <- syscall.SIGINT
		}
	}()
	defer server.Shutdown(context.TODO())

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	log.Println("Stopping bot")
}

func registerCommandRouter(discord *discordgo.Session, repo repository.Repository) {
	router := dgc.Create(&dgc.Router{
		Prefixes: []string{
			"!",
		},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands: []*dgc.Command{
			&commands.MainCommand,
		},
		Middlewares: []dgc.Middleware{},
	})
	router.RegisterMiddleware(func(next dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			ctx.CustomObjects.Set("repo", repo)
			next(ctx)
		}
	})
	adminOnlyMiddleware := middleware.AdminOnlyMiddleware{}
	router.RegisterMiddleware(adminOnlyMiddleware.Exec)
	router.RegisterMiddleware(middleware.GuildOnlyMiddleware)

	router.Initialize(discord)
}
