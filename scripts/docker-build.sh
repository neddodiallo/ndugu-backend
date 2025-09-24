#!/bin/bash

# Script de construction Docker pour Ndugu Backend
# Usage: ./scripts/docker-build.sh [dev|prod]

set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifier les arguments
BUILD_TYPE=${1:-prod}

if [[ "$BUILD_TYPE" != "dev" && "$BUILD_TYPE" != "prod" ]]; then
    log_error "Usage: $0 [dev|prod]"
    exit 1
fi

log_info "Construction Docker pour Ndugu Backend (mode: $BUILD_TYPE)"

# Aller dans le répertoire du projet
cd "$(dirname "$0")/.."

# Vérifier que Docker est installé
if ! command -v docker &> /dev/null; then
    log_error "Docker n'est pas installé"
    exit 1
fi

# Vérifier que Docker Compose est installé
if ! command -v docker-compose &> /dev/null; then
    log_error "Docker Compose n'est pas installé"
    exit 1
fi

# Nettoyer les images existantes
log_info "Nettoyage des images existantes..."
docker-compose down --volumes --rmi all 2>/dev/null || true

# Construire l'image selon le type
if [[ "$BUILD_TYPE" == "dev" ]]; then
    log_info "Construction de l'image de développement..."
    docker build -f Dockerfile.dev -t ndugu-backend-dev .
    log_success "Image de développement construite: ndugu-backend-dev"
else
    log_info "Construction de l'image de production..."
    docker build -t ndugu-backend .
    log_success "Image de production construite: ndugu-backend"
fi

# Vérifier que l'image a été construite
if [[ "$BUILD_TYPE" == "dev" ]]; then
    if docker images | grep -q "ndugu-backend-dev"; then
        log_success "Image de développement vérifiée"
    else
        log_error "Échec de la construction de l'image de développement"
        exit 1
    fi
else
    if docker images | grep -q "ndugu-backend"; then
        log_success "Image de production vérifiée"
    else
        log_error "Échec de la construction de l'image de production"
        exit 1
    fi
fi

# Afficher les informations sur l'image
log_info "Informations sur l'image construite:"
if [[ "$BUILD_TYPE" == "dev" ]]; then
    docker images | grep "ndugu-backend-dev"
else
    docker images | grep "ndugu-backend"
fi

# Proposer de démarrer les services
echo ""
log_info "Construction terminée avec succès!"
echo ""
log_info "Commandes disponibles:"
if [[ "$BUILD_TYPE" == "dev" ]]; then
    echo "  - Démarrer en mode développement: make docker-dev"
    echo "  - Voir les logs: make docker-logs-dev"
    echo "  - Arrêter: make docker-down-dev"
else
    echo "  - Démarrer en mode production: make docker-up"
    echo "  - Voir les logs: make docker-logs"
    echo "  - Arrêter: make docker-down"
fi
echo ""
log_info "URLs d'accès:"
echo "  - API: http://localhost:8080"
echo "  - API REST: http://localhost:8080/api/v1/"
echo "  - Health Check: http://localhost:8080/health"
echo ""

# Demander si l'utilisateur veut démarrer les services
read -p "Voulez-vous démarrer les services maintenant? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [[ "$BUILD_TYPE" == "dev" ]]; then
        log_info "Démarrage des services de développement..."
        make docker-dev
    else
        log_info "Démarrage des services de production..."
        make docker-up
    fi
    log_success "Services démarrés!"
    echo ""
    log_info "Vous pouvez maintenant accéder à:"
    echo "  - API REST: http://localhost:8080/api/v1/"
fi

log_success "Script terminé avec succès!"
