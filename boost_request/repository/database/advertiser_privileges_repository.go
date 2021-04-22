package database

import (
	"database/sql"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetAdvertiserPrivilegesForGuild(guildID string) ([]*repository.AdvertiserPrivileges, error) {
	privileges, err := repo.getAdvertiserPrivileges("WHERE guild_id = ? ORDER BY weight DESC", guildID)
	return privileges, err
}

func (repo *dbRepository) GetAdvertiserPrivilegesForRole(guildID, roleID string) (*repository.AdvertiserPrivileges, error) {
	privileges, err := repo.getAdvertiserPrivileges("WHERE guild_id = ? AND role_id = ?", guildID, roleID)
	if err != nil {
		return nil, err
	}
	if len(privileges) == 0 {
		return nil, repository.ErrNoResults
	}
	return privileges[0], nil
}

func (repo *dbRepository) getAdvertiserPrivileges(where string, args ...interface{}) ([]*repository.AdvertiserPrivileges, error) {
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

	privileges := make([]*repository.AdvertiserPrivileges, 0, 15)
	for res.Next() {
		p := repository.AdvertiserPrivileges{}
		res.Scan(&p.ID, &p.GuildID, &p.RoleID, &p.Weight, &p.Delay)
		privileges = append(privileges, &p)
	}
	if res.Err() != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return privileges, nil
}

func (repo *dbRepository) InsertAdvertiserPrivileges(privileges *repository.AdvertiserPrivileges) error {
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

func (repo *dbRepository) DeleteAdvertiserPrivileges(privileges *repository.AdvertiserPrivileges) error {
	_, err := repo.db.Exec(`DELETE FROM advertiser_privileges WHERE id = ?`, privileges.ID)
	return err
}
