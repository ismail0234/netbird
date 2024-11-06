//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"fmt"
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

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0.40",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          "netbird",
			"MYSQL_PASSWORD":      "mysql",
			"MYSQL_ROOT_PASSWORD": "mysqlroot",
			"MYSQL_DATABASE":      "netbird",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)

	if err != nil {
		log.Printf("[DEBUG] CreateMyDB ErrorX: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] CreateMyDB SCS:")

	if container == nil {
		return nil, nil
	}

	mysql.WithPassword("c")

	talksConn, err := ConnectionString(ctx, container, req.Env["MYSQL_USER"], req.Env["MYSQL_DATABASE"], req.Env["MYSQL_PASSWORD"])

	log.Printf("[DEBUG] CreateMyDB - ConnectionString: %s, Error: %s", talksConn, err)

	return GetContextDB(ctx, container, talksConn, err, "NETBIRD_STORE_ENGINE_MYSQL_DSN")
}

func ConnectionString(ctx context.Context, c testcontainers.Container, user string, db string, pass string) (string, error) {
	containerPort, err := c.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	extraArgs := ""

	if extraArgs != "" {
		extraArgs = "?" + extraArgs
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", user, pass, host, containerPort.Port(), db, extraArgs)
	return connectionString, nil
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
