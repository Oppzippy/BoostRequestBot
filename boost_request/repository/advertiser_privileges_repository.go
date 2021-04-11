package repository

import (
	"database/sql"
	"time"
)

type AdvertiserPrivilegesRepository interface {
	GetAdvertiserPrivilegesForGuild(guildID string) ([]*AdvertiserPrivileges, error)
	GetAdvertiserPrivilegesForRole(guildID, roleID string) (*AdvertiserPrivileges, error)
	InsertAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
	DeleteAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
}

func (repo *dbRepository) GetAdvertiserPrivilegesForGuild(guildID string) ([]*AdvertiserPrivileges, error) {
	privileges, err := repo.getAdvertiserPrivileges("WHERE guild_id = ?", guildID)
	return privileges, err
}

func (repo *dbRepository) GetAdvertiserPrivilegesForRole(guildID, roleID string) (*AdvertiserPrivileges, error) {
	privileges, err := repo.getAdvertiserPrivileges("WHERE guild_id = ? AND role_id = ?", guildID, roleID)
	if err != nil {
		return nil, err
	}
	if len(privileges) == 0 {
		return nil, nil
	}
	return privileges[0], nil
}

func (repo *dbRepository) getAdvertiserPrivileges(where string, args ...interface{}) ([]*AdvertiserPrivileges, error) {
	res, err := repo.db.Query(
		`SELECT
			id,
			guild_id,
			role_id,
			weight,
			delay
		FROM advertiser_privileges `+where,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	privileges := make([]*AdvertiserPrivileges, 0, 15)
	for res.Next() {
		p := AdvertiserPrivileges{}
		res.Scan(&p.ID, &p.GuildID, &p.RoleID, &p.Weight, &p.Delay)
		privileges = append(privileges, &p)
	}
	if res.Err() != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return privileges, nil
}

func (repo *dbRepository) InsertAdvertiserPrivileges(privileges *AdvertiserPrivileges) error {
	res, err := repo.db.Exec(
		`INSERT INTO advertiser_privileges (
			guild_id,
			role_id,
			weight,
			delay,
			created_at
		) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			weight = VALUES(weight),
			delay = VALUES(delay)`,
		privileges.GuildID,
		privileges.RoleID,
		privileges.Weight,
		privileges.Delay,
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	privileges.ID = id
	return nil
}

func (repo *dbRepository) DeleteAdvertiserPrivileges(privileges *AdvertiserPrivileges) error {
	_, err := repo.db.Exec(`DELETE FROM advertiser_privileges WHERE id = ?`, privileges.ID)
	return err
}
