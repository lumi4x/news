package feedtype

import (
	"context"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	dbnewsitem "github.com/lumi4x/news/src/services/app/internal/ent/build/newsitem"
	"github.com/lumi4x/news/src/services/app/internal/modules/newsitem"
)

type Latest struct {
}

func NewLatest() *Latest {
	return &Latest{}
}

func (f *Latest) FetchItems(ctx context.Context, limit int) ([]*build.NewsItem, error) {
	return newsitem.GetAll(ctx, false, limit, dbnewsitem.FieldItemPublishTime)
}
