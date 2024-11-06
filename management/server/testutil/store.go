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
)

func CreatePGDB() (func(), error) {

	log.Printf("[DEBUG] CreatePGDB")

	ctx := context.Background()
	c, err := postgres.Run(ctx, "postgres:16-alpine", testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	talksConn, err := c.ConnectionString(ctx)

	return GetContextDB(ctx, c, talksConn, err, "NETBIRD_STORE_ENGINE_POSTGRES_DSN")
}

func CreateMyDB() (func(), error) {

	log.Printf("[DEBUG] CreateMyDB")

	ctx := context.Background()
	c, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("foo"),
		mysql.WithUsername("root"),
		mysql.WithPassword(""),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).WithStartupTimeout(30*time.Second)),
	)

	log.Printf("[DEBUG] CreateMyDB - 2")

	if err != nil {
		log.Printf("[DEBUG] CreateMyDB Error: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] CreateMyDB - 3")

	talksConn, err := c.ConnectionString(ctx)

	log.Printf("[DEBUG] CreateMyDB - ConnectionString: %s, Error: %s", talksConn, err)

	return GetContextDB(ctx, c, talksConn, err, "NETBIRD_STORE_ENGINE_MYSQL_DSN")
}

func GetContextDB(ctx context.Context, c testcontainers.Container, talksConn string, err error, dsn string) (func(), error) {

	cleanup := func() {
		timeout := 10 * time.Second
		err = c.Stop(ctx, &timeout)
		if err != nil {
			log.WithContext(ctx).Warnf("failed to stop container: %s", err)
		}
	}

	if err != nil {
		return cleanup, err
	}

	return cleanup, os.Setenv(dsn, talksConn)
}
