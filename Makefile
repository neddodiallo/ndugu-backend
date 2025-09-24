# Makefile pour Ndugu Backend

.PHONY: help build test test-verbose test-coverage clean run docker-build docker-up docker-down lint fmt start test-grpc test-project

# Variables
BINARY_NAME=ndugu-backend
DOCKER_IMAGE=ndugu-backend
DOCKER_TAG=latest

# Couleurs pour les messages
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Affiche l'aide
	@echo "$(GREEN)Ndugu Backend - Commandes disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

proto3:
	protoc -Iproto api/*.proto --go_out=./internal/grpc --go_opt=paths=source_relative --go-grpc_out=./internal/grpc --go-grpc_opt=paths=source_relative -I . \
		-I /usr/local/lib \
		-I /opt/include/ \
		-I /opt/include/google \
		-I /opt/include/protoc-gen-openapiv2/options/

build: ## Compile l'application
	@echo "$(GREEN)Compilation de l'application...$(NC)"
	go build -o $(BINARY_NAME) ./services/coreapi/
	@echo "$(GREEN)Compilation terminée: $(BINARY_NAME)$(NC)"

test: ## Exécute les tests unitaires
	@echo "$(GREEN)Exécution des tests unitaires...$(NC)"
	go test -v ./...

test-verbose: ## Exécute les tests avec plus de détails
	@echo "$(GREEN)Exécution des tests détaillés...$(NC)"
	go test -v -race ./...

test-coverage: ## Exécute les tests avec couverture de code
	@echo "$(GREEN)Exécution des tests avec couverture...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Rapport de couverture généré: coverage.html$(NC)"

test-short: ## Exécute les tests courts
	@echo "$(GREEN)Exécution des tests courts...$(NC)"
	go test -short ./...


clean: ## Nettoie les fichiers générés
	@echo "$(GREEN)Nettoyage des fichiers...$(NC)"
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean

run: build ## Compile et exécute l'application
	@echo "$(GREEN)Démarrage de l'application...$(NC)"
	./$(BINARY_NAME)

run-dev: ## Exécute l'application en mode développement
	@echo "$(GREEN)Démarrage en mode développement...$(NC)"
	go run ./services/coreapi/

lint: ## Exécute le linter
	@echo "$(GREEN)Exécution du linter...$(NC)"
	golangci-lint run

fmt: ## Formate le code
	@echo "$(GREEN)Formatage du code...$(NC)"
	go fmt ./...

mod-tidy: ## Nettoie les dépendances
	@echo "$(GREEN)Nettoyage des dépendances...$(NC)"
	go mod tidy

mod-update: ## Met à jour les dépendances
	@echo "$(GREEN)Mise à jour des dépendances...$(NC)"
	go get -u ./...

docker-build: ## Construit l'image Docker
	@echo "$(GREEN)Construction de l'image Docker...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-build-dev: ## Construit l'image Docker de développement
	@echo "$(GREEN)Construction de l'image Docker de développement...$(NC)"
	docker build -f Dockerfile.dev -t ndugu-backend-dev .

docker-up: ## Démarre les services Docker
	@echo "$(GREEN)Démarrage des services Docker...$(NC)"
	docker-compose up -d

docker-dev: ## Démarre les services Docker en mode développement
	@echo "$(GREEN)Démarrage des services Docker en mode développement...$(NC)"
	docker-compose -f docker-compose.dev.yml up -d

docker-down: ## Arrête les services Docker
	@echo "$(GREEN)Arrêt des services Docker...$(NC)"
	docker-compose down

docker-down-dev: ## Arrête les services Docker de développement
	@echo "$(GREEN)Arrêt des services Docker de développement...$(NC)"
	docker-compose -f docker-compose.dev.yml down

docker-logs: ## Affiche les logs Docker
	@echo "$(GREEN)Logs des services Docker...$(NC)"
	docker-compose logs -f

docker-logs-dev: ## Affiche les logs Docker de développement
	@echo "$(GREEN)Logs des services Docker de développement...$(NC)"
	docker-compose -f docker-compose.dev.yml logs -f

docker-restart: ## Redémarre les services Docker
	@echo "$(GREEN)Redémarrage des services Docker...$(NC)"
	docker-compose restart

docker-clean: ## Nettoie les images et conteneurs Docker
	@echo "$(GREEN)Nettoyage des images et conteneurs Docker...$(NC)"
	docker-compose down --volumes --rmi all
	docker system prune -f

test-integration: ## Exécute les tests d'intégration
	@echo "$(GREEN)Exécution des tests d'intégration...$(NC)"
	@echo "$(YELLOW)Assurez-vous que les services Docker sont démarrés$(NC)"
	./test_endpoints.sh

test-grpc: ## Teste la connexion gRPC
	@echo "$(GREEN)Test de la connexion gRPC...$(NC)"
	go run test_grpc.go

install-tools: ## Installe les outils de développement
	@echo "$(GREEN)Installation des outils de développement...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

check-deps: ## Vérifie les dépendances
	@echo "$(GREEN)Vérification des dépendances...$(NC)"
	go mod verify

security-scan: ## Exécute un scan de sécurité
	@echo "$(GREEN)Scan de sécurité...$(NC)"
	gosec ./...

benchmark: ## Exécute les benchmarks
	@echo "$(GREEN)Exécution des benchmarks...$(NC)"
	go test -bench=. ./...

# Migrations Ory
migrate-kratos: ## Applique les migrations Kratos
	@echo "$(GREEN)Application des migrations Kratos...$(NC)"
	docker-compose run --rm -e DSN="postgres://user:password@db:5432/ndugu?sslmode=disable" kratos migrate sql -e --yes
	@echo "$(GREEN)Migrations Kratos appliquées$(NC)"

migrate-hydra: ## Applique les migrations Hydra
	@echo "$(GREEN)Application des migrations Hydra...$(NC)"
	docker-compose run --rm -e DSN="postgres://user:password@db:5432/ndugu?sslmode=disable" hydra migrate sql -e --yes
	@echo "$(GREEN)Migrations Hydra appliquées$(NC)"

migrate-keto: ## Applique les migrations Keto
	@echo "$(GREEN)Application des migrations Keto...$(NC)"
	docker-compose run --rm keto migrate up --yes -c /etc/config/keto/keto.yml
	@echo "$(GREEN)Migrations Keto appliquées$(NC)"

migrate-all: migrate-kratos migrate-hydra migrate-keto ## Applique toutes les migrations Ory
	@echo "$(GREEN)Toutes les migrations Ory ont été appliquées$(NC)"

# Configuration APISIX
setup-apisix: ## Configure les routes APISIX via l'Admin API
	@echo "$(GREEN)Configuration des routes APISIX...$(NC)"
	@./scripts/setup-apisix-routes.sh
	@echo "$(GREEN)Routes APISIX configurées$(NC)"

# Démarrage et tests
start: ## Démarre le projet complet
	@echo "$(GREEN)Démarrage du projet Ndugu Backend...$(NC)"
	@chmod +x start-project.sh
	@./start-project.sh

test-project: ## Teste la configuration du projet
	@echo "$(GREEN)Test de la configuration du projet...$(NC)"
	@chmod +x test-project.sh
	@./test-project.sh

test-grpc: ## Teste les endpoints gRPC
	@echo "$(GREEN)Test des endpoints gRPC...$(NC)"
	@chmod +x test-grpc.sh
	@./test-grpc.sh

# Règles par défaut
.DEFAULT_GOAL := help
