package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
)

type CreateTenderInput struct {
	Name            string `json:"name" validate:"required,min=3"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"serviceType" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  string `json:"organizationId" validate:"required"`
	CreatorUsername string `json:"creatorUsername" validate:"required"`
}

type PatchTenderInput struct {
	Name            string `json:"name" validate:"omitempty,min=3"`
	Description     string `json:"description" validate:"omitempty"`
	ServiceType     string `json:"serviceType" validate:"omitempty"`
	Status          string `json:"status" validate:"omitempty"`
	OrganizationId  string `json:"organizationId" validate:"omitempty"`
	CreatorUsername string `json:"creatorUsername" validate:"omitempty"`
}

type GetParams struct {
	Limit  int `schema:"limit" validate:"required,gt=0"`
	Offset int `schema:"offset" validate:"omitempty,gte=0"`
}

type UserParam struct {
	Username string `schema:"username" validate:"required"`
}

type GetTendersParams struct {
	GetParams
	ServiceType []string `schema:"service_type" validate:"required,dive,alpha"`
}

type GetUserTendersParams struct {
	GetParams
	UserParam
}

type UpdateTenderStatusParams struct {
	UserParam
	Status string `json:"status" validate:"required"`
}

type CreateBidInput struct {
	Name        string `json:"name" validate:"required,alpha"`
	Description string `json:"description" validate:"required"`
	TenderId    string `json:"tenderId" validate:"required"`
	AuthorType  string `json:"authorType" validate:"required,alpha"`
	AuthorId    string `json:"authorId" validate:"required"`
}

type GetUserBidParams struct {
	GetParams
	UserParam
}

type GetBidsForTenderParams struct {
	GetParams
	UserParam
}

type UpdateBidStatusParams struct {
	UserParam
	Status string `json:"status" validate:"required"`
}

type SubmitBidFeedbackParams struct {
	UserParam
	BidFeedback string `json:"bidFeedback" validate:"required"`
}

type Tender interface {
	CreateTender(input CreateTenderInput) (entity.Tender, error)
	GetTenders(gtp GetTendersParams) ([]entity.Tender, error)
	GetUserTenders(gutp GetUserTendersParams) ([]entity.Tender, error)
	GetTenderStatus(up UserParam, tenderId string) (string, error)
	EditTender(up UserParam, tenderId string, params map[string]interface{}) (entity.Tender, error)
	UpdateTenderStatus(utsp UpdateTenderStatusParams, id string) (entity.Tender, error)
	RollbackTender(u UserParam, tenderId string, version int) (entity.Tender, error)
}

type Bid interface {
	CreateBid(input CreateBidInput) (entity.Bid, error)
	GetUserBids(params GetUserBidParams) ([]entity.Bid, error)
	GetBidsForTender(params GetBidsForTenderParams, tenderId string) ([]entity.Bid, error)
	GetBidStatus(param UserParam, bidId string) (string, error)
	UpdateBidStatus(params UpdateBidStatusParams, bidId string) (entity.Bid, error)
	EditBid(param UserParam, bidId string, params map[string]interface{}) (entity.Bid, error)
	SubmitBidFeedback(params SubmitBidFeedbackParams, bidId string) (entity.Bid, error)
	RollbackBid(param UserParam, bidId string, version int) (entity.Bid, error)
}

type Services struct {
	Tender Tender
	Bid    Bid
}

func NewServices(r *repo.Repositories) *Services {
	return &Services{
		Tender: NewTenderService(r.Tender, r.User, r.Responsible),
		Bid:    NewBidService(r.Bid, r.User, r.Responsible, r.Tender),
	}
}
