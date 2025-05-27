.PHONY: seed
seed:
	@go run cmd/migrate/seed/seed.go
