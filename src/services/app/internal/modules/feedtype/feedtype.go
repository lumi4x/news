package feedtype

import (
	"context"
	"fmt"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
)

type FeedType interface {
	FetchItems(ctx context.Context, limit int) ([]*build.NewsItem, error)
}

var feedTypes = map[string]FeedType{
	"foryou": NewForYou(),
	"latest": NewLatest(),
}

func New(feedType string) (FeedType, error) {
	feed, ok := feedTypes[feedType]
	if !ok {
		return nil, fmt.Errorf("invalid feed type: %s", feedType)
	}

	return feed, nil
}
