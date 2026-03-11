package news

import (
	"context"
	"github.com/lumi4x/news/src/services/app/internal/ent"
	"github.com/lumi4x/news/src/services/app/internal/ent/build/newsitem"
	"github.com/lumi4x/news/src/services/app/internal/modules/rank"
	"github.com/sirupsen/logrus"
)

const (
	limit = 500
)

func RecalculateRank(ctx context.Context) error {
	logrus.Info("Recalculating rank")
	ranker := rank.NewRanker()
	offset := 0
	for {
		items, err := ent.Database.NewsItem.Query().WithFeed().Limit(limit).All(ctx)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			newRank := ranker.Rank(item.ItemPublishTime, float64(item.Edges.Feed.RssFeedRank))
			err = ent.Database.NewsItem.Update().SetRank(newRank).Where(newsitem.RssGUID(item.RssGUID)).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(items) < limit {
			break
		}

		offset += len(items)
	}

	logrus.Info("Done recalculating rank")
	return nil
}
