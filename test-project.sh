#!/bin/bash

echo "🚀 Test du projet Ndugu Backend"
echo "================================"

# Couleurs pour la sortie
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Test 1: Vérifier la compilation Go
echo "📦 Test de compilation Go..."
if go build -o backend ./services/coreapi/; then
    print_result 0 "Compilation Go réussie"
    rm -f backend
else
    print_result 1 "Erreur de compilation Go"
    exit 1
fi

# Test 2: Vérifier les dépendances
echo "📋 Test des dépendances..."
if go mod tidy; then
    print_result 0 "Dépendances Go nettoyées"
else
    print_result 1 "Erreur lors du nettoyage des dépendances"
fi

# Test 3: Vérifier Docker
echo "🐳 Test de Docker..."
if docker --version > /dev/null 2>&1; then
    print_result 0 "Docker est installé"
else
    print_result 1 "Docker n'est pas installé ou accessible"
fi

# Test 4: Vérifier Docker Compose
echo "🐳 Test de Docker Compose..."
if docker-compose --version > /dev/null 2>&1; then
    print_result 0 "Docker Compose est installé"
else
    print_result 1 "Docker Compose n'est pas installé ou accessible"
fi

# Test 5: Vérifier la configuration
echo "⚙️  Test de la configuration..."
if [ -f "docker-compose.yml" ]; then
    print_result 0 "docker-compose.yml trouvé"
else
    print_result 1 "docker-compose.yml manquant"
fi

if [ -f "Dockerfile" ]; then
    print_result 0 "Dockerfile trouvé"
else
    print_result 1 "Dockerfile manquant"
fi

if [ -f "apisix/config.yaml" ]; then
    print_result 0 "Configuration APISIX trouvée"
else
    print_result 1 "Configuration APISIX manquante"
fi

# Test 6: Vérifier les routes APISIX
echo "🛣️  Test des routes APISIX..."
if [ -f "apisix/routes.yaml" ]; then
    print_result 0 "Routes APISIX trouvées"
    echo "   Routes configurées:"
    grep -E "uri:|scheme:" apisix/routes.yaml | sed 's/^/   /'
else
    print_result 1 "Routes APISIX manquantes"
fi

echo ""
echo "🎯 Résumé des tests:"
echo "   - Backend: gRPC uniquement (port 50051)"
echo "   - APISIX: Gateway gRPC (port 9080)"
echo "   - ETCD: Configuration (port 2379)"
echo "   - Services Ory: Kratos, Hydra, Keto"
echo ""
echo "🧪 Pour tester le projet:"
echo "   1. docker-compose up -d"
echo "   2. grpcurl -plaintext localhost:9080 ndugu.v1.AuthService/CreateUser"
echo "   3. grpcurl -plaintext localhost:9080 ndugu.v1.CustomerService/CreateCustomer"
echo ""
echo "✅ Tests terminés!"
