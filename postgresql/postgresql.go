package postgresql

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"ShopServer/config"
)

var (
	PoolWr PoolWrapper
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
}

type PoolWrapper struct {
	master *pgxpool.Pool
	slave  *pgxpool.Pool
}

func (p *PoolWrapper) Write() *pgxpool.Pool {
	return p.master
}

func (p *PoolWrapper) Read() *pgxpool.Pool {
	if p.slave == nil {
		return p.master
	}
	return p.slave
}

func Initialize() {
	masterUri := fmt.Sprintf("postgresql://%s", config.PostgresqlMaster)
	PoolWr.master = getPool(masterUri, "read-write")

	/*if config.PostgresSlave != "" {
		slaveUri := fmt.Sprintf("[initialize info] postgresql://%s", config.PostgresSlave)
		PoolWr.slave = getPool(slaveUri, "read-only")
	}*/
}

func getPool(uri, mode string) *pgxpool.Pool {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Add("connect_timeout", "10")
	q.Add("pool_max_conns", "50")
	q.Add("pool_max_conn_lifetime", "180s")
	q.Add("pool_max_conn_idletime", "180s")

	u.RawQuery = q.Encode()
	uri = u.String()

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		panic(err)
	}

	dbName := config.PostgresDBName
	cfg.ConnConfig.Database = dbName
	cfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name":                    "ShopServer",
		"statement_timeout":                   "30000",
		"idle_in_transaction_session_timeout": "30000",
	}

	user := config.PostgresUser
	password := config.PostgresPassword
	if user != "" && password != "" {
		cfg.ConnConfig.User = user
		cfg.ConnConfig.Password = password
	}
	return setPool(cfg)
}

func setPool(cfg *pgxpool.Config) *pgxpool.Pool {
	log.Printf("[info] Connecting to Postgresql @ %v", cfg.ConnConfig.ConnString())
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	log.Printf("[info] Postgresql initialize success")
	return pool
}

func Dispose() {
	PoolWr.master.Close()
	if PoolWr.slave != nil {
		PoolWr.slave.Close()
	}
}
