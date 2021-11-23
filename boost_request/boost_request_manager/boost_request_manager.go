package boost_request_manager

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/api/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/active_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/boost_request/sequences"
	"github.com/oppzippy/BoostRequestBot/boost_request/webhook"
)

type BoostRequestManager struct {
	discord        *discordgo.Session
	repo           repository.Repository
	bundle         *i18n.Bundle
	messenger      *messenger.BoostRequestMessenger
	webhookManager *webhook.WebhookManager
	activeRequests *sync.Map
	isLoaded       bool
	isLoadedLock   *sync.Mutex
}

func NewBoostRequestManager(discord *discordgo.Session, repo repository.Repository, bundle *i18n.Bundle) *BoostRequestManager {
	brm := &BoostRequestManager{
		discord:        discord,
		repo:           repo,
		bundle:         bundle,
		messenger:      messenger.NewBoostRequestMessenger(discord, bundle),
		webhookManager: webhook.NewWebhookManager(repo),
		activeRequests: new(sync.Map),
		isLoadedLock:   new(sync.Mutex),
	}

	return brm
}

func (brm *BoostRequestManager) Destroy() {
	brm.webhookManager.Destroy()
	brm.messenger.Destroy()
}

func (brm *BoostRequestManager) LoadBoostRequests() {
	brm.isLoadedLock.Lock()
	defer brm.isLoadedLock.Unlock()
	if !brm.isLoaded {
		boostRequests, err := brm.repo.GetUnresolvedBoostRequests()
		if err != nil {
			log.Printf("Unable to load boost requests: %v", err)
			return
		}
		brm.isLoaded = true
		for _, br := range boostRequests {
			brm.activeRequests.Store(br.BackendMessageID, active_request.NewActiveRequest(*br, brm.setWinner))
		}
	}
}

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

func (brm *BoostRequestManager) CreateBoostRequest(
	brc *repository.BoostRequestChannel, brPartial BoostRequestPartial,
) (*repository.BoostRequest, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	br := &repository.BoostRequest{
		ExternalID:             &id,
		Channel:                *brc,
		RequesterID:            brPartial.RequesterID,
		Message:                brPartial.Message,
		Price:                  brPartial.Price,
		AdvertiserCut:          brPartial.AdvertiserCut,
		AdvertiserRoleCuts:     brPartial.AdvertiserRoleCuts,
		Discount:               brPartial.Discount,
		EmbedFields:            brPartial.EmbedFields,
		PreferredAdvertiserIDs: brPartial.PreferredAdvertiserIDs,
		CreatedAt:              time.Now().UTC(),
	}

	if brc.UsesBuyerMessage {
		br.BackendMessageID = brPartial.BackendMessageID
	}

	// Essentially checking if the user is a bot
	// TODO add IsBot to BoostRequest
	if br.EmbedFields == nil && br.Price == 0 {
		roleDiscounts, err := brm.getRoleDiscountsForUser(br.Channel.GuildID, br.RequesterID)
		if err != nil {
			// They won't get their discounts, but we don't have to abort
			log.Printf("Error searching roles for discounts: %v", err)
		}
		br.RoleDiscounts = roleDiscounts
	}

	sequenceArgs := sequences.CreateSequenceArgs{
		Repository:        brm.repo,
		BoostRequest:      br,
		Discord:           brm.discord,
		Messenger:         brm.messenger,
		ActiveRequests:    brm.activeRequests,
		SetWinnerCallback: brm.setWinner,
	}
	if br.EmbedFields == nil {
		err := sequences.RunCreateHumanRequesterSequence(sequenceArgs)
		if err != nil {
			return nil, fmt.Errorf("boost request creation failed with human requester: %v", err)
		}
	} else {
		err := sequences.RunCreateBotRequesterSequence(sequenceArgs)
		if err != nil {
			return nil, fmt.Errorf("boost request creation failed with bot requester: %v", err)
		}
	}
	return br, nil
}

