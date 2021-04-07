package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/commands"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func main() {
	log.Println("Starting bot")
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
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
		log.Fatalln("Error connecting to database", err)
	}
	defer db.Close()

	err = MigrateUp("mysql://" + dataSourceName + "?multiStatements=true")
	if err != nil {
		log.Fatalln("Error running database migrations", err)
	}

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalln("Error creating discord connection", err)
	}
	defer discord.Close()

	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Connect) {
		log.Println("Connected to discord")
	})
	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Disconnect) {
		log.Println("Disconnected from discord")
	})

	repo := repository.NewDBRepository(db)
	brm := boost_request.NewBoostRequestManager(discord, repo)

	defer brm.Destroy()

	err = discord.Open()
	if err != nil {
		log.Fatalln("Error connecting to discord", err)
	}

	registerCommands(discord, repo)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	log.Println("Stopping bot")
}

func registerCommands(discord *discordgo.Session, repo repository.Repository) {
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
	router.RegisterMiddleware(func(next dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			// TODO check permissions
			next(ctx)
		}
	})

	router.Initialize(discord)
}
