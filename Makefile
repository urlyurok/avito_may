include .env
export

migrate-up: ### migration up
	migrate -path /migrations -database '$(POSTGRES_CONN)?sslmode=disable' up
.PHONY: migrate-up