package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func main() {
	log.Println("Starting bot")
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
	db, err := sql.Open("mysql", dataSourceName)
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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	log.Println("Stopping bot")
}
