package ent

import (
	"context"
	"fmt"

	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"github.com/sirupsen/logrus"
)

func WithTx(ctx context.Context, client *build.Client, fn func(tx *build.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			err = tx.Rollback()
			if err != nil {
				logrus.WithError(err).Error("rollback failed")
			}
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
