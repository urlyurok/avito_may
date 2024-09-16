package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
)

type BidService struct {
	bidRepo         repo.Bid
	userRepo        repo.User
	responsibleRepo repo.Responsible
	tenderRepo      repo.Tender
}

func NewBidService(bidRepo repo.Bid, userRepo repo.User, responsibleRepo repo.Responsible, tenderRepo repo.Tender) *BidService {
	return &BidService{
		bidRepo:         bidRepo,
		userRepo:        userRepo,
		responsibleRepo: responsibleRepo,
		tenderRepo:      tenderRepo,
	}
}

func (bs *BidService) CreateBid(bi CreateBidInput) (entity.Bid, error) {

	user, err := bs.userRepo.GetById(context.Background(), bi.AuthorId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, bi.TenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.tenderRepo.IsTenderExists(context.Background(), bi.TenderId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	bid, err := bs.bidRepo.CreateBid(context.Background(), bi.Name, bi.Description, bi.TenderId, bi.AuthorType, bi.AuthorId)
	if err != nil {
		return entity.Bid{}, err
	}

	return bid, nil
}

func (bs *BidService) GetUserBids(ubp GetUserBidParams) ([]entity.Bid, error) {
	_, err := bs.userRepo.GetByName(context.Background(), ubp.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Bid{}, ErrUserNotExists
		}
		return []entity.Bid{}, err
	}

	return bs.bidRepo.GetUserBids(context.Background(), ubp.Username, ubp.Limit, ubp.Offset)
}

func (bs *BidService) GetBidsForTender(bftp GetBidsForTenderParams, tenderId string) ([]entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), bftp.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Bid{}, ErrUserNotExists
		}
		return []entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Bid{}, ErrUserIsNotResposible
		}
		return []entity.Bid{}, err
	}
	if !isResponsibe {
		return []entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return []entity.Bid{}, ErrTenderNotFound
	}

	return bs.bidRepo.GetBidsForTender(context.Background(), tenderId, bftp.Limit, bftp.Offset)
}

// TODO: aDD IS RESPONSEL BY BID ID
func (bs *BidService) GetBidStatus(u UserParam, bidId string) (string, error) {
	user, err := bs.userRepo.GetByName(context.Background(), u.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserNotExists
		}
		return "", err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserIsNotResposible
		}
		return "", err
	}
	if !isResponsibe {
		return "", ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return "", ErrTenderNotFound
	}

	return bs.bidRepo.GetBidStatus(context.Background(), bidId)
}

func (bs *BidService) UpdateBidStatus(params UpdateBidStatusParams, bidId string) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	return bs.bidRepo.UpdateBidStatus(context.Background(), params.Status, bidId)
}

func (bs *BidService) EditBid(up UserParam, bidId string, params map[string]interface{}) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), up.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}
	return bs.bidRepo.EditBid(context.Background(), bidId, params)
}

func (bs *BidService) SubmitBidFeedback(params SubmitBidFeedbackParams, bidId string) (entity.Bid, error) {

	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	err = bs.bidRepo.CreateBidFeedback(context.Background(), params.BidFeedback, bidId)
	if err != nil {
		return entity.Bid{}, err
	}

	return bs.bidRepo.GetBid(context.Background(), bidId)

}

func (bs *BidService) RollbackBid(up UserParam, bidId string, version int) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), up.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	bid, err := bs.bidRepo.RollbackBidVersion(context.Background(), bidId, version)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrBidOrVersionNotFound
		}
		return entity.Bid{}, err
	}

	return bid, err
}
