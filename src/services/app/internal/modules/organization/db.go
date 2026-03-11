package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/app/internal/ent"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/organization"
	"github.com/pkg/errors"
)

func GetAll(ctx context.Context) ([]*build.Organization, error) {
	posts, err := ent.Database.Organization.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting orgs from db")
	}

	return posts, nil
}

func UpsertOrganization(ctx context.Context, orgUUID uuid.UUID, uniqueName, name, description, websiteUrl string, logoImageURL string) error {
	creator := ent.Database.Organization.Create().
		SetID(orgUUID).
		SetName(name).
		SetUniqueName(uniqueName).
		SetDescription(description).
		SetWebsiteURL(websiteUrl).
		OnConflictColumns(organization.FieldID).
		UpdateName().
		UpdateDescription().
		UpdateWebsiteURL()

	if logoImageURL != "" {
		creator.SetLogoImageURL(logoImageURL)
	}

	err := creator.Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed upserting organization in db")
	}

	return nil
}
