package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type ResponsibleRepo struct {
	*postgres.Postgres
}

func NewResponsibleRepo(pg *postgres.Postgres) *ResponsibleRepo {
	return &ResponsibleRepo{pg}
}

func (r *ResponsibleRepo) GetAllResponsiblesByUserId(ctx context.Context, userId string) ([]entity.Responsible, error) {
	const fn = "repo.pgrepo.responsible.GetAllResponsiblesByUserId"

	sql := `
	SELECT *
	FROM organization_responsible
	WHERE user_id=$1::uuid
	`

	rows, err := r.Pool.Query(ctx, sql, userId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Responsible{}, repoerrs.ErrNotFound
		}
		return []entity.Responsible{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var responsibles []entity.Responsible
	for rows.Next() {
		var responsible entity.Responsible
		err := rows.Scan(
			&responsible.Id,
			&responsible.OrganizationId,
			&responsible.UserId,
		)
		if err != nil {
			return []entity.Responsible{}, fmt.Errorf("%s: %v", fn, err)
		}
		responsibles = append(responsibles, responsible)
	}

	return responsibles, nil
}

func (r *ResponsibleRepo) IsUserResponsibleForOrganizationByOrganizationId(ctx context.Context, userId, organizationId string) (bool, error) {
	const fn = "repo.pgrepo.responsible.IsUserResponsibleForOrganizationByOrganizationId"

	responsibles, err := r.GetAllResponsiblesByUserId(ctx, userId)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return false, err
		}
		return false, fmt.Errorf("%s: %v", fn, err)
	}

	for _, responsible := range responsibles {
		if responsible.OrganizationId == organizationId {
			return true, nil
		}
	}

	return false, nil
}

func (r *ResponsibleRepo) IsUserResponsibleForOrganizationByTenderId(ctx context.Context, userId, tenderId string) (bool, error) {
	const fn = "repo.pgrepo.responsible.IsUserResponsibleForOrganizationByTenderId"

	sql := `
	SELECT 
    org_res.user_id
	FROM
		tender t
	JOIN
		organization org ON t.organization_id = org.id
	JOIN
		organization_responsible org_res ON org.id = org_res.organization_id
	WHERE
		t.id = $1
		AND org_res.user_id = $2
	`

	var id string
	err := r.Pool.QueryRow(ctx, sql, tenderId, userId).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, repoerrs.ErrNotFound
		}
		return false, fmt.Errorf("%s: %v", fn, err)
	}

	return true, nil
}

func (r *ResponsibleRepo) IsUserResponsibleForOrganizationByBidId(ctx context.Context, userId, bidId string) (bool, error) {
	const fn = "repo.pgrepo.responsible.IsUserResponsibleForOrganizationByBidId"

	sql := `

	`
	_ = sql

	return true, nil

}
