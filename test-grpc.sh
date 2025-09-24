#!/bin/bash

echo "🧪 Test des endpoints gRPC"
echo "=========================="

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Fonction pour tester un endpoint gRPC
test_grpc_endpoint() {
    local service=$1
    local method=$2
    local description=$3
    
    echo "🔍 Test: $description"
    echo "   Service: $service/$method"
    
    if grpcurl -plaintext localhost:9080 list | grep -q "$service"; then
        echo -e "   ${GREEN}✅ Service $service disponible${NC}"
        
        # Test de la méthode
        if grpcurl -plaintext localhost:9080 list "$service" | grep -q "$method"; then
            echo -e "   ${GREEN}✅ Méthode $method disponible${NC}"
        else
            echo -e "   ${YELLOW}⚠️  Méthode $method non trouvée${NC}"
        fi
    else
        echo -e "   ${RED}❌ Service $service non disponible${NC}"
    fi
    echo ""
}

# Vérifier si grpcurl est installé
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}❌ grpcurl n'est pas installé${NC}"
    echo "   Installation: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Attendre que les services soient prêts
echo "⏳ Attente que les services soient prêts..."
sleep 5

# Test 1: Vérifier la connectivité APISIX
echo "🌐 Test de connectivité APISIX..."
if curl -s http://localhost:9080 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ APISIX accessible${NC}"
else
    echo -e "${RED}❌ APISIX non accessible${NC}"
    echo "   Vérifiez que les services sont démarrés: docker-compose ps"
    exit 1
fi

# Test 2: Lister les services disponibles
echo "📋 Services gRPC disponibles:"
grpcurl -plaintext localhost:9080 list
echo ""

# Test 3: Tester les services spécifiques
test_grpc_endpoint "ndugu.v1.AuthService" "CreateUser" "Création d'utilisateur"
test_grpc_endpoint "ndugu.v1.AuthService" "GetUser" "Récupération d'utilisateur"
test_grpc_endpoint "ndugu.v1.AuthService" "ValidateSession" "Validation de session"
test_grpc_endpoint "ndugu.v1.CustomerService" "CreateCustomer" "Création de client"

# Test 4: Test de réflexion gRPC
echo "🔍 Test de réflexion gRPC..."
if grpcurl -plaintext localhost:9080 list | grep -q "grpc.reflection"; then
    echo -e "${GREEN}✅ Réflexion gRPC activée${NC}"
else
    echo -e "${YELLOW}⚠️  Réflexion gRPC non activée${NC}"
fi

echo ""
echo "🎯 Tests terminés!"
echo ""
echo "💡 Pour des tests plus détaillés:"
echo "   grpcurl -plaintext localhost:9080 describe ndugu.v1.AuthService"
echo "   grpcurl -plaintext localhost:9080 describe ndugu.v1.CustomerService"