// Best is defined as the role with the highest weight
func (brm *BoostRequestManager) GetBestRolePrivileges(guildID string, roles []string) *repository.AdvertiserPrivileges {
	var bestPrivileges *repository.AdvertiserPrivileges = nil
	for _, role := range roles {
		privileges, err := brm.repo.GetAdvertiserPrivilegesForRole(guildID, role)
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching privileges: %v", err)
		}
		if privileges != nil {
			if bestPrivileges == nil || privileges.Weight > bestPrivileges.Weight {
				bestPrivileges = privileges
			}
		}
	}

	return bestPrivileges
}

func (brm *BoostRequestManager) AddAdvertiserToBoostRequest(br *repository.BoostRequest, userID string) (err error) {
	// TODO cache roles
	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		return fmt.Errorf("fetching guild member: %v", err)
	}
	privileges := brm.GetBestRolePrivileges(br.Channel.GuildID, guildMember.Roles)
	if privileges == nil {
		return ErrNoPrivileges
	}
	if len(br.PreferredAdvertiserIDs) == 0 {
		brm.signUp(br, userID, privileges)
	} else {
		var isPreferredAdvertiser bool
		for id := range br.PreferredAdvertiserIDs {
			if id == userID {
				isPreferredAdvertiser = true
				break
			}
		}
		if !isPreferredAdvertiser {
			return ErrNotPreferredAdvertiser
		}

		req, ok := brm.activeRequests.Load(br.BackendMessageID)
		if !ok {
			log.Printf("AddAdvertiserToBoostRequest: activeRequest not found")
			return
		}
		r := req.(*active_request.ActiveRequest)
		r.SetAdvertiser(userID)
	}
	return nil
}

func (brm *BoostRequestManager) RemoveAdvertiserFromBoostRequest(backendMessageID string, userID string) (removed bool) {
	req, ok := brm.activeRequests.Load(backendMessageID)
	if ok {
		ar, ok := req.(*active_request.ActiveRequest)
		if ok && ar.HasSignup(userID) {
			ar.RemoveSignup(userID)
			return true
		}
	}
	return false
}

func (brm *BoostRequestManager) IsAdvertiserSignedUpForBoostRequest(backendMessageID, userID string) bool {
	req, ok := brm.activeRequests.Load(backendMessageID)
	if ok {
		ar, ok := req.(*active_request.ActiveRequest)
		if ok {
			return ar.HasSignup(userID)
		}
	}
	return false
}

func (brm *BoostRequestManager) StealBoostRequest(br *repository.BoostRequest, userID string) (ok, usedCredits bool) {
	if len(br.PreferredAdvertiserIDs) > 0 {
		var isPreferredAdvertiser bool
		for id := range br.PreferredAdvertiserIDs {
			if id == userID {
				isPreferredAdvertiser = true
				break
			}
		}
		if !isPreferredAdvertiser {
			return false, false
		}
	}

	credits, err := brm.repo.GetStealCreditsForUser(br.Channel.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits: %v", err)
		return false, false
	}
	if credits <= 0 {
		return false, false
	}

	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching guild member: %v", err)
		return false, false
	}
	privileges := brm.GetBestRolePrivileges(br.Channel.GuildID, guildMember.Roles)
	if privileges == nil {
		return false, false
	}

	req, ok := brm.activeRequests.Load(br.BackendMessageID)
	if !ok {
		return false, false
	}
	r := req.(*active_request.ActiveRequest)
	endTime := br.CreatedAt.Add(time.Duration(privileges.Delay) * time.Second)
	now := time.Now()
	ok = r.SetAdvertiser(userID)
	// Don't subtract steal credits if the 👍 button would have had the same effect
	if ok && now.Before(endTime) {
		err := brm.repo.AdjustStealCreditsForUser(br.Channel.GuildID, userID, repository.OperationSubtract, 1)
		if err != nil {
			log.Printf("Error subtracting boost request credits after use: %v", err)
			return false, false
		}
		usedCredits = true
	}
	return ok, usedCredits
}

