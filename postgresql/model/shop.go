package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"

	"ShopServer/postgresql"
)

type Shop struct {
	ID      int    `json:"-"`
	ShopID  string `json:"shop_id,omitempty"`
	Account string `json:"account,omitempty"`
}

func GetShopIDByAccount(ctx context.Context, db postgresql.DB, account string) (*int, error) {
	query := `
		SELECT
			id
		FROM
			shops
		WHERE account = $1
	`

	rows, err := db.Query(ctx, query, account)
	if err != nil {
		return nil, err
	}

	var shopID *int

	err = pgxscan.ScanOne(&shopID, rows)
	if err != nil {
		return nil, err
	}
	return shopID, nil
}

func (s Shop) Upsert(ctx context.Context, db postgresql.DB) error {
	var (
		index  = 1
		values = []string{}
		args   = []interface{}{}
	)

	cols := []string{
		"shop_id",
		"account",
	}

	byted, err := json.Marshal(s)
	if err != nil {
		return err
	}

	var sMap map[string]interface{}

	err = json.Unmarshal(byted, &sMap)
	if err != nil {
		return nil
	}

	for _, col := range cols {
		values = append(values, fmt.Sprintf("$%d", index))
		args = append(args, sMap[col])
		index++
	}

	query := fmt.Sprintf(`
		INSERT INTO
			shops
				(%s)
		VALUES
			(%s)

	`,
		strings.Join(cols, ","),
		strings.Join(values, ","),
	)

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
