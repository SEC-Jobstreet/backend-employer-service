DB_URL=postgresql://admin:admin@localhost:5432/employer_service_jobstreet?sslmode=disable

docker_compose_build:
	docker-compose build

docker_build:
	docker buildx build -t thanhquy1105/backend-jobstreet-employer-service-prod:latest .

docker_push:
	docker push thanhquy1105/backend-jobstreet-employer-service-prod

docker_build_run:
	docker-compose up

# generate a new migration
new_migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

# run postgres container with network 
run_postgres:
	-docker network create jobstreet-network
	docker run --name postgres --network jobstreet-network -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:13.12 

start_postgres:
	docker start postgres

build_app:
	docker build -t thanhquy1105/backend-jobstreet-employer-service-prod:latest .

run_app:
	docker run --name backend-jobstreet-employer-service-prod --network jobstreet-network -p 4001:4001 -e DB_SOURCE="postgresql://admin:admin@postgres:5432/employer_service_jobstreet?sslmode=disable" thanhquy1105/backend-jobstreet-employer-service-prod:latest

start_app:
	docker start backend-jobstreet-employer-service-prod

push_app:
	docker push thanhquy1105/backend-jobstreet-employer-service-prod

# create employer_service_jobstreet database on postgres container
createdb:
	docker exec -it postgres createdb --username=admin --owner=admin employer_service_jobstreet

# drop employer_service_jobstreet database on postgres container
dropdb:
	docker exec -it postgres dropdb --username=admin employer_service_jobstreet

# migrate employer_service_jobstreet database from app to postgres container
migrate:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

# generate queries to golang code
sqlc:
	docker run --rm -v "${CURDIR}:/src" -w /src sqlc/sqlc:1.20.0 generate

# generate swagger
swagger:
	swag init --parseDependency -g main.go

# run test
test:
	go test -v -cover -short ./...

.PHONY: build_run_prod new_migrate run_postgres migrate dropdb createdb start_postgres sqlc evans swagger proto