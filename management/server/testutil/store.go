//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	mysqlGorm "gorm.io/driver/mysql"
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

	talksConn, err := container.ConnectionString(ctx)

	postgresContainer = container
	postgresContainerString = talksConn

	return GetContextDB(ctx, container, talksConn, err, "NETBIRD_STORE_ENGINE_POSTGRES_DSN", false)
}

func CreateOrGetMysqlTestContainer() (func(), error) {

	if mysqlContainer != nil && mysqlContainer.IsRunning() && mysqlContainerString != "" {

		db, err := gorm.Open(mysqlGorm.Open(mysqlContainerString))
		if err != nil {
			return nil, err
		}

		RefreshDatabase(db, "NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
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

	talksConn, err := container.ConnectionString(ctx)

	os.Setenv("NB_SQL_MAX_OPEN_CONNS", "20")

	mysqlContainer = container
	mysqlContainerString = talksConn

	return GetContextDB(ctx, container, talksConn, err, "NETBIRD_STORE_ENGINE_MYSQL_DSN", true)
}

func RefreshDatabase(db *gorm.DB, dsn string, connectionString string) {
	db.Exec("DROP DATABASE netbird")
	db.Exec("CREATE DATABASE netbird")

	sqlDB, _ := db.DB()
	sqlDB.Close()

	os.Setenv(dsn, connectionString)
}

func CreatePGDB() (func(), error) {

	ctx := context.Background()
	c, err := postgres.Run(ctx, "postgres:16-alpine", testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	talksConn, err := c.ConnectionString(ctx)

	return GetContextDB(ctx, c, talksConn, err, "NETBIRD_STORE_ENGINE_POSTGRES_DSN", false)
}

func GetContextDB(ctx context.Context, c testcontainers.Container, talksConn string, err error, dsn string, clearCleanUp bool) (func(), error) {

	cleanup := func() {
		timeout := 10 * time.Second
		err = c.Stop(ctx, &timeout)
		if err != nil {
			log.WithContext(ctx).Warnf("failed to stop container: %s", err)
		}
	}

	if clearCleanUp {
		cleanup := func() {
			_ = 1
		}

		return cleanup, os.Setenv(dsn, talksConn)
	}

	if err != nil {
		return cleanup, err
	}

	return cleanup, os.Setenv(dsn, talksConn)
}
