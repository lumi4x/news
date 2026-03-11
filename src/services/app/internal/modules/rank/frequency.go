package rank

import (
	"math"

	"github.com/lumi4x/news/src/services/app/client/models"
)

type FrequencyRanker struct {
}

func NewFrequencyRanker() *FrequencyRanker {
	return &FrequencyRanker{}
}

func (r *FrequencyRanker) Rank(item *models.RSSItem) int64 {
	return item.Rank
}

func (r *FrequencyRanker) ContextualizeRankAsc(currentIndex int, item *models.RSSItem, allItems []*models.RSSItem) int64 {
	nearbyItems := allItems //[max(0, currentIndex-2):min(len(allItems)-1, currentIndex+2)]

	rank := item.Rank
	derankedCount := 0
	for _, iterItem := range nearbyItems {
		if iterItem.Id == item.Id {
			continue
		}

		// Only derank if more articles of the same organization are nearby
		if iterItem.OrganizationUuid != item.OrganizationUuid {
			continue
		}

		// Derank
		rank -= int64(math.Abs(float64(rank)) * 0.1)
		derankedCount++
	}

	return rank
}

func (r *FrequencyRanker) ContextualizeRankDesc(currentIndex int, item *models.RSSItem, allItems []*models.RSSItem) int64 {
	return item.Rank
}
