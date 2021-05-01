package database

import (
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetBestDiscountsForRoles(guildID string, roleIDs []string) ([]*repository.RoleDiscount, error) {
	if len(roleIDs) == 0 {
		return make([]*repository.RoleDiscount, 0), nil
	}
	args := make([]interface{}, 0, len(roleIDs)+1)
	args = append(args, guildID)
	for _, roleID := range roleIDs {
		args = append(args, roleID)
	}

	rows, err := repo.db.Query(
		`SELECT
			d.id,
			d.guild_id,
			d.role_id,
			d.boost_type,
			d.discount
		FROM
			(
			SELECT
				rd.id,
				rd.guild_id,
				rd.role_id,
				rd.boost_type,
				rd.discount,
				RANK() OVER(PARTITION BY rd.boost_type
				ORDER BY
					rd.discount DESC,
					rd.id ASC) AS discount_rank
				FROM
					role_discount rd
				WHERE
					rd.guild_id = ?
					AND rd.role_id IN `+SQLSet(len(roleIDs))+`
					AND not_deleted = 1
			) d
		WHERE
			d.discount_rank = 1`,
		args...,
	)
	if err != nil {
		return nil, err
	}

	discounts := make([]*repository.RoleDiscount, 0, 10)
	for rows.Next() {
		var rd repository.RoleDiscount
		err := rows.Scan(&rd.ID, &rd.GuildID, &rd.RoleID, &rd.BoostType, &rd.Discount)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, &rd)
	}
	return discounts, nil
}

func (repo *dbRepository) GetRoleDiscountsForRole(guildID, roleID string) ([]*repository.RoleDiscount, error) {
	discounts, err := repo.getRoleDiscounts("WHERE guild_id = ? AND role_id = ? AND not_deleted = 1", guildID, roleID)
	if err != nil {
		return nil, err
	}
	return discounts, nil
}

func (repo *dbRepository) GetRoleDiscountForBoostType(guildID, roleID, boostType string) (*repository.RoleDiscount, error) {
	discounts, err := repo.getRoleDiscounts("WHERE guild_id = ? AND role_id = ? AND boost_type = ? AND not_deleted = 1", guildID, roleID, boostType)
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
	rows, err := repo.db.Query("SELECT id, guild_id, role_id, boost_type, discount FROM role_discount "+where, args...)
	if err != nil {
		return nil, err
	}
	discounts := make([]*repository.RoleDiscount, 0, 1)
	for rows.Next() {
		var rd repository.RoleDiscount
		err := rows.Scan(&rd.ID, &rd.GuildID, &rd.RoleID, &rd.BoostType, &rd.Discount)
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
		"UPDATE role_discount SET deleted_at = ? WHERE guild_id = ? AND role_id = ? AND boost_type = ? AND not_deleted = 1",
		time.Now().UTC(),
		rd.GuildID,
		rd.RoleID,
		rd.BoostType,
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}
		return err
	}
	res, err := tx.Exec(
		`INSERT INTO role_discount (guild_id, role_id, boost_type, discount, created_at) 
			VALUES (?, ?, ?, ?, ?)`,
		rd.GuildID,
		rd.RoleID,
		rd.BoostType,
		rd.Discount,
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
