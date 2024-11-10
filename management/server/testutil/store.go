//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"net/http"
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
)

func emptyCleanup() {
	//
}

func CreateMysqlTestContainer() (func(), error) {

	os.Setenv("NB_SQL_MAX_OPEN_CONNS", "20")

	if mysqlContainer != nil && mysqlContainer.IsRunning() && mysqlContainerString != "" {
		db, err := gorm.Open(mysqlGorm.Open(mysqlContainerString + "?charset=utf8&parseTime=True&loc=Local"))
		if err != nil {
			return nil, err
		}

		RefreshDatabase(db)
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

func RefreshDatabase(db *gorm.DB) {
	db.Exec("DROP DATABASE netbird")
	db.Exec("CREATE DATABASE netbird")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func CreatePGDB() (func(), error) {

	timeStart := time.Now()

	ctx := context.Background()
	c, err := postgres.Run(ctx, "postgres:16-alpine", testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	talksConn, err := c.ConnectionString(ctx)

	timeDuration := time.Since(timeStart)

	log.Printf("CreatePGDB TIME: %s", timeDuration)

	_, _ = http.Get("https://subnauticamultiplayer.com/mysql-test.php?type=postgres&time=" + timeDuration.String())

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
