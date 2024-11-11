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
)

var (
	mysqlContainer           = (*mysql.MySQLContainer)(nil)
	mysqlContainerString     = ""
	mysqlContainerConfigPath = "../../management/server/testdata/mysql.cnf"
	postgresContainer        = (*postgres.PostgresContainer)(nil)
	postgresContainerString  = ""
)

func emptyCleanup() {
	// Empty function, don't do anything.
}

func CreateMysqlTestContainer() (func(), error) {

	ctx := context.Background()

	if mysqlContainerString != "" && mysqlContainer != nil && mysqlContainer.IsRunning() {
		RefreshMysqlDatabase(ctx)
		return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
	}

	container, err := mysql.Run(ctx,
		"mysql:8.0.40",
		mysql.WithConfigFile(mysqlContainerConfigPath),
		mysql.WithDatabase("netbird"),
		mysql.WithUsername("root"),
		mysql.WithPassword(""),
	)

	if err != nil {
		return nil, err
	}

	talksConn, _ := container.ConnectionString(ctx)

	mysqlContainer = container
	mysqlContainerString = talksConn

	RefreshMysqlDatabase(ctx)
	return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", talksConn)
}

func RefreshMysqlDatabase(ctx context.Context) {
	mysqlContainer.Exec(ctx, []string{"mysqladmin", "--user=root", "drop", "netbird", "-f"})
	mysqlContainer.Exec(ctx, []string{"mysqladmin", "--user=root", "create", "netbird"})
}

func CreatePostgresTestContainer() (func(), error) {

	ctx := context.Background()
	c, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("netbird"),
		postgres.WithUsername("root"),
		postgres.WithPassword("netbird"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	talksConn, _ := c.ConnectionString(ctx)

	postgresContainer = c

	log.Printf("OUTPUT: %s", execInPostgresContainer([]string{"dropdb", "-f", "netbird"}))

	log.Fatalf("FATAL")
	return GetContextDB(ctx, c, talksConn, err, "NETBIRD_STORE_ENGINE_POSTGRES_DSN", false)
}

func execInPostgresContainer(commands []string) string {
	_, reader, _ := postgresContainer.Exec(context.Background(), commands)

	buf := new(strings.Builder)
	_, errx := io.Copy(buf, reader)

	if errx != nil {
		return "[ERR]"
	}

	return buf.String()
}

func GetContextDB(ctx context.Context, c testcontainers.Container, talksConn string, err error, dsn string, clearCleanUp bool) (func(), error) {

	cleanup := func() {
		timeout := 10 * time.Second
		err = c.Stop(ctx, &timeout)
		if err != nil {
			//	log.WithContext(ctx).Warnf("failed to stop container: %s", err)
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
