MIGRATE_PATH = cmd/migrate/migrations

.PHONY: new-migration
migrate:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

.PHONY:migration-up
migrate-up:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:F%40%21%24%24%40l343477s@/thegosocialnetwork up

.PHONY:migration-down
migrate-down:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:F%40%21%24%24%40l343477s@/thegosocialnetwork down

.PHONY:migration-back
migrate-back:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:F%40%21%24%24%40l343477s@/thegosocialnetwork force $(no)


.PHONY: seed
seed:
	@go run cmd/migrate/seed/seed.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
