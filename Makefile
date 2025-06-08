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

build_mini_gateway:
	@eval $(minikube docker-env) && docker build -t gateway:latest -f gateway/Dockerfile .

build_mini_auth:
	@eval $(minikube docker-env) && docker build -t auth-service:latest -f auth-service/Dockerfile .

build_mini_account:
	@eval $(minikube docker-env) && docker build -t account-service:latest -f account-service/Dockerfile .

build_mini_transaction:
	@eval $(minikube docker-env) && docker build -t transaction-service:latest -f transaction-service/Dockerfile .

build_mini_mailer:
	@eval $(minikube docker-env) && docker build -t mailer-service:latest -f mailer-service/Dockerfile .

deploy_gateway:
	@kubectl apply -f ./deployments/k8s/base/gateway/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/gateway/service.yaml

deploy_auth:
	@kubectl apply -f ./deployments/k8s/base/auth-service/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/auth-service/service.yaml

deploy_account:
	@kubectl apply -f ./deployments/k8s/base/account-service/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/account-service/service.yaml

deploy_transaction:
	@kubectl apply -f ./deployments/k8s/base/transaction-service/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/transaction-service/service.yaml

deploy_mailer:
	@kubectl apply -f ./deployments/k8s/base/mailer-service/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/mailer-service/service.yaml

deploy_redis:
	@kubectl apply -f ./deployments/k8s/base/mailer-service/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/mailer-service/service.yaml

deploy_mailhog:
	@kubectl apply -f ./deployments/k8s/base/mailhog/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/mailhog/service.yaml

deploy_rabbitmq:
	@kubectl apply -f ./deployments/k8s/base/rabbitmq/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/rabbitmq/service.yaml

deploy_postgres:
	@kubectl apply -f ./deployments/k8s/base/postgres/pvc.yaml
	@kubectl apply -f ./deployments/k8s/base/postgres/deployment.yaml
	@kubectl apply -f ./deployments/k8s/base/postgres/service.yaml


build_mini_all: build_mini_gateway build_mini_auth build_mini_account build_mini_transaction build_mini_mailer

deploy_all: deploy_gateway deploy_auth deploy_account deploy_transaction deploy_mailer deploy_postgres deploy_redis deploy_mailhog deploy_rabbitmq 


