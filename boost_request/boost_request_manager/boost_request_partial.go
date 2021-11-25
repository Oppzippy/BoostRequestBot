package boost_request_manager

import (
	"strconv"

	"github.com/oppzippy/BoostRequestBot/api/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestPartial struct {
	RequesterID            string
	Message                string
	EmbedFields            []*repository.MessageEmbedField
	PreferredAdvertiserIDs map[string]struct{}
	BackendMessageID       string
	Price                  int64
	AdvertiserCut          int64
	AdvertiserRoleCuts     map[string]int64
	Discount               int64
}

func FromModelBoostRequestPartial(br *models.BoostRequestPartial) (*BoostRequestPartial, error) {
	roleCuts := make(map[string]int64)
	for roleID, cutStr := range br.AdvertiserRoleCuts {
		cut, err := strconv.ParseInt(cutStr, 10, 64)
		if err != nil {
			return nil, err
		}
		roleCuts[roleID] = cut
	}

	preferredAdvertiserIDs := make(map[string]struct{})
	for _, advertiserID := range br.PreferredAdvertiserIDs {
		preferredAdvertiserIDs[advertiserID] = struct{}{}
	}

	return &BoostRequestPartial{
		RequesterID:            br.RequesterID,
		Message:                br.Message,
		PreferredAdvertiserIDs: preferredAdvertiserIDs,
		Price:                  br.Price,
		AdvertiserCut:          br.AdvertiserCut,
		AdvertiserRoleCuts:     roleCuts,
		Discount:               br.Discount,
	}, nil
}
