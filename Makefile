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
tidy-mailer:
	@cd ./mailer-service/ && go mod tidy
tidy-shared:
	@cd ./shared/ && go mod tidy

tidy: tidy-shared tidy-auth tidy-gateway tidy-account tidy-transaction tidy-mailer

up:
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) up -d --build

down:
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) down

reset: down tidy up

minikube_docker:
	@eval $(minikube docker-env)

minikube_docker_unset:
	@eval $(minikube docker-env --unset)

build_image_auth:
	@docker build -t auth-service:latest -f auth-service/Dockerfile .

build_image_gateway:
	@docker build -t gateway:latest -f gateway/Dockerfile .


