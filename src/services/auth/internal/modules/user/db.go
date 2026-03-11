package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/auth/internal/config"
	"github.com/lumi4x/news/src/services/auth/internal/ent"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/role"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/user"
	"github.com/lumi4x/news/src/services/auth/internal/ent/mixins"
	"golang.org/x/crypto/pbkdf2"
)

func Get(ctx context.Context, email string) (*build.User, error) {
	return ent.Database.User.Query().Where(user.EmailEQ(strings.ToLower(email))).Only(ctx)
}

func GetIncludeBanned(ctx context.Context, email string) (*build.User, error) {
	return ent.Database.User.Query().Where(user.EmailEQ(strings.ToLower(email))).Only(mixins.SkipBanFilter(ctx))
}

func Exist(ctx context.Context, email string) (bool, error) {
	return ent.Database.User.Query().Where(user.EmailEQ(strings.ToLower(email))).Exist(ctx)
}

func SetRoles(ctx context.Context, appliedUser *build.User, roles []*build.Role) error {
	_, err := ent.Database.Role.Delete().Where(role.HasUserWith(user.ID(appliedUser.ID))).Exec(ctx)
	if err != nil {
		return err
	}

	err = ent.Database.User.Update().AddRoles(roles...).Where(user.ID(appliedUser.ID)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Register(ctx context.Context, email, password string) (*build.User, error) {
	if len(password) < config.PasswordMinLength {
		return nil, fmt.Errorf("password too short")
	}

	// Generate a random salt
	salt := make([]byte, config.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, api.NewInternalError(err, "failed generating salt")
	}

	// Hash the password using PBKDF2
	hashedPassword := pbkdf2.Key([]byte(password), salt, config.Iterations, config.PasswordHashLength, sha256.New)

	// Encode the salt and hashed password as base64 strings
	hashedPasswordString := base64.StdEncoding.EncodeToString(hashedPassword)
	saltString := base64.StdEncoding.EncodeToString(salt)

	user, err := ent.Database.User.Create().
		SetEmail(strings.ToLower(email)).
		SetPassword(hashedPasswordString).
		SetSalt(saltString).
		Save(ctx)
	if err != nil {
		return nil, err // Return raw ent error for later parsing
	}

	return user, nil
}
