include .env

stop_containers:
	@echo "Stopping other docker container"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "no containers running..."; \
	fi

create_container:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${DB_USERNAME} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:12-alpine

create_db:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USERNAME} --owner=${DB_USERNAME} ${DB_NAME}

start_container:
	docker start ${DB_DOCKER_CONTAINER}

create_migrations:
	sqlx migrate add -r init

migrate_up:
	sqlx migrate run --database-url "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

migrate_down:
	sqlx migrate revert --database-url "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

build:
	@echo "Building binary..."
	go build -o ${BINARY_NAME} cmd/server/*.go
	@echo "Binary built!"

run: build stop_containers start_container
	@echo "Startin api"
	@env SERVER_PORT=${SERVER_PORT} DSN=${DSN} ./${BINARY_NAME} &
	@echo "api started!"
# @echo "api started..."

stop:
	@echo "stopping server.."
	@-pkill -SIGTERM -f "./${BINARY}"
	@echo "server stopped..."

start: run

restart: stop start