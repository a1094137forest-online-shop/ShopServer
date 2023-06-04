package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"

	"ShopServer/postgresql"
)

type Product struct {
	ID          int    `json:"-"`
	ProductID   string `json:"product_id" db:"product_id"`
	ShopID      int    `json:"shop_id" db:"shop_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`

	Created_at *time.Time `json:"created_at" db:"created_at"`
}

func GetProductsList(ctx context.Context, db postgresql.DB, shopID int) (*[]Product, error) {
	var ps []Product

	query := `
		SELECT
			product_id, title, description, created_at
		FROM
			products
		WHERE
			shop_id = $1
	`

	rows, err := db.Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&ps, rows)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func GetProduct(ctx context.Context, db postgresql.DB, productID string) (*Product, error) {
	var p Product
	cols := []string{
		"product_id",
		"title",
		"description",
		"created_at",
	}

	query := fmt.Sprintf(`
		SELECT
			%s
		FROM
			products
		WHERE
			product_id = $1
	`,
		strings.Join(cols, ","),
	)

	rows, err := db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanOne(&p, rows)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Product) Upsert(ctx context.Context, db postgresql.DB) error {
	var (
		index  = 1
		values = []string{}
		exCols []string
		args   []interface{}
	)

	cols := []string{
		"product_id",
		"shop_id",
		"title",
		"description",
	}

	byted, _ := json.Marshal(p)

	var setMap map[string]interface{}

	_ = json.Unmarshal(byted, &setMap)

	for _, col := range cols {
		exCols = append(exCols, "EXCLUDED."+col)
		values = append(values, fmt.Sprintf("$%d", index))
		args = append(args, setMap[col])
		index++
	}

	query := fmt.Sprintf(`
		INSERT INTO
			products
				(%s)
		VALUES
			(%s)
		ON CONFLICT
			(product_id)
		DO UPDATE
			SET
				(%s) = (%s)
	`,
		strings.Join(cols, ","),
		strings.Join(values, ","),
		strings.Join(cols, ","),
		strings.Join(exCols, ","),
	)

	_, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete(ctx context.Context, db postgresql.DB) error {
	query := `
		DELETE
			FROM
				products
		WHERE
			shop_id = $1 AND product_id = $2
	`

	_, err := db.Exec(ctx, query, p.ShopID, p.ProductID)
	if err != nil {
		return err
	}

	return nil
}
