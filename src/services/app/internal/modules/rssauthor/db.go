package rssauthor

import (
	"context"

	"github.com/lumi4x/news/src/services/app/internal/ent"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/organization"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/rssauthor"
	"github.com/sirupsen/logrus"
)

func Upsert(ctx context.Context, org *build.Organization, nameToEmail map[string]string) (map[string]*build.RSSAuthor, error) {
	// Populate creators
	var creators []*build.RSSAuthorCreate
	for name, email := range nameToEmail {
		if name == "" {
			logrus.WithField("email", email).Errorf("cannot upsert author with empty name")
			continue
		}

		creator := ent.Database.RSSAuthor.Create().
			SetName(name).
			SetOrganization(org)

		if email != "" {
			creator.SetEmail(email)
		}

		creators = append(creators, creator)
	}

	// Upsert
	err := ent.Database.RSSAuthor.CreateBulk(creators...).
		OnConflictColumns(rssauthor.FieldName, rssauthor.OrganizationColumn).
		UpdateEmail().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Query authors
	var names []string
	for name := range nameToEmail {
		names = append(names, name)
	}

	authors, err := ent.Database.RSSAuthor.Query().
		Where(
			rssauthor.HasOrganizationWith(organization.ID(org.ID)),
			rssauthor.NameIn(names...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Map name to author
	nameToAuthor := map[string]*build.RSSAuthor{}
	for _, author := range authors {
		nameToAuthor[author.Name] = author
	}

	return nameToAuthor, nil
}
