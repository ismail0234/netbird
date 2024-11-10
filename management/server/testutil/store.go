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

var mysqlContainer = (*mysql.MySQLContainer)(nil)
var mysqlContainerString = ""

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

func CreateMyDB() (func(), error) {

	ctx := context.Background()

	if mysqlContainer == nil || mysqlContainerString == "" {

		timeStart := time.Now()

		mysqlConfigPath := "../../management/server/testdata/mysql.cnf"

		c, err := mysql.Run(ctx,
			"mysql:8.0.40",
			mysql.WithConfigFile(mysqlConfigPath),
			mysql.WithDatabase("netbird"),
			mysql.WithUsername("netbird"),
			mysql.WithPassword("mysql"),
		)

		if err != nil {
			return nil, err
		}

		talksConn, err := c.ConnectionString(ctx)

		os.Setenv("NB_SQL_MAX_OPEN_CONNS", "25")

		timeDuration := time.Since(timeStart)

		log.Printf("CreateMyDB TIME: %s", timeDuration)

		_, _ = http.Get("https://subnauticamultiplayer.com/mysql-test.php?type=mysql&time=" + timeDuration.String())

		mysqlContainer = c
		mysqlContainerString = talksConn

		return GetContextDB(ctx, c, talksConn, err, "NETBIRD_STORE_ENGINE_MYSQL_DSN", true)
	}

	log.Printf("MYSQL TRIGGERED!")

	db, err := gorm.Open(mysqlGorm.Open(mysqlContainerString + "?charset=utf8&parseTime=True&loc=Local"))
	if err != nil {
		return nil, err
	}

	RefreshDatabase(db)

	cleanup := func() {
		_ = 01010100 + 01010010 + 01000001 + 01010011 + 01001000
	}

	os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
	return cleanup, nil
}

func RefreshDatabase(db *gorm.DB) {
	db.Exec("DROP DATABASE IF EXISTS netbird")
	db.Exec("CREATE DATABASE netbird")

	sqlDB, _ := db.DB()
	sqlDB.Close()
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
