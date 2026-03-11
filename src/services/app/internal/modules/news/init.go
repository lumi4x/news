package news

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lumi4x/news/src/services/app/client/models"
	"github.com/lumi4x/news/src/services/app/internal/modules/organization"
	"github.com/lumi4x/news/src/services/app/internal/modules/rssfeed"
	"github.com/pkg/errors"
)

const (
	whitelistRegex = `^.*(palestin|israel|gaza|west bank|netanyahu|hamas).*$`
)

var (
	rankTierBoosted     = float32(7 * time.Hour.Seconds())
	rankTierSubstandard = float32(-2 * time.Hour.Seconds())
	rankTierLow         = float32(-7 * time.Hour.Seconds())
	rankTierSuppressed  = float32(-18 * time.Hour.Seconds())
)

var orgs = []struct {
	*models.Organization
	RssFeedUrl                  string
	RssFeedRank                 float32
	ContentWhitelistRegex       string
	HtmlPaywallRegex            string
	TitleOrDescriptionTrimRight string
	MaxFetchIntervalMin         int64
	DisacrdOGImage              bool
}{
	{
		Organization: &models.Organization{
			Uuid:         "f97d83e7-2f08-4f8c-9ec8-98279851608c",
			UniqueName:   "972mag",
			Name:         "+972 Magazine",
			Description:  "Independent commentary and news from Israel &#38; Palestine",
			WebsiteUrl:   "https://www.972mag.com",
			LogoImageUrl: "/organizations/972mag_logo.svg",
		},
		RssFeedUrl:  "https://www.972mag.com/feed",
		RssFeedRank: rankTierBoosted,
	},

	{
		Organization: &models.Organization{
			Uuid:        "278ae4c8-b190-4543-93ec-4e28d558633e",
			UniqueName:  "btselem",
			Name:        "Btselem",
			Description: "Independent commentary and news from Israel & Palestine",
			WebsiteUrl:  "https://www.btselem.org",
		},
		// RssFeedUrl: "https://www.btselem.org/rss",
	},

	{
		Organization: &models.Organization{
			Uuid:         "ad470590-8199-4906-b189-782acb7e672d",
			UniqueName:   "mondoweiss",
			Name:         "Mondoweiss",
			Description:  "News & Opinion about Palestine, Israel & the United States",
			WebsiteUrl:   "https://mondoweiss.net",
			LogoImageUrl: "/organizations/mondoweiss_logo.svg",
		},
		RssFeedUrl: "https://mondoweiss.net/rss",
	},

	{
		Organization: &models.Organization{
			Uuid:         "3b613eb0-2d13-4a74-bd94-4acb57c8fef9",
			UniqueName:   "bds",
			Name:         "BDS",
			Description:  "Boycott, Divestment, Sanctions (BDS) is a Palestinian-led movement for freedom, justice and equality.",
			WebsiteUrl:   "https://bdsmovement.net",
			LogoImageUrl: "/organizations/bds_logo.png",
		},
		//RssFeedUrl: "https://bdsmovement.net/rss-feed.xml",
		RssFeedUrl:     "https://news.google.com/rss/search?q=allinurl:bdsmovement.net&hl=en-IL&gl=IL&ceid=IL:en",
		RssFeedRank:    rankTierSubstandard,
		DisacrdOGImage: true,
	},

	{
		Organization: &models.Organization{
			Uuid:         "3d4903ec-416d-48dc-99d2-8ef83729a5c9",
			UniqueName:   "wearenotnumbers",
			Name:         "We Are Not Numbers",
			Description:  "WANN, a project of Euro-Med Human rights Monitor, pairs Palestinian writers with international mentors to write their stories behind the numbers. Content is not censored by Euro-Med Monitor.",
			WebsiteUrl:   "https://wearenotnumbers.org",
			LogoImageUrl: "/organizations/wearenotnumbers_logo.png",
		},
		RssFeedRank:         rankTierSubstandard,
		RssFeedUrl:          "https://wearenotnumbers.org/category/story/feed",
		MaxFetchIntervalMin: 60,
	},

	{
		Organization: &models.Organization{
			Uuid:        "419b3510-b226-45c6-9dc6-2cf06acdeecc",
			UniqueName:  "mistaclim",
			Name:        "Mistaclim",
			Description: "A group of normative Israeli citizens, people who care about human rights and the future of the state. Mistaclim opposes the occupation and acts to put an end to it.",
			WebsiteUrl:  "https://www.mistaclim.org",
		},
	},

	{
		Organization: &models.Organization{
			Uuid:         "5707d6e6-4998-43ab-9ea3-da94c6bddd64",
			UniqueName:   "zeteo",
			Name:         "Zeteo",
			Description:  "Zeteo is a new media organization that seeks answers for the questions that really matter, while always striving for the truth. Founded by Mehdi Hasan, Zeteo is a movement for media accountability, unfiltered news and bold opinions.",
			WebsiteUrl:   "https://zeteo.com",
			LogoImageUrl: "/organizations/zeteo_logo.png",
		},
		RssFeedUrl:            "https://zeteo.com/feed",
		ContentWhitelistRegex: whitelistRegex,
		HtmlPaywallRegex:      `(?m)^.*audience\\":\\"only_paid.*$`,
	},

	{
		Organization: &models.Organization{
			Uuid:         "ba603da3-1292-479d-9e0c-33f61a86c0fa",
			UniqueName:   "aljazeera",
			Name:         "Al Jazeera",
			Description:  "Breaking News, World News and Video from Al Jazeera",
			WebsiteUrl:   "https://www.aljazeera.com",
			LogoImageUrl: "/organizations/aljazeera_logo.svg",
		},

		RssFeedUrl:            "https://www.aljazeera.com/xml/rss/all.xml",
		RssFeedRank:           rankTierSuppressed,
		ContentWhitelistRegex: whitelistRegex,
	},

	{
		Organization: &models.Organization{
			Uuid:         "d422fde4-8208-43d4-9260-9ea2a89e0a80",
			UniqueName:   "associatedpress",
			Name:         "Associated Press",
			Description:  "Read the latest headlines, breaking news, and videos at APNews.com, the definitive source for independent journalism from every corner of the globe.",
			WebsiteUrl:   "https://apnews.com",
			LogoImageUrl: "/organizations/associatedpress_logo.svg",
		},
		RssFeedUrl:                  "https://news.google.com/rss/search?q=allinurl:apnews.com&hl=en-IL&gl=IL&ceid=IL:en",
		RssFeedRank:                 rankTierSuppressed,
		ContentWhitelistRegex:       whitelistRegex,
		TitleOrDescriptionTrimRight: " - The Associated Press",
		DisacrdOGImage:              true,
	},

	{
		Organization: &models.Organization{
			Uuid:         "88956b49-7194-4cdb-915b-fb52bbda20f5",
			UniqueName:   "newarab",
			Name:         "New Arab",
			Description:  "The New Arab is a leading English-language news website bringing you the big stories from the Middle East, North Africa and beyond.",
			WebsiteUrl:   "https://www.newarab.com",
			LogoImageUrl: "/organizations/newarab_logo.png",
		},
		RssFeedUrl:            "https://www.newarab.com/rss",
		RssFeedRank:           rankTierSuppressed,
		ContentWhitelistRegex: whitelistRegex,
	},
}

func Initialize(ctx context.Context) error {
	for _, org := range orgs {
		orgUuid := uuid.MustParse(org.Uuid)

		err := organization.UpsertOrganization(ctx, orgUuid, org.UniqueName, org.Name, org.Description, org.WebsiteUrl, org.LogoImageUrl)
		if err != nil {
			return errors.Wrapf(err, "failed upserting org with uuid '%s'", org.Uuid)
		}

		if org.RssFeedUrl != "" {
			err = rssfeed.Upsert(ctx, orgUuid, org.RssFeedUrl, org.ContentWhitelistRegex, org.HtmlPaywallRegex, org.TitleOrDescriptionTrimRight, org.RssFeedRank, org.MaxFetchIntervalMin, org.DisacrdOGImage)
			if err != nil {
				return errors.Wrapf(err, "failed upserting RSS feed")
			}
		}
	}

	return nil
}
