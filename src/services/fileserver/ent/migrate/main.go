//go:build ignore

package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/lumi4x/news/src/common/database"
	"github.com/lumi4x/news/src/services/fileserver/config"
	"github.com/lumi4x/news/src/services/fileserver/ent"
	"github.com/lumi4x/news/src/services/fileserver/ent/build/migrate"
	"github.com/sirupsen/logrus"
)

// Usage: go generate ./... && go run -mod=mod ./src/services/<service_name>/ent/migrate/main.go migration_name

// To generate an initial migration, run:
// go generate ./... && atlas migrate new --dir file://<path_to_migrations_dir> --dir-format golang-migrate

// To fix checksum mismatch, make sure you don't have conflicts within your migrations, and then run:
// atlas migrate hash --dir file://<path_to_migrations_dir>
func main() {
	err := database.CreateMigration(config.FileServerDirectoryName, config.FileServerDatabaseName, os.Args[1], ent.Connect, migrate.NamedDiff)
	if err != nil {
		logrus.WithError(err).Fatalf("failed running migration")
	}
}
