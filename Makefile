DB_DSN = host=localhost port=5432 user=root password=root dbname=tasks sslmode=disable
MIGRATIONS_DIR = migrations
BIN_DIR = bin
CLI_APP = taskmanager
WORKER_APP = taskcleaner

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

db-connect:
	docker exec -it techno-db psql -U root -d tasks

db-tables:
	docker exec -it techno-db psql -U root -d tasks -c "\dt"

db-describe:
	docker exec -it techno-db psql -U root -d tasks -c "\d tasks"

up:
	docker-compose up -d

down:
	docker-compose down

build: build-cli build-worker

build-cli:
	go build -o $(BIN_DIR)/$(CLI_APP) ./cmd/app

build-worker:
	go build -o $(BIN_DIR)/$(WORKER_APP) ./cmd/worker

clean:
	rm -rf $(BIN_DIR)