package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/app/internal/ent/mixins"
)

type NewsItem struct {
	ent.Schema
}

func (NewsItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.CreateTime{},
		mixin.UpdateTime{},
		mixins.SoftDeleteMixin{},
	}
}

func (NewsItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("rss_guid").Unique(),
		index.Fields("rank"),
	}
}

func (NewsItem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New).Unique().StorageKey("uuid"),
		field.String("rss_guid").Unique(),
		field.String("title"),
		field.String("description"),
		field.String("content"),
		field.String("link"),
		field.Strings("links"),
		field.Time("item_publish_time").Optional(),
		field.Time("item_update_time").Optional(),
		field.String("image_url"),
		field.String("image_title"),
		field.Strings("categories"),
		field.Bool("blur").Default(false),
		field.Bool("paywalled"),
		field.Int64("rank"),
	}
}

func (NewsItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("authors", RSSAuthor.Type),
		edge.To("feed", RSSFeed.Type).Unique(),
	}
}
