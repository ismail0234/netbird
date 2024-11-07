//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"os"
	"path/filepath"
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

func fileExists(filename string) string {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return "NO"
	}

	if info.IsDir() {
		return "YES"
	}

	return "NO"
}

func CreateMyDB() (func(), error) {

	filename := "mysql.cnf"

	log.Printf("mysql: %s", fileExists(filepath.Join("..", "server", "testdata", filename)))
	log.Printf("mysql 1: %s", fileExists(filepath.Join("server", "testdata", filename)))
	log.Printf("mysql 2: %s", fileExists(filepath.Join("..", "testdata", filename)))
	log.Printf("mysql 3: %s", fileExists(filepath.Join("testdata", filename)))
	log.Printf("mysql 4: %s", fileExists(filepath.Join(filename)))

	ctx := context.Background()
	c, err := mysql.Run(ctx,
		"mysql:8.0.40",
		mysql.WithConfigFile(filepath.Join("..", "testdata", filename)),
		mysql.WithDatabase("netbird"),
		mysql.WithUsername("netbird"),
		mysql.WithPassword("mysql"),
	)

	if err != nil {
		return nil, err
	}

	talksConn, err := c.ConnectionString(ctx)

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
