package config

import "github.com/spf13/viper"

var (
	RunMode         string
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int

	OnlineShopServerUrl string
	ShopServerUrl       string
	UserServerUrl       string

	PostgresqlMaster string
	PostgresSlave    string
	PostgresDBName   string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string

	LocalServerName string
	LocalServerUrl  string
)

func Initialize() {
	RunMode = viper.GetString("RUN_MODE")

	Port = viper.GetString("PORT")

	ReadTimeout = viper.GetInt("READ_TIMEOUT")
	WriteTimeout = viper.GetInt("WRITE_TIMEOUT")
	ShutdownTimeout = viper.GetInt("SHUTDOWN_TIMEOUT")

	OnlineShopServerUrl = viper.GetString("ONLINE_SHOP_SERVER_URL")
	ShopServerUrl = viper.GetString("SHOP_SERVER_URL")
	UserServerUrl = viper.GetString("USER_SERVER_URL")

	PostgresDBName = viper.GetString("POSTGRES_DB_NAME")
	PostgresSlave = viper.GetString("POSTGRES_SLAVE")
	PostgresUser = viper.GetString("POSTGRES_USER")
	PostgresPassword = viper.GetString("POSTGRES_PASSWORD")
	PostgresqlMaster = viper.GetString("POSTGRES_MASTER")
	PostgresHost = viper.GetString("POSTGRES_HOST")

	LocalServerName = viper.GetString("LOCAL_SERVER_NAME")
	LocalServerUrl = viper.GetString("LOCAL_SERVER_URL")
}
