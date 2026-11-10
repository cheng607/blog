package model

import "time"

type Tag struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Count       int        `json:"count"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type Article struct {
	ID          int        `json:"id,omitempty"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CoverImage  *string    `json:"coverImage,omitempty"`
	PublishTime *time.Time `json:"publishTime,omitempty"`
	UpdateTime  *time.Time `json:"updateTime,omitempty"`
	UserID      int        `json:"userId"`
	ViewsCount  *int       `json:"viewsCount,omitempty"`
	TagIDs      []int      `json:"tagIds,omitempty"`
	Tags        []Tag      `json:"tags,omitempty"`
}

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Avatar    string     `json:"avatar"`
	Signature string     `json:"signature"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type Hollow struct {
	ID   int        `json:"id,omitempty"`
	Word string     `json:"word"`
	Date *time.Time `json:"date"`
}

type ArticlePayload struct {
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	CoverImage *string `json:"coverImage"`
	UserID     int     `json:"userId"`
	TagIDs     []int   `json:"tagIds"`
	Targets    []int   `json:"targets"`
}

func (p *ArticlePayload) ResolvedTagIDs() []int {
	if len(p.TagIDs) > 0 {
		return p.TagIDs
	}
	return p.Targets
}
