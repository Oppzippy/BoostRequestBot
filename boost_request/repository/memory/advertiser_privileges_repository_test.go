package memory_test

import (
	"reflect"
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/memory"
)

func TestGetAdvertiserPrivilegesForGuild(t *testing.T) {
	repo := memory.NewRepository()
	p, err := repo.GetAdvertiserPrivilegesForGuild("")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if p == nil {
		t.Error("returned nil, expected empty slice")
		return
	}
	if len(p) != 0 {
		t.Errorf("expected empty array but it had a length of %d", len(p))
		return
	}
}

func TestGetAdvertiserPrivilegesForRole(t *testing.T) {
	repo := memory.NewRepository()
	p, err := repo.GetAdvertiserPrivilegesForRole("", "")
	if err != repository.ErrNoResults {
		t.Errorf("expected ErrNoResults but received: %v", err)
		return
	}
	if p != nil {
		t.Errorf("expected nil, received: %v", p)
		return
	}
}

func TestInsertAdvertiserPrivileges(t *testing.T) {
	repo := memory.NewRepository()

	adPriv := &repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "Advertiser",
		Weight:  1,
		Delay:   60,
	}
	eliteAdPriv := &repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "Elite Advertiser",
		Weight:  2,
		Delay:   20,
	}

	err := repo.InsertAdvertiserPrivileges(adPriv)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	err = repo.InsertAdvertiserPrivileges(eliteAdPriv)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	if adPriv.ID != 0 {
		t.Errorf("expected adPriv to be id 0 but received id %d", adPriv.ID)
		return
	}
	if eliteAdPriv.ID != 1 {
		t.Errorf("expected adPriv to be id 1 but received id %d", eliteAdPriv.ID)
		return
	}

	retrievedAdPriv, err := repo.GetAdvertiserPrivilegesForRole("guild", "Advertiser")
	if err != nil {
		t.Errorf("expected to be able to retrieve inserted adPriv: %v", err)
		return
	}
	retrievedEliteAdPriv, err := repo.GetAdvertiserPrivilegesForRole("guild", "Elite Advertiser")
	if err != nil {
		t.Errorf("expected to be able to retrieve inserted eliteAdPriv: %v", err)
		return
	}
	if !reflect.DeepEqual(*adPriv, *retrievedAdPriv) {
		t.Error("expected inserted and retrieved advertiser structs to be equal")
		return
	}
	if !reflect.DeepEqual(*eliteAdPriv, *retrievedEliteAdPriv) {
		t.Error("expected inserted and retrieved elite advertiser structs to be equal")
		return
	}
}

func TestDeleteAdvertiserPrivileges(t *testing.T) {
	repo := memory.NewRepository()
	p := &repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "test",
		Weight:  1,
		Delay:   20,
	}
	err := repo.InsertAdvertiserPrivileges(p)
	if err != nil {
		t.Errorf("failed to insert privileges: %v", err)
		return
	}
	err = repo.DeleteAdvertiserPrivileges(p)
	if err != nil {
		t.Errorf("failed to delete privileges: %v", err)
		return
	}
	p, err = repo.GetAdvertiserPrivilegesForRole("guild", "test")
	if err == nil && p != nil {
		t.Errorf("successfuly retrieved privileges that should have been deleted")
	}
}
