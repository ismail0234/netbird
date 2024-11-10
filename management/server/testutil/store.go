//go:build !ios
// +build !ios

package testutil

import (
	"context"
	"io"
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

func CreateMysqlTestContainer() (func(), error) {

	os.Setenv("NB_SQL_MAX_OPEN_CONNS", "20")

	if mysqlContainerString != "" && mysqlContainer != nil && mysqlContainer.IsRunning() {
		execInMysqlContainer([]string{"mysqladmin", "--user=root", "drop", "netbird", "-f"})
		execInMysqlContainer([]string{"mysqladmin", "--user=root", "create", "netbird"})
		return emptyCleanup, os.Setenv("NETBIRD_STORE_ENGINE_MYSQL_DSN", mysqlContainerString)
	}

	ctx := context.Background()
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

	execInMysqlContainer([]string{"mysqladmin", "--user=root", "drop", "netbird", "-f"})
	execInMysqlContainer([]string{"mysqladmin", "--user=root", "create", "netbird"})
	//log.Printf("TEST 2: %s", execInMysqlContainer([]string{"mysqladmin", "--user=root", "drop", "netbird", "-f"}))
	//log.Printf("TEST 2: %s", execInMysqlContainer([]string{"mysqladmin", "--user=root", "create", "netbird"}))

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

func execInPostgresContainer(commands []string) string {
	_, reader, _ := postgresContainer.Exec(context.Background(), commands)

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
