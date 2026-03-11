package news

import (
	"context"
	"github.com/lumi4x/news/src/services/app/internal/modules/rank"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"github.com/lumi4x/news/src/services/app/internal/modules/gofeed"
	"github.com/lumi4x/news/src/services/app/internal/modules/html"
	"github.com/lumi4x/news/src/services/app/internal/modules/newsitem"
	"github.com/lumi4x/news/src/services/app/internal/modules/organization"
	"github.com/lumi4x/news/src/services/app/internal/modules/rssauthor"
	"github.com/lumi4x/news/src/services/app/internal/modules/rssfeed"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SyncRSSFeeds(ctx context.Context) error {
	// TODO: Optimize
	//err := RecalculateRank(ctx)
	//if err != nil {
	//	return err
	//}

	err := syncRSSFeeds(ctx)
	if err != nil {
		logrus.WithError(err).Errorf("failed initial syncing of rss feeds")
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := syncRSSFeeds(ctx)
			if err != nil {
				logrus.WithError(err).Errorf("failed syncing rss feeds")
				continue
			}
		}

	}
}

func syncRSSFeeds(ctx context.Context) error {
	orgs, err := organization.GetAll(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed getting orgs for sync")
	}

	ranker := rank.NewRanker()

	logrus.Info("Syncing RSS feeds")

	workGroup := new(sync.WaitGroup)

	for _, org := range orgs {
		workGroup.Add(1)
		go func() {
			defer workGroup.Done()

			feeds, err := rssfeed.GetByOrganization(ctx, org)
			if err != nil {
				logrus.WithError(err).WithField("orgUUID", org.ID.String()).Errorf("failed getting feeds by org for rss feed sync")
				return
			}

			for _, feed := range feeds {
				logrus := logrus.WithField("rssFeedURL", feed.RssFeedURL)
				logrus.Infof("Syncing feed")

				maxTimeAgo := time.Now().Add(-time.Duration(feed.MaxFetchIntervalMin) * time.Minute)
				if feed.UpdateTime.After(maxTimeAgo) {
					logrus.WithField("MaxFetchIntervalMin", feed.MaxFetchIntervalMin).Info("Not updating feed due to max fetch interval")
					continue
				}

				goFeed, err := gofeed.New(feed.RssFeedURL)
				if err != nil {
					logrus.WithError(err).Errorf("failed parsing go feed")
					continue
				}

				guidToRank, err := newsitem.RankByGUIDBestEffort(ctx, goFeed.GetGUIDs())
				if err != nil {
					logrus.WithError(err).Errorf("failed fetching exists by guid")
					continue
				}

				nameToAuthor, err := rssauthor.Upsert(ctx, org, goFeed.GetAuthorsNameToEmail())
				if err != nil {
					logrus.WithError(err).Errorf("failed to bulk upsert authors")
					continue
				}

				// TODO: Optimize by bulking
				for _, item := range goFeed.Items {
					logrus := logrus.WithField("guid", item.GUID).WithField("link", item.Link)

					rank := ranker.Rank(*item.PublishedParsed, float64(feed.RssFeedRank))

					// Minimal update to items that already exist
					if currentRank, exists := guidToRank[item.GUID]; exists {
						if currentRank != rank {
							err := newsitem.UpdateRank(ctx, item.GUID, rank)
							if err != nil {
								logrus.WithError(err).Errorf("failed to update item rank")
								continue
							}
						}
						continue
					}

					if feed.TitleTrimRight != "" {
						item.Title = strings.TrimSuffix(item.Title, feed.TitleTrimRight)
						item.Description = strings.TrimSuffix(item.Title, feed.TitleTrimRight)
					}

					imageUrl := ""
					imageTitle := ""
					if item.Image != nil {
						imageUrl = item.Image.URL
						imageTitle = item.Image.Title
					}

					itemHTML := html.NewHTMLFromURL(item.Link)
					if imageUrl == "" && !feed.DiscardOgImage {
						ogImage := itemHTML.GetOGImage()
						if ogImage != "" {
							imageUrl = ogImage
						}
					}

					// Get DB authors
					var dbAuthors []*build.RSSAuthor
					for _, author := range item.Authors {
						dbAuthors = append(dbAuthors, nameToAuthor[author.Name])
					}

					// Skip based on whitelist
					isDeleted := false
					if feed.ContentWhitelistRegex != "" {
						whitelistRegex := regexp.MustCompile(feed.ContentWhitelistRegex)
						if !whitelistRegex.MatchString(strings.ToLower(item.Title)) && !whitelistRegex.MatchString(strings.ToLower(item.Description)) {
							isDeleted = true
						}
					}

					isPaywalled := false
					if feed.HTMLPaywallRegex != "" {
						htmlPaywallRegex := regexp.MustCompile(feed.HTMLPaywallRegex)

						respContent, err := itemHTML.GetHTML()
						if err != nil {
							logrus.WithError(err).Errorf("failed fetching HTML for item, not including item")
							continue
						}

						isPaywalled = htmlPaywallRegex.MatchString(strings.ToLower(respContent))
					}

					finalUrl, err := itemHTML.GetFinalURL()
					if err != nil {
						logrus.WithError(err).Errorf("failed getting final URL, using default")
						finalUrl = item.Link
					}

					if item.PublishedParsed == nil {
						logrus.Errorf("item does not have publish date, skipping")
						continue
					}

					err = newsitem.Upsert(ctx, feed.ID, item.GUID, item.Title, item.Description, item.Content, finalUrl, item.Links, item.PublishedParsed, item.UpdatedParsed, imageUrl, imageTitle, item.Categories, dbAuthors, isPaywalled, isDeleted, rank)
					if err != nil {
						logrus.WithError(err).WithField("itemGUID", item.GUID).Errorf("failed upserting rss item")
						continue
					}
				}

				err = rssfeed.UpdateUpdateTime(ctx, feed)
				if err != nil {
					logrus.WithError(err).Errorf("failed updating feed update time")
				}

				logrus.Infof("Done syncing feed")
			}
		}()
	}

	workGroup.Wait()
	logrus.Info("Done syncing RSS feeds")
	return nil
}

type ItemHTML struct {
	url  string
	html string
}

func newItemHTML(url string) *ItemHTML {
	return &ItemHTML{
		url: url,
	}
}
