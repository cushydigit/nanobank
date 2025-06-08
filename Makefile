.PHONY: tidy up down help reset start stop status build deploy delete restart full-reset logs logs-auth logs-account open-gateway open-rabbitmq open-mailhog delete-pvcs

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

## DEVELOPMENT

COMPOSE_PROJECT_NAME=nanobank
COMPOSE_FILE=./deployments/docker-compose.yml
SERVICE_NAME_1=gateway
SERVICE_NAME_2=auth-service
SERVICE_NAME_3=account-service
SERVICE_NAME_4=transaction-service
SERVICE_NAME_5=mailer-service
SERVICE_SHARED=shared
SERVICE_NAME_6=postgres
SERVICE_NAME_7=redis
SERVICE_NAME_8=mailhog
SERVICE_NAME_9=rabbitmq

tidy: ## Run go mod tidy in all services
	@cd ./$(SERVICE_NAME_1)/ && go mod tidy
	@cd ./$(SERVICE_NAME_2)/ && go mod tidy
	@cd ./$(SERVICE_NAME_3)/ && go mod tidy
	@cd ./$(SERVICE_NAME_4)/ && go mod tidy
	@cd ./$(SERVICE_NAME_5)/ && go mod tidy
	@cd ./$(SERVICE_SHARED)/ && go mod tidy

up: ## Start Docker services with build
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) up -d --build

down: ## Stop and remove Docker containers
	docker-compose -p $(COMPOSE_PROJECT_NAME) -f $(COMPOSE_FILE) down

reset: down tidy up ## Reset Docker stack (down + tidy + up)

## DEPLOYMENT

KUBECTL=kubectl
MINIKUBE=minikube
K8S_DIR=./deployments/k8s/base

start: ## Start Minikube
	@$(MINIKUBE) start

stop: ## Stop Minikube
	@$(MINIKUBE) stop

status: ## Show status of Minikube resources
	@$(MINIKUBE) get all

build: ## Build Docker images for services inside Minikube
	@eval $(minikube docker-env) && docker build -t $(SERVICE_NAME_1):latest -f $(SERVICE_NAME_1)/Dockerfile .
	@eval $(minikube docker-env) && docker build -t $(SERVICE_NAME_2):latest -f $(SERVICE_NAME_2)/Dockerfile .
	@eval $(minikube docker-env) && docker build -t $(SERVICE_NAME_3):latest -f $(SERVICE_NAME_3)/Dockerfile .
	@eval $(minikube docker-env) && docker build -t $(SERVICE_NAME_4):latest -f $(SERVICE_NAME_4)/Dockerfile .
	@eval $(minikube docker-env) && docker build -t $(SERVICE_NAME_5):latest -f $(SERVICE_NAME_5)/Dockerfile .

deploy: ## Apply Kubernetes manifests for all services
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_1)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_2)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_3)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_4)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_5)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_6)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_7)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_8)
	@$(KUBECTL) apply -f $(K8S_DIR)/$(SERVICE_NAME_9)

delete: ## Delete Kubernetes resources
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_1) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_2) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_3) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_4) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_5) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_6) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_7) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_8) || true
	@$(KUBECTL) delete -f $(K8S_DIR)/$(SERVICE_NAME_9) || true

delete-pvcs: ## Delete all PVCs in current namespace
	@$(KUBECTL) delete pvc --all

restart: delete build deploy ## Restart Kubernetes services

full-reset: delete delete-pvcs stop start build deploy ## Full Minikube + Kubernetes reset

logs: ## Tail logs of gateway service
	@$(KUBECTL) logs -l app=$(SERVICE_NAME_1) --tail=100 -f

logs-auth: ## Tail logs of auth service
	@$(KUBECTL) logs -l app=$(SERVICE_NAME_2) --tail=100 -f

logs-account: ## Tail logs of account service
	@$(KUBECTL) logs -l app=$(SERVICE_NAME_3) --tail=100 -f

open-gateway: ## Open gateway in Minikube
	@$(MINIKUBE) service $(SERVICE_NAME_1)

open-rabbitmq: ## Open RabbitMQ web interface in Minikube
	@$(MINIKUBE) service $(SERVICE_NAME_9)

open-mailhog: ## Open MailHog web interface in Minikube
	@$(MINIKUBE) service $(SERVICE_NAME_8)

