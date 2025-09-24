#!/bin/bash

echo "🧪 Test des endpoints Ndugu Backend avec Ory"
echo "=============================================="

# Couleurs pour les résultats
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour tester un endpoint
test_endpoint() {
    local name="$1"
    local url="$2"
    local method="${3:-GET}"
    local data="$4"
    
    echo -n "Testing $name... "
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$url")
    else
        response=$(curl -s -w "%{http_code}" "$url")
    fi
    
    http_code="${response: -3}"
    body="${response%???}"
    
    if [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
        echo -e "${GREEN}✅ OK (HTTP $http_code)${NC}"
        if [ -n "$body" ] && [ "$body" != "$http_code" ]; then
            echo "   Response: $body"
        fi
    else
        echo -e "${RED}❌ FAILED (HTTP $http_code)${NC}"
        if [ -n "$body" ] && [ "$body" != "$http_code" ]; then
            echo "   Error: $body"
        fi
    fi
    echo
}

echo "1. 🔍 Vérification des services Docker"
echo "--------------------------------------"
docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
echo

echo "2. 🌐 Test des services Ory"
echo "---------------------------"

# Test Kratos
test_endpoint "Kratos Health" "http://localhost:4433/health/alive"
test_endpoint "Kratos Ready" "http://localhost:4433/health/ready"

# Test Hydra
test_endpoint "Hydra Health" "http://localhost:4444/health/ready"

# Test Keto
test_endpoint "Keto Health" "http://localhost:4466/health/ready"

echo "3. 🚀 Test du Backend (via conteneur)"
echo "------------------------------------"

# Test du backend via le conteneur
echo -n "Testing Backend Health (via container)... "
backend_health=$(docker exec ndugu-backend-repo-backend-1 wget -qO- http://localhost:8080/health 2>/dev/null || echo "FAILED")
if [ "$backend_health" = "Services Ory opérationnels" ]; then
    echo -e "${GREEN}✅ OK${NC}"
    echo "   Response: $backend_health"
else
    echo -e "${RED}❌ FAILED${NC}"
    echo "   Response: $backend_health"
fi
echo

echo "4. 📡 Test gRPC"
echo "---------------"

# Test gRPC avec grpcurl
echo -n "Testing gRPC Server... "
grpc_test=$(~/go/bin/grpcurl -plaintext localhost:50051 list 2>/dev/null | grep -c "grpc.reflection" || echo "0")
if [ "$grpc_test" -gt 0 ]; then
    echo -e "${GREEN}✅ OK${NC}"
    echo "   Services disponibles:"
    ~/go/bin/grpcurl -plaintext localhost:50051 list | sed 's/^/     /'
else
    echo -e "${RED}❌ FAILED${NC}"
fi
echo

echo "5. 🌍 Test API Gateway (APISIX)"
echo "-------------------------------"
test_endpoint "APISIX Health" "http://localhost:9080/health"

echo "6. 📊 Résumé des tests"
echo "====================="

# Compter les services qui fonctionnent
working_services=0
total_services=6

# Vérifier chaque service
if curl -s http://localhost:4433/health/alive >/dev/null 2>&1; then ((working_services++)); fi
if curl -s http://localhost:4444/health/ready >/dev/null 2>&1; then ((working_services++)); fi
if curl -s http://localhost:4466/health/ready >/dev/null 2>&1; then ((working_services++)); fi
if docker exec ndugu-backend-repo-backend-1 wget -qO- http://localhost:8080/health >/dev/null 2>&1; then ((working_services++)); fi
if ~/go/bin/grpcurl -plaintext localhost:50051 list >/dev/null 2>&1; then ((working_services++)); fi
if curl -s http://localhost:9080/health >/dev/null 2>&1; then ((working_services++)); fi

echo "Services fonctionnels: $working_services/$total_services"

if [ "$working_services" -eq "$total_services" ]; then
    echo -e "${GREEN}🎉 Tous les services fonctionnent parfaitement !${NC}"
elif [ "$working_services" -gt 3 ]; then
    echo -e "${YELLOW}⚠️  La plupart des services fonctionnent, quelques ajustements nécessaires${NC}"
else
    echo -e "${RED}❌ Plusieurs services ont des problèmes${NC}"
fi

echo
echo "📝 Notes:"
echo "- Les services Ory peuvent prendre quelques minutes pour être complètement prêts"
echo "- La base de données doit être initialisée pour Kratos"
echo "- APISIX nécessite ETCD pour fonctionner"
echo "- Les endpoints gRPC nécessitent la génération du code protobuf"
