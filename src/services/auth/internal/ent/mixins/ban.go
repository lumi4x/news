package mixins

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/intercept"
)

// BanMixin implements the ban pattern for schemas.
type BanMixin struct {
	mixin.Schema
}

// Fields of the BanMixin.
func (BanMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("ban_time").Optional(),
	}
}

type banKey struct{}

// SkipBanFilter returns a new context that skips the soft-delete interceptor/mutators.
func SkipBanFilter(parent context.Context) context.Context {
	return context.WithValue(parent, banKey{}, true)
}

// Interceptors of the BanMixin.
func (d BanMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			// Skip ban filter, means include banned entities.
			if skip, _ := ctx.Value(banKey{}).(bool); skip {
				return nil
			}

			d.P(q)
			return nil
		}),
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d BanMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldIsNull(d.Fields()[0].Descriptor().Name),
	)
}
