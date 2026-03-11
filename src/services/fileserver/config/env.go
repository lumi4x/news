package config

import (
	_ "embed"

	"github.com/lumi4x/news/src/common/config"
)

var (
	FileServerDatabaseName = config.RequireEnv("FILESERVER_DATABASE_NAME")
)
