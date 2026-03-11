package main

import (
	"context"

	"github.com/lumi4x/news/src/common/config"
	"github.com/lumi4x/news/src/common/maingroup"
	"github.com/lumi4x/news/src/services/auth/internal/ent"
	"github.com/lumi4x/news/src/services/auth/internal/modules/role"
	"github.com/lumi4x/news/src/services/auth/internal/modules/user"
	"github.com/lumi4x/news/src/services/auth/internal/router"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("Starting auth service in %s mode", config.AppEnv)
	defer logrus.Infof("Finished gracefully shutting down")

	ctx := context.Background()

	err := ent.Connect(ctx)
	if err != nil {
		logrus.WithError(err).Panicf("failed connecting to database")
	}

	err = role.Initialize(ctx)
	if err != nil {
		logrus.WithError(err).Panicf("failed to initialize role module")
	}

	err = user.Initialize(ctx)
	if err != nil {
		logrus.WithError(err).Panicf("failed to initialize user module")
	}

	mainGroup := maingroup.New(ctx)

	mainGroup.Go(router.ListenAndServe)

	logrus.Panic(mainGroup.Wait())
}
