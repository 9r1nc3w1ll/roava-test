MIGRATE_PATH ?= "migrations"
MIGRATE_DB ?= "postgres://postgres:@db/roava_test?sslmode=disable"

proto_generate:
	mkdir -p pb && protoc --proto_path=proto --go_out=plugins=grpc:. proto/*.proto

proto_clean:
	rm pb/*.pb.go

migrate: migrate-up

migrate-fresh: migrate-drop migrate-up

migrate-up:
	migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) down

migrate-drop:
	migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) drop

migrate-version:
	migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) version
