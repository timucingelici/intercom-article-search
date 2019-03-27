package helpdocs

// Helpdocs Components
type Search struct {
	Articles       []Article `json:"articles"`
	Query          string    `json:"query"`
	QueryTruncated bool      `json:"query_truncated"`
	TotalHits      int       `json:"totalHits"`
	MaxScore       float64   `json:"maxScore"`
	StemmedQuery   string    `json:"stemmedQuery"`
}

type Article struct {
	AccountId        string      `json:"account_id"`
	ArticleId        string      `json:"article_id"`
	Body             string      `json:"body"`
	Canonical        string      `json:"canonical"`
	CategoryTitle    string      `json:"category_title"`
	Description      string      `json:"description"`
	IsPrivate        bool        `json:"is_private"`
	IsPublished      bool        `json:"is_published"`
	PermissionGroups []string    `json:"permission_groups"`
	RelativeUrl      string      `json:"relative_url"`
	Search           SearchScore `json:"search"`
	ShortVersion     string      `json:"short_version"`
	Slug             string      `json:"slug"`
	Tags             []string    `json:"tags"`
	Title            string      `json:"title"`
	Url              string      `json:"url"`
}

type ArticleResponse struct {
	Article Article `json:"article"`
}

type SearchScore struct {
	score float64
}
