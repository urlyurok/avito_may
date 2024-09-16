package entity

import (
	"strings"
	"time"
)

type Tender struct {
	Id             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	ServiceType    string    `db:"service_type" json:"serviceType"`
	Status         string    `db:"status" json:"status"`
	OrganizationId string    `db:"organization_id" json:"organizationId"`
	Version        int       `db:"version" json:"version"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

type VersionedTender struct {
	Name           string `db:"name"`
	Description    string `db:"description"`
	ServiceType    string `db:"service_type"`
	Status         string `db:"status"`
	OrganizationId string `db:"organization_id"`
}

func (t *Tender) Capitalize() {
	t.ServiceType = strings.Title(strings.ToLower(t.ServiceType))
	t.Status = strings.Title(strings.ToLower(t.Status))
}