func (brm *BoostRequestManager) setWinner(event *active_request.AdvertiserChosenEvent) {
	br := event.BoostRequest
	userID := event.UserID

	brm.activeRequests.Delete(br.BackendMessageID)
	br.AdvertiserID = userID
	br.IsResolved = true
	br.ResolvedAt = time.Now()
	err := brm.repo.ResolveBoostRequest(&br)
	if err != nil {
		// Log the error but try to keep things running.
		// There will be data loss, but that is better than a lost sale.
		log.Printf("Error resolving boost request: %v", err)
	}

	_, err = brm.messenger.SendBackendAdvertiserChosenMessage(&br)
	if err != nil {
		log.Printf("Error sending message to boost request backend: %v", err)
	}

	_, err = brm.messenger.SendAdvertiserChosenDMToAdvertiser(&br)
	if err != nil {
		log.Printf("Error sending advertsier chosen DM to advertiser: %v", err)
	}

	if !br.Channel.SkipsBuyerDM {
		_, err = brm.messenger.SendAdvertiserChosenDMToRequester(&br)
		if err != nil {
			log.Printf("Error sending advertiser chosen DM to requester: %v", err)
		}
	}

	rollChannel, err := brm.repo.GetRollChannel(br.Channel.GuildID)
	if err != nil {
		if err != repository.ErrNoResults {
			log.Printf("Error fetching log channel: %v", err)
		}
	} else if event.RollResults != nil {
		_, err := brm.messenger.SendRoll(rollChannel, &br, event.RollResults)
		if err != nil {
			log.Printf("Error sending roll message: %v", err)
		}
	}

	preferredAdvertiserIDs := make([]string, 0, len(br.PreferredAdvertiserIDs))
	for id := range br.PreferredAdvertiserIDs {
		preferredAdvertiserIDs = append(preferredAdvertiserIDs, id)
	}

	err = brm.webhookManager.QueueToSend(br.Channel.GuildID, &webhook.WebhookEvent{
		Event: webhook.AdvertiserChosenEvent,
		Payload: models.BoostRequest{
			ID:                     br.ExternalID.String(),
			RequesterID:            br.RequesterID,
			IsAdvertiserSelected:   br.IsResolved,
			AdvertiserID:           br.AdvertiserID,
			BackendChannelID:       br.Channel.BackendChannelID,
			BackendMessageID:       br.BackendMessageID,
			Message:                br.Message,
			Price:                  br.Price,
			AdvertiserCut:          br.AdvertiserCut,
			PreferredAdvertiserIDs: preferredAdvertiserIDs,
			CreatedAt:              br.CreatedAt.Format(time.RFC3339),
			AdvertiserSelectedAt:   br.ResolvedAt.Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Printf("error queueing webhook: %v", err)
	}
}

func (brm *BoostRequestManager) getRoleDiscountsForUser(guildID, userID string) ([]*repository.RoleDiscount, error) {
	member, err := brm.discord.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}
	bestDiscounts, err := brm.repo.GetBestDiscountsForRoles(guildID, member.Roles)
	return bestDiscounts, err
}

func (brm *BoostRequestManager) signUp(br *repository.BoostRequest, userID string, privileges *repository.AdvertiserPrivileges) {
	req, ok := brm.activeRequests.Load(br.BackendMessageID)
	if !ok {
		return
	}
	r := req.(*active_request.ActiveRequest)
	r.AddSignup(userID, *privileges)
}

func (brm *BoostRequestManager) CancelBoostRequest(br *repository.BoostRequest) error {
	err := brm.discord.ChannelMessageDelete(br.Channel.BackendChannelID, br.BackendMessageID)
	if err != nil {
		log.Printf("failed to delete backend message when cancelling boost request: %v", err)
		// it's not an ideal situation but we can still continue
	}
	err = brm.repo.DeleteBoostRequest(br)
	if err != nil {
		return err
	}
	activeRequestInterface, loaded := brm.activeRequests.LoadAndDelete(br.BackendMessageID)
	if !loaded {
		return nil
	}
	activeRequest := activeRequestInterface.(*active_request.ActiveRequest)
	activeRequest.Destroy()
	return nil
}
