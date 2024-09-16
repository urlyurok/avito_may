package entity

import "time"

type BidReview struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
