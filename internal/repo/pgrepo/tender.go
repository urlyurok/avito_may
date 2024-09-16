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

type TenderRepo struct {
	*postgres.Postgres
}

func NewTenderRepo(pg *postgres.Postgres) *TenderRepo {
	return &TenderRepo{pg}
}

func (r *TenderRepo) CreateTender(ctx context.Context, name, description, serviceType, status, organizationId string) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.CreateTender"

	sql := `
	INSERT INTO tender (name, description, service_type, status, organization_id)
	VALUES ($1, $2, UPPER($3)::service_type, UPPER($4)::tender_status, $5) 
	RETURNING id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
	`

	var t entity.Tender
	err := r.Pool.QueryRow(ctx, sql, name, description, serviceType, status, organizationId).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		log.Debug("err: ", err)
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	log.Info("CreateTender tender: ", "tender", t)

	return t, nil

}

func (r *TenderRepo) GetTenders(ctx context.Context, limit, offset int, serviceType []string) ([]entity.Tender, error) {
	const fn = "repo.pgrepo.tender.GetTenders"

	var rows pgx.Rows
	var err error

	if len(serviceType) == 0 {
		sql := `
		SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
		FROM tender
		LIMIT $1
		OFFSET $2
		`
		rows, err = r.Pool.Query(ctx, sql, limit, offset)
	} else {
		sql := `
		SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
		FROM tender
		WHERE service_type::text = ANY($1)
		LIMIT $2
		OFFSET $3
		`
		rows, err = r.Pool.Query(ctx, sql, serviceType, limit, offset)

	}

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Tender{}, repoerrs.ErrNotFound
		}
		return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var t entity.Tender
		err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.ServiceType,
			&t.Status,
			&t.OrganizationId,
			&t.Version,
			&t.CreatedAt,
		)
		if err != nil {
			return []entity.Tender{}, fmt.Errorf("%s: %v", err)
		}
		tenders = append(tenders, t)
	}

	log.Debug("GetTenders: ", tenders)

	return tenders, nil
}

func (r *TenderRepo) GetUserTenders(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error) {
	const fn = "repo.pgrepo.tender.GetUserTenders"

	sql := `
		SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
		FROM tender
		WHERE organization_id in (
			SELECT o.id
			FROM organization_responsible ores
					 JOIN organization o ON ores.organization_id = o.id
					 JOIN employee e ON ores.user_id = e.id
			WHERE e.username = $1)
		LIMIT $2
		OFFSET $3
		`

	rows, err := r.Pool.Query(ctx, sql, username, limit, offset)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Tender{}, repoerrs.ErrNotFound
		}
		return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var t entity.Tender
		err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.ServiceType,
			&t.Status,
			&t.OrganizationId,
			&t.Version,
			&t.CreatedAt,
		)
		if err != nil {
			return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
		}
		tenders = append(tenders, t)
	}

	log.Debug("GetUserTenders: ", tenders)

	return tenders, nil
}

func (r *TenderRepo) GetTenderStatus(ctx context.Context, tenderId string) (string, error) {
	const fn = "repo.pgrepo.tender.GetTenderStatus"

	sql := `
		SELECT INITCAP(status::text) AS status
		FROM tender
		WHERE id = $1
		`

	var status string
	err := r.Pool.QueryRow(ctx, sql, tenderId).Scan(&status)
	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		return "", fmt.Errorf("%s: %v", fn, err)
	}

	return status, nil
}

func (r *TenderRepo) UpdateTender(ctx context.Context, tenderId string, params map[string]interface{}) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.UpdateTender"

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, _ := builder.
		Update("tender").
		SetMap(params).
		Where("id = ?", tenderId).
		Suffix("RETURNING id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at").
		ToSql()

	var t entity.Tender
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", err.Error())
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	return t, nil

}

func (r *TenderRepo) UpdateTenderStatus(ctx context.Context, status, tenderId string) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.UpdateTenderStatus"

	sql := `
		UPDATE tender
		SET status = UPPER($1)::tender_status
		WHERE id = $2
		RETURNING id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at`

	var t entity.Tender
	err := r.Pool.QueryRow(ctx, sql, status, tenderId).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		log.Debug("err: ", err)
		if err == pgx.ErrNoRows {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}
	log.Debug("Upadted tender: ", t)

	return t, nil
}

func (r *TenderRepo) GetTenderById(ctx context.Context, tenderId string) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.GetTenderById"

	sql := `
	SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
	FROM tender
	WHERE id = $1
	`

	var t entity.Tender
	err := r.Pool.QueryRow(ctx, sql, tenderId).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", fn, err)
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}
	log.Debug("Upadted tender: ", t)

	return t, nil
}

func (r *TenderRepo) RollbackTenderVersion(ctx context.Context, tenderId string, version int) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.RollbackTenderVersion"

	sql := `
	SELECT name, description, service_type, status, organization_id 
	FROM tender_versions
	WHERE tender_id = $1 AND version = $2
    `

	var vc entity.VersionedTender
	err := r.Pool.QueryRow(ctx, sql, tenderId, version).Scan(
		&vc.Name,
		&vc.Description,
		&vc.ServiceType,
		&vc.Status,
		&vc.OrganizationId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", fn, err)
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	// Увеличение версии и сохранение как новой
	sql = `
        UPDATE tender 
        SET name = $1, description = $2, service_type = $3, 
            status = $4, organization_id = $5
        WHERE id = $6
        RETURNING id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
        `

	var t entity.Tender
	err = r.Pool.QueryRow(ctx, sql, vc.Name, vc.Description, vc.ServiceType, vc.Status, vc.OrganizationId, tenderId).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		log.Debug("err: ", fn, err)
		return entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	return t, nil
}

func (r *TenderRepo) IsTenderExists(ctx context.Context, tenderId string) (bool, error) {
	const fn = "repo.pgrepo.tender.CheckTenderExists"

	sql := `SELECT EXISTS (SELECT 1 FROM tender WHERE id = $1)`

	var exists bool
	err := r.Pool.QueryRow(ctx, sql, tenderId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %v", fn, err)
	}

	return exists, nil
}
