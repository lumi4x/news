package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/app/client/models"
	"github.com/lumi4x/news/src/services/app/internal/modules/feedtype"
	"github.com/lumi4x/news/src/services/app/internal/modules/news"
	"github.com/lumi4x/news/src/services/app/internal/modules/organization"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	limit = 100
)

func GetNews(w *api.Writer, r *http.Request, p httprouter.Params, j *api.JWT) (*models.GetNewsResponse, *api.Error) {
	feedTypeValue := p.ByName("type")
	feedType, err := feedtype.New(feedTypeValue)
	if err != nil {
		logrus.WithError(err).Errorf("invalid feed type when getting news, falling back to foryou")
		feedType = feedtype.NewForYou()
	}

	// Get featured
	featuredItems, err := feedtype.NewForYou().FetchItems(r.Context(), 1)
	if err != nil {
		return nil, api.NewInternalError(err, "failed fetching featured items")
	}

	if len(featuredItems) == 0 {
		return nil, api.NewNotFoundError(nil, "no featured items found")
	}

	items, err := feedType.FetchItems(r.Context(), limit)
	if err != nil {
		return nil, api.NewInternalError(err, "failed getting news items")
	}

	orgs, err := organization.GetAll(r.Context())
	if err != nil {
		return nil, api.NewInternalError(err, "failed getting orgs")
	}

	response := models.GetNewsResponse{
		Organizations: organization.TranslateOrganizationsToMap(orgs),
		Featured:      news.TranslateRSSItem(featuredItems[0], featuredItems[0].Edges.Feed.Edges.Organization.ID),
		Latest:        &models.RSSFeed{},
	}

	// Translate articles
	var translatedItems []*models.RSSItem
	for _, item := range items {
		// Don't include featured item
		if item.RssGUID == response.Featured.Guid {
			continue
		}

		translatedItems = append(translatedItems, news.TranslateRSSItem(item, item.Edges.Feed.Edges.Organization.ID))
	}

	orgToCount := map[string]int{}
	orgToOldestArticle := map[string]*models.RSSItem{}

	for _, item := range translatedItems {
		if len(response.Latest.Items) >= 3 && len(response.Latest.Items)%3 == 0 && orgToCount[item.OrganizationUuid] >= 3 && orgToOldestArticle[item.OrganizationUuid] != nil {
			orgToOldestArticle[item.OrganizationUuid].SeeAlsoItem = item
			orgToOldestArticle[item.OrganizationUuid] = nil
			continue
		}

		response.Latest.Items = append(response.Latest.Items, item)

		orgToCount[item.OrganizationUuid]++
		orgToOldestArticle[item.OrganizationUuid] = item
	}

	// Limit
	response.Latest.Items = response.Latest.Items[:min(len(response.Latest.Items), 40)]

	return &response, nil
}
