package boost_request_manager

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	models_v1 "github.com/oppzippy/BoostRequestBot/api/v1/models"
	"github.com/oppzippy/BoostRequestBot/api/v2/models"
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

func NewBoostRequestManager(
	discord *discordgo.Session,
	repo repository.Repository,
	bundle *i18n.Bundle,
	messenger *messenger.BoostRequestMessenger,
) *BoostRequestManager {
	brm := &BoostRequestManager{
		discord:        discord,
		repo:           repo,
		bundle:         bundle,
		messenger:      messenger,
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
			brm.activeRequests.Store(br.ID, active_request.NewActiveRequest(*br, brm.setWinner))
		}
	}
}

func (brm *BoostRequestManager) CreateBoostRequest(
	brc *repository.BoostRequestChannel,
	brPartial *BoostRequestPartial,
) (*repository.BoostRequest, error) {
	br, err := brm.partialToBoostRequest(brc, brPartial)
	if err != nil {
		return nil, err
	}

	err = brm.dispatchBoostRequest(br)
	if err != nil {
		return nil, err
	}

	go brm.applyAutoSignups(br)

	return br, nil
}

func (brm *BoostRequestManager) partialToBoostRequest(brc *repository.BoostRequestChannel, brPartial *BoostRequestPartial) (*repository.BoostRequest, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	br := &repository.BoostRequest{
		GuildID:                brPartial.GuildID,
		BackendChannelID:       brPartial.BackendChannelID,
		ExternalID:             &id,
		Channel:                brc,
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

	if brc != nil && brc.UsesBuyerMessage {
		br.BackendMessages = []*repository.BoostRequestBackendMessage{
			{
				ChannelID: brc.FrontendChannelID,
				MessageID: brPartial.BackendMessageID,
			},
		}
	}

	// Essentially checking if the user is a bot
	// TODO add IsBot to BoostRequest
	if br.EmbedFields == nil && br.Price == 0 {
		roleDiscounts, err := brm.getRoleDiscountsForUser(br.GuildID, br.RequesterID)
		if err != nil {
			// They won't get their discounts, but we don't have to abort
			log.Printf("Error searching roles for discounts: %v", err)
		}
		br.RoleDiscounts = roleDiscounts
	}
	return br, nil
}

func (brm *BoostRequestManager) dispatchBoostRequest(br *repository.BoostRequest) error {
	sequenceArgs := sequences.CreateSequenceArgs{
		Repository:               brm.repo,
		BoostRequest:             br,
		Discord:                  brm.discord,
		Messenger:                brm.messenger,
		ActiveRequests:           brm.activeRequests,
		SetWinnerCallback:        brm.setWinner,
		BackendMessageChannelIDs: make(map[string]struct{}),
	}
	if len(sequenceArgs.BackendMessageChannelIDs) == 0 {
		if br.Channel != nil {
			sequenceArgs.BackendMessageChannelIDs[br.Channel.BackendChannelID] = struct{}{}
		} else if len(br.PreferredAdvertiserIDs) == 1 {
			// Only one preferred advertiser is set so just dm them the request
			for userID := range br.PreferredAdvertiserIDs {
				channel, err := brm.discord.UserChannelCreate(userID)
				if err != nil {
					return err
				}
				sequenceArgs.BackendMessageChannelIDs[channel.ID] = struct{}{}
			}
		} else {
			sequenceArgs.BackendMessageChannelIDs[br.BackendChannelID] = struct{}{}
		}
	}
	if br.EmbedFields == nil {
		err := sequences.RunCreateHumanRequesterSequence(sequenceArgs)
		if err != nil {
			return fmt.Errorf("boost request creation failed with human requester: %v", err)
		}
	} else {
		err := sequences.RunCreateBotRequesterSequence(sequenceArgs)
		if err != nil {
			return fmt.Errorf("boost request creation failed with bot requester: %v", err)
		}
	}
	return nil
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
	guildMember, err := brm.discord.GuildMember(br.GuildID, userID)
	if err != nil {
		return fmt.Errorf("fetching guild member: %v", err)
	}
	privileges := brm.GetBestRolePrivileges(br.GuildID, guildMember.Roles)
	if privileges == nil {
		return ErrNoPrivileges
	}
	if len(br.PreferredAdvertiserIDs) == 0 {
		brm.signUp(br, userID, privileges)
	} else {
		_, isPreferredAdvertiser := br.PreferredAdvertiserIDs[userID]
		if !isPreferredAdvertiser {
			return ErrNotPreferredAdvertiser
		}

		req, ok := brm.activeRequests.Load(br.ID)
		if !ok {
			log.Printf("AddAdvertiserToBoostRequest: activeRequest not found")
			return
		}
		r := req.(*active_request.ActiveRequest)
		r.SetAdvertiser(userID)
	}
	return nil
}

func (brm *BoostRequestManager) RemoveAdvertiserFromBoostRequest(br *repository.BoostRequest, userID string) (removed bool) {
	req, ok := brm.activeRequests.Load(br.ID)
	if ok {
		ar, ok := req.(*active_request.ActiveRequest)
		if ok && ar.HasSignup(userID) {
			ar.RemoveSignup(userID)
			return true
		}
	}
	return false
}

func (brm *BoostRequestManager) IsAdvertiserSignedUpForBoostRequest(br *repository.BoostRequest, userID string) bool {
	req, ok := brm.activeRequests.Load(br.ID)
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
		_, isPreferredAdvertiser := br.PreferredAdvertiserIDs[userID]
		if !isPreferredAdvertiser {
			return false, false
		}
	}

	credits, err := brm.repo.GetStealCreditsForUser(br.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits: %v", err)
		return false, false
	}
	if credits <= 0 {
		return false, false
	}

	guildMember, err := brm.discord.GuildMember(br.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching guild member: %v", err)
		return false, false
	}
	privileges := brm.GetBestRolePrivileges(br.GuildID, guildMember.Roles)
	if privileges == nil {
		return false, false
	}

	req, ok := brm.activeRequests.Load(br.ID)
	if !ok {
		return false, false
	}
	r := req.(*active_request.ActiveRequest)
	endTime := br.CreatedAt.Add(time.Duration(privileges.Delay) * time.Second)
	now := time.Now()
	ok = r.SetAdvertiser(userID)
	// Don't subtract steal credits if the 👍 button would have had the same effect
	if ok && now.Before(endTime) {
		err := brm.repo.AdjustStealCreditsForUser(br.GuildID, userID, repository.OperationSubtract, 1)
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

	err := brm.cancelDelayedMessages(&br)
	if err != nil {
		log.Printf("error canceling delayed messages: %v", err)
	}

	brm.activeRequests.Delete(br.ID)
	br.AdvertiserID = userID
	br.IsResolved = true
	br.ResolvedAt = time.Now()
	err = brm.repo.ResolveBoostRequest(&br)
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

	if br.Channel == nil || !br.Channel.SkipsBuyerDM {
		_, err = brm.messenger.SendAdvertiserChosenDMToRequester(&br)
		if err != nil {
			log.Printf("Error sending advertiser chosen DM to requester: %v", err)
		}
	}

	rollChannel, err := brm.repo.GetRollChannel(br.GuildID)
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

	// v1
	err = brm.webhookManager.QueueToSend(br.GuildID, &webhook.WebhookEvent{
		Event:   webhook.AdvertiserChosenEvent,
		Payload: models_v1.FromRepositoryBoostRequest(&br),
	})
	if err != nil {
		log.Printf("error queueing webhook: %v", err)
	}

	// v2
	err = brm.webhookManager.QueueToSend(br.GuildID, &webhook.WebhookEvent{
		Event:   webhook.AdvertiserChosenEventV2,
		Payload: models.FromRepositoryBoostRequest(&br),
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
	req, ok := brm.activeRequests.Load(br.ID)
	if !ok {
		return
	}
	r := req.(*active_request.ActiveRequest)
	r.AddSignup(userID, *privileges)
}

func (brm *BoostRequestManager) CancelBoostRequest(br *repository.BoostRequest) error {
	err := brm.cancelDelayedMessages(br)
	if err != nil {
		return err
	}
	for _, message := range br.BackendMessages {
		err := brm.discord.ChannelMessageDelete(message.ChannelID, message.MessageID)
		if err != nil {
			// it's not an ideal situation but we can still continue
			log.Printf("failed to delete backend message when cancelling boost request: %v", err)
		}
	}
	err = brm.repo.DeleteBoostRequest(br)
	if err != nil {
		return err
	}
	activeRequestInterface, loaded := brm.activeRequests.LoadAndDelete(br.ID)
	if !loaded {
		return nil
	}
	activeRequest := activeRequestInterface.(*active_request.ActiveRequest)
	activeRequest.Destroy()
	return nil
}

func (brm *BoostRequestManager) cancelDelayedMessages(br *repository.BoostRequest) error {
	delayedMessageIDs, err := brm.repo.GetBoostRequestDelayedMessageIDs(br)
	if err != nil {
		return err
	}
	for _, id := range delayedMessageIDs {
		err := brm.messenger.CancelDelayedMessage(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (brm *BoostRequestManager) EnableAutoSignup(guildID, userID string, duration time.Duration) error {
	brm.CancelAutoSignup(guildID, userID)

	expiresAt := time.Now().Add(duration)
	autoSignup, err := brm.repo.EnableAutoSignup(guildID, userID, expiresAt)
	if err != nil {
		return err
	}

	delayedMessages, errChannel := brm.messenger.SendAutoSignupMessages(guildID, userID, expiresAt)
	brm.repo.InsertAutoSignupDelayedMessages(autoSignup, delayedMessages)
	go func() {
		for err := range errChannel {
			log.Printf("error sending auto signup expired message: %v", err)
		}
	}()
	return nil
}

func (brm *BoostRequestManager) CancelAutoSignup(guildID, userID string) error {
	err := brm.repo.CancelAutoSignup(guildID, userID)
	if err != nil {
		return err
	}
	ids, err := brm.repo.GetAutoSignupDelayedMessageIDs(guildID, userID)
	if err != nil {
		return err
	}
	for _, id := range ids {
		err := brm.messenger.CancelDelayedMessage(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (brm *BoostRequestManager) applyAutoSignups(br *repository.BoostRequest) {
	autoSignupSessions, err := brm.repo.GetEnabledAutoSignupsInGuild(br.GuildID)
	if err != nil {
		log.Printf("Error fetching auto signup sessions for guild: %v", err)
	}
	for _, s := range autoSignupSessions {
		err := brm.AddAdvertiserToBoostRequest(br, s.AdvertiserID)
		if err != nil {
			log.Printf("Error auto signing up for boost request: %v", err)
		}
	}
}
