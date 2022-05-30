package main

type Model struct {
	Username       string  `json:"username"`
	Photo          string  `json:"photo"`
	Biography      string  `json:"biography"`
	Followers      int64   `json:"followers"`
	Following      int64   `json:"following"`
	AmountPosts    int64   `json:"amount_posts"`
	EngagementRate float64 `json:"engagement_rate"`
	AvgLikes       int64   `json:"avg_likes"`
	AvgComments    int64   `json:"avg_comments"`
	Posts          Posts   `json:"posts"`
}

type Post struct {
	Photo    string `json:"photo"`
	Likes    int64  `json:"likes"`
	Comments int64  `json:"comments"`
}

type Posts []Post

func (p Posts) AvgLikes() int64 {
	likes := int64(0)
	for _, v := range p {
		likes += v.Likes
	}

	return likes / int64(len(p))
}

func (p Posts) AvgComments() int64 {
	comments := int64(0)
	for _, v := range p {
		comments += v.Comments
	}

	return comments / int64(len(p))
}
