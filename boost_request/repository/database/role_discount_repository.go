package database

import (
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetRoleDiscountForRole(guildID, roleID string) (*repository.RoleDiscount, error) {
	discounts, err := repo.getRoleDiscounts("WHERE guild_id = ? AND role_id = ? AND not_deleted = 1", guildID, roleID)
	if err != nil {
		return nil, err
	}
	if len(discounts) == 0 {
		return nil, repository.ErrNoResults
	}
	return discounts[0], nil
}

func (repo *dbRepository) GetRoleDiscountsForGuild(guildID string) ([]*repository.RoleDiscount, error) {
	discounts, err := repo.getRoleDiscounts("WHERE guild_id = ? AND not_deleted = 1", guildID)
	return discounts, err
}

func (repo *dbRepository) getRoleDiscounts(where string, args ...interface{}) ([]*repository.RoleDiscount, error) {
	rows, err := repo.db.Query("SELECT id, guild_id, role_id, discount FROM role_discount "+where, args...)
	if err != nil {
		return nil, err
	}
	discounts := make([]*repository.RoleDiscount, 0, 1)
	for rows.Next() {
		var rd repository.RoleDiscount
		err := rows.Scan(&rd.ID, &rd.GuildID, &rd.RoleID, &rd.Discount)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, &rd)
	}
	return discounts, nil
}

// If a discount for the role already exists, it will be replaced.
func (repo *dbRepository) InsertRoleDiscount(rd *repository.RoleDiscount) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"UPDATE role_discount SET deleted_at = ? WHERE guild_id = ? AND role_id = ? AND not_deleted = 1",
		time.Now().UTC(),
		rd.GuildID,
		rd.RoleID,
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}
		return err
	}
	res, err := tx.Exec(
		`INSERT INTO role_discount (guild_id, role_id, discount, created_at) 
			VALUES (?, ?, ?, ?)`,
		rd.GuildID,
		rd.RoleID,
		rd.Discount.String(),
		time.Now().UTC(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	rd.ID = id
	return nil
}

func (repo *dbRepository) DeleteRoleDiscount(rd *repository.RoleDiscount) error {
	_, err := repo.db.Exec("UPDATE role_discount SET deleted_at = ? WHERE id = ?", time.Now().UTC(), rd.ID)
	return err
}
