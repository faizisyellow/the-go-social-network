.PHONY:migration-up
migrate-up:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:F%40%21%24%24%40l343477s@/thegosocialnetwork up

.PHONY:migration-down
migrate-down:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:F%40%21%24%24%40l343477s@/thegosocialnetwork down


.PHONY: seed
seed:
	@go run cmd/migrate/seed/seed.go
