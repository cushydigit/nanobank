.PHONY: tidy up down
tidy-gateway:
	@cd ./gateway/ && go mod tidy
tidy-auth:
	@cd ./auth-service/ && go mod tidy
tidy-shared:
	@cd ./shared/ && go mod tidy
tidy: tidy-shared tidy-auth tidy-gateway

up:
	COMPOSE_PROJECT_NAME=nanabank docker-compose -f ./deployments/docker-compose.yml up -d --build

down:
	docker-compose -f ./deployments/docker-compose.yml down

reset: down up
