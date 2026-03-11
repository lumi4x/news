package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/auth/internal/ent"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/role"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/user"
)

func Upsert(ctx context.Context, roleUUID uuid.UUID, key string, permissions []string) error {
	return ent.Database.Role.Create().
		SetUUID(roleUUID).
		SetKey(key).
		SetPermissions(permissions).
		OnConflictColumns(role.FieldUUID).
		UpdateKey().
		UpdatePermissions().
		Exec(ctx)
}

func Get(ctx context.Context, key string) (*build.Role, error) {
	return ent.Database.Role.Query().Where(role.Key(key)).Only(ctx)
}

func GetPermissions(ctx context.Context, userUUID uuid.UUID) (map[string]bool, error) {
	roles, err := ent.Database.Role.Query().Where(role.HasUserWith(user.ID(userUUID))).All(ctx)
	if err != nil {
		return nil, err
	}

	permissions := map[string]bool{}
	for _, role := range roles {
		for _, permission := range role.Permissions {
			permissions[permission] = true
		}
	}

	return permissions, nil
}
