//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	mysqlContainer           = (*mysql.MySQLContainer)(nil)
	mysqlContainerString     = ""
	mysqlContainerConfigPath = "../../management/server/testdata/mysql.cnf"
	postgresContainer        = (*postgres.PostgresContainer)(nil)
	postgresContainerString  = ""
)

func emptyCleanup() {
	//
}

func CreatePostgresTestContainer() (func(), error) {

	if postgresContainer != nil && postgresContainer.IsRunning() && postgresContainerString != "" {
		/*	mysqlConnection err := gorm.Open(postgresGorm.Open(postgresContainerString))
			if err != nil {
				return nil, err
			}

			RefreshDatabase(db)*/
		return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_POSTGRES_DSN", postgresContainerString)
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

	return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_POSTGRES_DSN", talksConn)
}

func CreateMysqlTestContainer() (func(), error) {

	os.Setenv("NB_SQL_MAX_OPEN_CONNS", "20")

	if mysqlContainerString != "" && mysqlContainer != nil && mysqlContainer.IsRunning() {

		log.Printf("TEST 1: %s", execInMysqlContainer([]string{"mysqladmin --user=root --password=mysql drop netbird", "/root"}))
		log.Printf("TEST 2: %s", execInMysqlContainer([]string{"mysqladmin --user=root --password=mysql version", "/"}))
		log.Fatal("FATAL ERROR! => ")

		return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
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

	mysqlContainer = container
	mysqlContainerString = talksConn

	return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", talksConn)
}

func execInMysqlContainer(commands []string) string {
	_, reader, _ := mysqlContainer.Exec(context.Background(), commands)

	buf := new(strings.Builder)
	_, errx := io.Copy(buf, reader)

	if errx != nil {
		return "[ERR]"
	}

	return buf.String()
}

func RefreshDatabase(db *gorm.DB) {
	db.Exec("DROP DATABASE IF EXISTS netbird")
	db.Exec("CREATE DATABASE netbird")
}
