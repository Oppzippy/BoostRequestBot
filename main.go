package main

import (
	"context"
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
	"github.com/oppzippy/BoostRequestBot/initialization"
	"github.com/oppzippy/BoostRequestBot/locales"
)

func main() {
	fmt.Println("Starting bot")
	rand.Seed(time.Now().UnixNano())
	localeBundle := locales.Bundle()

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	db, err := initialization.GetDBC()
	if err != nil {
		log.Fatalf("Failed to acquire database connection: %v", err)
	}

	discord, err := setUpDiscord()
	if err != nil {
		log.Fatalf("Error setting up discord: %v", err)
	}
	defer discord.Close()

	var repo repository.Repository = db_repository.NewRepository(db)
	brm := boost_request.NewBoostRequestManager(discord, repo, localeBundle)
	defer brm.Destroy()
	brm.LoadBoostRequests()

	_ = boost_request.NewBoostRequestDiscordHandler(discord, repo, brm)
	registerCommandRouter(discord, repo)

	server := api.NewWebAPI(repo, brm, os.Getenv("HTTP_LISTEN_ADDRESS"))

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
	fmt.Println("Stopping bot")
}

func setUpDiscord() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("creating discord connection: %v", err)
	}
	discord.Identify.Intents = discordgo.IntentsNone

	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Connect) {
		fmt.Println("Connected to discord")
	})
	discord.AddHandler(func(_ *discordgo.Session, event *discordgo.Disconnect) {
		fmt.Println("Disconnected from discord")
	})
	return discord, nil
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
