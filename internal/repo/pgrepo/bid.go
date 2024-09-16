package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	log "log/slog"
)

type BidRepo struct {
	*postgres.Postgres
}

func NewBidRepo(pg *postgres.Postgres) *BidRepo {
	return &BidRepo{pg}
}

func (r *BidRepo) CreateBid(ctx context.Context, name string, description string, tenderId string, authorType string, authorId string) (entity.Bid, error) {
	sql := `
	INSERT INTO bid (name, description, tender_id, author_type, author_id)
	VALUES ($1, $2, $3, UPPER($4)::authore_type, $5) 
	RETURNING id, name, description, INITCAP(status::text), tender_id, INITCAP(author_type::text), author_id, version, created_at
	`

	var b entity.Bid
	err := r.Pool.QueryRow(ctx, sql, name, description, tenderId, authorType, authorId).Scan(
		&b.Id,
		&b.Name,
		&b.Description,
		&b.Status,
		&b.TenderId,
		&b.AuthorType,
		&b.AuthorId,
		&b.Version,
		&b.CreatedAt,
	)
	if err != nil {
		return entity.Bid{}, repoerrs.ErrUnableToInsert
	}

	return b, nil
}

func (r *BidRepo) GetUserBids(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error) {
	const fn = "repo.pgrepo.bid.GetUserBids"

	sql := `
	SELECT b.id, b.name, b.description, INITCAP(b.status::text), b.tender_id, INITCAP(b.author_type::text), b.author_id, b.version, b.created_at
	FROM bid b
			 JOIN employee e ON b.author_id = e.id
	WHERE e.username = $1
	LIMIT $2
	OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, sql, username, limit, offset)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Bid{}, repoerrs.ErrNotFound
		}
		return []entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var bids []entity.Bid
	for rows.Next() {
		var b entity.Bid
		err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Description,
			&b.Status,
			&b.TenderId,
			&b.AuthorType,
			&b.AuthorId,
			&b.Version,
			&b.CreatedAt,
		)
		if err != nil {
			return []entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
		}
		bids = append(bids, b)
	}
	return bids, nil

}

func (r *BidRepo) GetBidsForTender(ctx context.Context, tenderId string, limit int, offset int) ([]entity.Bid, error) {
	const fn = "repo.pgrepo.bid.GetBidsForTender"

	sql := `
	SELECT b.id, b.name, b.description, INITCAP(b.status::text), b.tender_id, INITCAP(b.author_type::text), b.author_id, b.version, b.created_at
	FROM bid b
	WHERE tender_id = $1
	LIMIT $2
	OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, sql, tenderId, limit, offset)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Bid{}, repoerrs.ErrNotFound
		}
		return []entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var bids []entity.Bid
	for rows.Next() {
		var b entity.Bid
		err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Description,
			&b.Status,
			&b.TenderId,
			&b.AuthorType,
			&b.AuthorId,
			&b.Version,
			&b.CreatedAt,
		)
		if err != nil {
			return []entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
		}
		bids = append(bids, b)
	}
	return bids, nil

}

func (r *BidRepo) GetBidStatus(ctx context.Context, bidId string) (string, error) {
	const fn = "repo.pgrepo.bid.GetBidStatus"

	sql := `
	SELECT INITCAP(status::text) as status
	FROM bid
	WHERE id = $1
	`

	var status string
	err := r.Pool.QueryRow(ctx, sql, bidId).Scan(&status)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		return "", fmt.Errorf("%s: %v", fn, err)
	}

	return status, nil
}

