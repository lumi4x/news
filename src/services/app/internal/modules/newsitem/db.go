package newsitem

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/app/internal/ent"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/newsitem"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/organization"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/predicate"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/rssfeed"
)

func Get(ctx context.Context, itemUUID uuid.UUID) (*build.NewsItem, error) {
	return ent.Database.NewsItem.Get(ctx, itemUUID)
}

// RankByGUIDBestEffort gets each item's rank, but doesn't fail if not found.
func RankByGUIDBestEffort(ctx context.Context, guids []string) (map[string]int64, error) {
	items, err := ent.Database.NewsItem.
		Query().
		Where(newsitem.RssGUIDIn(guids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	guidExists := map[string]int64{}
	for _, item := range items {
		guidExists[item.RssGUID] = item.Rank
	}

	return guidExists, nil
}

func GetByOrganizations(ctx context.Context, orgUUIDs []uuid.UUID) ([]*build.NewsItem, error) {
	return ent.Database.NewsItem.Query().
		WithFeed(func(rq *build.RSSFeedQuery) {
			rq.WithOrganization()
		}).
		WithAuthors().
		Where(newsitem.HasFeedWith(rssfeed.HasOrganizationWith(organization.IDIn(orgUUIDs...)))).
		All(ctx)
}

func GetAll(ctx context.Context, includePaywalled bool, limit int, orderBy string, predicates ...predicate.NewsItem) ([]*build.NewsItem, error) {
	if !includePaywalled {
		predicates = append(predicates, newsitem.Paywalled(false))
	}

	return ent.Database.NewsItem.
		Query().
		WithAuthors().
		WithFeed(func(rq *build.RSSFeedQuery) {
			rq.WithOrganization()
		}).
		Where(predicates...).
		Limit(limit).
		Order(build.Desc(orderBy)).
		All(ctx)
}

func ToggleBlur(ctx context.Context, itemUUID uuid.UUID) error {
	newsItem, err := Get(ctx, itemUUID)
	if err != nil {
		return err
	}

	return ent.Database.NewsItem.Update().Where(newsitem.ID(itemUUID)).SetBlur(!newsItem.Blur).Exec(ctx)
}

func SoftDelete(ctx context.Context, itemUUID uuid.UUID) error {
	newsItem, err := Get(ctx, itemUUID)
	if err != nil {
		return err
	}

	return ent.Database.NewsItem.DeleteOne(newsItem).Exec(ctx)
}

func Upsert(ctx context.Context, feedUUID uuid.UUID, rssGuid, title, description, content, link string, links []string, publishTime, updateTime *time.Time, imageUrl, imageTitle string, categories []string, authors []*build.RSSAuthor, isPaywalled bool, deleted bool, rank int64) error {
	if rssGuid == "" {
		return fmt.Errorf("rss guid cannot be empty")
	}

	// NOTE: Don't forget to add Update below as well
	creator := ent.Database.NewsItem.
		Create().
		SetFeedID(feedUUID).
		SetRssGUID(rssGuid).
		SetTitle(title).
		SetDescription(description).
		SetContent(content).
		SetLink(link).
		SetLinks(links).
		SetImageURL(imageUrl).
		SetImageTitle(imageTitle).
		SetPaywalled(isPaywalled).
		SetDeleted(deleted).
		SetCategories(categories).
		SetRank(rank)

	if deleted {
		creator.SetDeleteTime(time.Now())
	}

	if publishTime != nil {
		creator.SetItemPublishTime(*publishTime)
	}

	if updateTime != nil {
		creator.SetItemUpdateTime(*updateTime)
	}

	// Don't override - blur
	err := creator.OnConflictColumns(newsitem.FieldRssGUID).
		UpdateTitle().
		UpdateDescription().
		UpdateContent().
		UpdateLink().
		UpdateLinks().
		UpdateImageURL().
		UpdateImageTitle().
		UpdatePaywalled().
		UpdateCategories().
		UpdateDeleted().
		UpdateDeleteTime().
		//UpdateItemPublishTime(). // Don't update publish time, as if often changes to reflect update time instead
		UpdateItemUpdateTime().
		UpdateRank().
		Exec(ctx)
	if err != nil {
		return err
	}

	err = ent.Database.NewsItem.Update().ClearAuthors().AddAuthors(authors...).Where(newsitem.RssGUID(rssGuid)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRank(ctx context.Context, guid string, rank int64) error {
	return ent.Database.NewsItem.Update().SetRank(rank).Where(newsitem.RssGUID(guid)).Exec(ctx)
}
