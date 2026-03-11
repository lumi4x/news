package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/auth/client"
	"github.com/lumi4x/news/src/services/auth/client/models"
	"github.com/pkg/errors"
)

const (
	KeyAdmin = "admin"
)

var roles = []*models.Role{
	{
		Uuid: "3ec70e30-afa7-4e85-b161-c7a71e8c4e0a",
		Key:  KeyAdmin,
		Permissions: []string{
			client.KeyNewsItemImageBlur,
			client.KeyNewsItemDelete,
		},
	},
}

func Initialize(ctx context.Context) error {
	for _, role := range roles {
		err := Upsert(ctx, uuid.MustParse(role.Uuid), role.Key, role.Permissions)
		if err != nil {
			return errors.Wrapf(err, "failed upserting role")
		}
	}

	return nil
}
