package config

import (
	_ "embed"

	"github.com/lumi4x/news/src/common/config"
)

var (
	AppDatabaseName = config.RequireEnv("APP_DATABASE_NAME")
)
