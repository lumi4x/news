package news

import (
	"html"
	"regexp"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/app/client/models"
	"github.com/lumi4x/news/src/services/app/internal/ent/build"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var htmlTagRegex = regexp.MustCompile("<[^>]*>")

func TranslateRSSItems(items []*build.NewsItem, orgUUID uuid.UUID) []*models.RSSItem {
	var translated []*models.RSSItem
	for _, item := range items {
		translated = append(translated, TranslateRSSItem(item, orgUUID))
	}
	return translated
}

func translateHTML(htmlText string) string {
	htmlText = html.UnescapeString(htmlText)
	htmlText = htmlTagRegex.ReplaceAllString(htmlText, "")
	return htmlText
}

func TranslateRSSItem(item *build.NewsItem, orgUUID uuid.UUID) *models.RSSItem {
	if item == nil {
		return &models.RSSItem{}
	}

	return &models.RSSItem{
		Title:            item.Title,
		Description:      translateHTML(item.Description),
		Content:          translateHTML(item.Content),
		Link:             item.Link,
		Links:            item.Links,
		UpdateDate:       timestamppb.New(item.ItemUpdateTime),
		PublishDate:      timestamppb.New(item.ItemPublishTime),
		Authors:          TranslateRSSAuthors(item.Edges.Authors),
		Guid:             item.RssGUID,
		Image:            TranslateRSSImage(item.ImageURL, item.ImageTitle, item.Blur),
		Categories:       item.Categories,
		OrganizationUuid: orgUUID.String(),
		Id:               item.ID.String(),
	}
}

func TranslateRSSAuthors(authors []*build.RSSAuthor) []*models.RSSAuthor {
	var translated []*models.RSSAuthor
	for _, author := range authors {
		translated = append(translated, TranslateRSSAuthor(author))
	}
	return translated
}

func TranslateRSSAuthor(author *build.RSSAuthor) *models.RSSAuthor {
	if author == nil {
		return &models.RSSAuthor{}
	}

	return &models.RSSAuthor{
		Email: author.Email,
		Name:  author.Name,
	}
}

func TranslateRSSImage(imageUrl, imageTitle string, blur bool) *models.RSSImage {
	return &models.RSSImage{
		Url:   imageUrl,
		Title: imageTitle,
		Blur:  blur,
	}
}