func (r *BidRepo) UpdateBidStatus(ctx context.Context, status string, bidId string) (entity.Bid, error) {
	const fn = "repo.pgrepo.bid.UpdateBidStatus"

	sql := `
	UPDATE bid
	SET status = UPPER($1)::bid_status
	WHERE id = $2
	RETURNING id, name, description, status, tender_id, INITCAP(author_type::text), author_id, version, created_at
	`

	var b entity.Bid
	err := r.Pool.QueryRow(ctx, sql, status, bidId).Scan(
		&b.Id,
		&b.Name,
		&b.Description,
		&b.Status,
		&b.TenderId,
		&b.AuthorType,
		&b.AuthorId,
		&b.Version,
		&b.CreatedAt,
	)
	if err != nil {
		log.Debug("err: ", err)
		if err == pgx.ErrNoRows {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	return b, nil

}

func (r *BidRepo) EditBid(ctx context.Context, bidId string, params map[string]interface{}) (entity.Bid, error) {
	const fn = "repo.pgrepo.bid.EditBid"

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, _ := builder.
		Update("bid").
		SetMap(params).
		Where("id = ?", bidId).
		Suffix("RETURNING id, name, description, status, tender_id, INITCAP(author_type::text), author_id, version, created_at").
		ToSql()

	var b entity.Bid
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&b.Id,
		&b.Name,
		&b.Description,
		&b.Status,
		&b.TenderId,
		&b.AuthorType,
		&b.AuthorId,
		&b.Version,
		&b.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", err.Error())
		return entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	return b, nil
}

func (r *BidRepo) CreateBidFeedback(ctx context.Context, feedback string, bidId string) error {
	const fn = "repo.pgrepo.bid.CreateBidFeedback"

	sql := `
	INSERT INTO bid_review
	(description)
	VALUES
	($1)
	RETURNING id
	`

	var reviewId string
	err := r.Pool.QueryRow(ctx, sql, feedback).Scan(
		&reviewId,
	)
	if err != nil {
		return repoerrs.ErrUnableToInsert
	}

	sql = `
	INSERT INTO bid_bidreview
	(bid_id, bid_review_id)
	VALUES
	($1, $2)
	`
	if _, err = r.Pool.Exec(ctx, sql, bidId, reviewId); err != nil {
		return repoerrs.ErrUnableToInsert

	}

	return nil

}

func (r *BidRepo) GetBid(ctx context.Context, id string) (entity.Bid, error) {
	const fn = "repo.pgrepo.bid.GetBid"

	sql := `
	SELECT id, name, description, status, tender_id, INITCAP(author_type::text), author_id, version, created_at
	FROM bid
	WHERE id = $1
	`

	var b entity.Bid
	err := r.Pool.QueryRow(ctx, sql, id).Scan(
		&b.Id,
		&b.Name,
		&b.Description,
		&b.Status,
		&b.TenderId,
		&b.AuthorType,
		&b.AuthorId,
		&b.Version,
		&b.CreatedAt,
	)
	if err != nil {
		log.Debug("err: ", err)
		if err == pgx.ErrNoRows {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}
	return b, nil
}

func (r *BidRepo) RollbackBidVersion(ctx context.Context, bidId string, version int) (entity.Bid, error) {
	const fn = "repo.pgrepo.bid.RollbackBidVersion"

	sql := `
	SELECT name, description, status, tender_id, author_type, author_id 
	FROM bid_versions
	WHERE bid_id = $1 AND version = $2
    `

	var vb entity.VersionedBid
	err := r.Pool.QueryRow(ctx, sql, bidId, version).Scan(
		&vb.Name,
		&vb.Description,
		&vb.Status,
		&vb.TenderId,
		&vb.AuthorType,
		&vb.AuthorId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", fn, err)
		return entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	sql = `
	UPDATE bid
	SET name = $1, description = $2, status = $3, tender_id = $4, author_type = $5, author_id = $6
	WHERE id = $7
	RETURNING id, name, description, INITCAP(status::text), tender_id, INITCAP(author_type::text), author_id, version, created_at
	`

	var b entity.Bid
	err = r.Pool.QueryRow(ctx, sql, vb.Name, vb.Description, vb.Status, vb.TenderId, vb.AuthorType, vb.AuthorId, bidId).Scan(
		&b.Id,
		&b.Name,
		&b.Description,
		&b.Status,
		&b.TenderId,
		&b.AuthorType,
		&b.AuthorId,
		&b.Version,
		&b.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", fn, err)
		return entity.Bid{}, fmt.Errorf("%s: %v", fn, err)
	}

	return b, nil
}

func (r *BidRepo) IsBidExists(ctx context.Context, bidId string) (bool, error) {
	const fn = "repo.pgrepo.bid.IsBidExists"

	sql := `SELECT EXISTS (SELECT 1 FROM bid WHERE id = $1)`

	var exists bool
	err := r.Pool.QueryRow(ctx, sql, bidId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %v", fn, err)
	}

	return exists, nil

}
