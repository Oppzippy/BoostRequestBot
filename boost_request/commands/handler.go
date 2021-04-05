package commands

import (
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type CommandHandler struct {
	repo repository.Repository
}

func NewCommandHandler(repo repository.Repository) CommandHandler {
	return CommandHandler{
		repo: repo,
	}
}
