package database

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	common_config "github.com/lumi4x/news/src/common/config"
	"github.com/sirupsen/logrus"
)

func init() {
	err := os.Setenv("APP_ENV", string(common_config.Development))
	if err != nil {
		logrus.WithError(err).Panicf("failed setting env for APP_ENV")
	}
}

func CreateMigration(serviceDirName, databaseName, migrationName string, connectFunc func(context.Context) error, namedDiffFunc func(context.Context, string, string, ...schema.MigrateOption) error) error {
	common_config.PostgresHost = "localhost"
	common_config.PostgresPort = "5433"
	common_config.AppEnv = common_config.Development

	rootPath, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("failed getting working directory")
	}

	dir, err := sqltool.NewGolangMigrateDir(path.Join(rootPath, common_config.ProjectMigrationsDirectory(serviceDirName)))
	if err != nil {
		logrus.Fatalf("failed creating atlas migration directory: %v", err)
	}

	// Migrate diff options
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}
	if len(os.Args) != 2 {
		return fmt.Errorf("migration name is required. Use: 'go run -mod=mod ./src/services/<service_name>/ent/migrate/main.go <name>'")
	}

	//#nosec
	cmd := exec.Command("docker", "run", "-d", "--rm", "--name", "migration", "-e", "POSTGRES_DB="+databaseName, "-e", "POSTGRES_PASSWORD="+common_config.PostgresPassword, "-e", "POSTGRES_USER="+common_config.PostgresUser, "-p", "5433:5432", "postgres")
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed creating postgres instance: %v", err)
	}

	containerId := strings.Trim(strings.TrimSpace(buf.String()), "\n")

	//#nosec
	cmd = exec.Command("docker", "rm", "-f", containerId)
	cmd.Stderr = os.Stderr

	defer func() {
		err = cmd.Run()
		if err != nil {
			logrus.WithError(err).Error("failed removing postgres instance")
		}
	}()

	time.Sleep(2 * time.Second)

	connectionString := common_config.PostgresConnectionString(databaseName)
	err = namedDiffFunc(context.Background(), connectionString, migrationName, opts...)
	if err != nil {
		logrus.Errorf("failed generating migration file: %v", err)

		if strings.Contains(err.Error(), "postgres: scanning system variables: EOF") {
			logrus.Panicf("This is a known issue - please re-run, and if it still fails restart your machine")
		}

		absPath, _ := filepath.Abs(fmt.Sprintf("./src/services/%s/internal/ent/migrations", serviceDirName))
		logrus.Panicf("\n\nIf this error is '\n\nIf this error is a checksum mismatch, please run the following command to fix it:\n" +
			"atlas migrate hash --dir file://" + absPath)
	}

	logrus.Info("Created migration. If no migration was created, make sure you ran 'go generate ./...' after modifying the schema")

	return nil
}
