COMPOSE_PROJECT_NAME=nanobank
COMPOSE_FILE=./deployments/docker-compose.yml

.PHONY: tidy up down	
tidy-gateway:
	@cd ./gateway/ && go mod tidy
tidy-auth:
	@cd ./auth-service/ && go mod tidy
tidy-account:
	@cd ./account-service/ && go mod tidy
tidy-transaction:
	@cd ./transaction-service/ && go mod tidy
tidy-shared:
	@cd ./shared/ && go mod tidy


tidy: tidy-shared tidy-auth tidy-gateway tidy-account tidy-transaction

up:
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) up -d --build

down:
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) down

reset: down up
