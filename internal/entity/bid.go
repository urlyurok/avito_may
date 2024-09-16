package entity

import "time"

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderId    string    `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorId    string    `json:"authorId"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

type VersionedBid struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      string `db:"status"`
	TenderId    string `db:"tender_id"`
	AuthorType  string `db:"author_type"`
	AuthorId    string `db:"author_id"`
}
