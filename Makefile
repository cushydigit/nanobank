.PHONY: tidy-auth tidy-shared tidy

tidy-auth:
	@cd ./auth-service/ && go mod tidy
tidy-shared:
	@cd ./shared/ && go mod tidy
tidy: tidy-shared tidy-auth

up:
	COMPOSE_PROJECT_NAME=nanabank docker-compose -f ./deployments/docker-compose.yml up -d --build

down:
	docker-compose -f ./deployments/docker-compose.yml down

reset: down up
