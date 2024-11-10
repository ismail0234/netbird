//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	mysqlGorm "gorm.io/driver/mysql"
	postgresGorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	mysqlContainer           = (*mysql.MySQLContainer)(nil)
	mysqlContainerString     = ""
	mysqlContainerConfigPath = "../../management/server/testdata/mysql.cnf"
	postgresContainer        = (*postgres.PostgresContainer)(nil)
	postgresContainerString  = ""
)

func emptyCleanUp() {
	_ = 01010100 + 01010010 + 01000001 + 01010011 + 01001000
}

func CreatePostgresTestContainer() (func(), error) {
	if postgresContainer != nil && postgresContainer.IsRunning() && postgresContainerString != "" {

		db, err := gorm.Open(postgresGorm.Open(postgresContainerString))
		if err != nil {
			return nil, err
		}

		RefreshContainer(db, "NETBIRD_STORE_ENGINE_POSTGRES_DSN", postgresContainerString)
		return emptyCleanUp, nil
	}

	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:16-alpine", testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	talksConn, _ := container.ConnectionString(ctx)

	postgresContainer = container
	postgresContainerString = talksConn

	return emptyCleanUp, os.Setenv("NETBIRD_STORE_ENGINE_POSTGRES_DSN", talksConn)
}

func CreateMysqlTestContainer() (func(), error) {

	if mysqlContainer != nil && mysqlContainer.IsRunning() && mysqlContainerString != "" {

		db, err := gorm.Open(mysqlGorm.Open(mysqlContainerString))
		if err != nil {
			return nil, err
		}

		RefreshContainer(db, "NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
		return emptyCleanUp, nil
	}

	ctx := context.Background()
	container, err := mysql.Run(ctx,
		"mysql:8.0.40",
		mysql.WithConfigFile(mysqlContainerConfigPath),
		mysql.WithDatabase("netbird"),
		mysql.WithUsername("netbird"),
		mysql.WithPassword("mysql"),
	)

	if err != nil {
		return nil, err
	}

	talksConn, _ := container.ConnectionString(ctx)

	os.Setenv("NB_SQL_MAX_OPEN_CONNS", "20")

	mysqlContainer = container
	mysqlContainerString = talksConn

	return emptyCleanUp, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", talksConn)
}

func RefreshContainer(db *gorm.DB, dsn string, connectionString string) {
	db.Exec("DROP DATABASE netbird")
	db.Exec("CREATE DATABASE netbird")

	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Setenv(dsn, connectionString)
}
