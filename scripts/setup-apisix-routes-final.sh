#!/bin/bash

# Script final pour configurer les routes APISIX avec le format correct

set -e

ETCD_URL="http://localhost:2379"
APISIX_PREFIX="/apisix"

echo "🚀 Configuration finale des routes APISIX..."

# Fonction pour attendre qu'APISIX soit synchronisé avec ETCD
wait_for_apisix_sync() {
    echo "⏳ Attente que APISIX soit synchronisé avec ETCD..."
    for i in {1..30}; do
        # Vérifier si APISIX peut lire les routes depuis ETCD
        if curl -s -f "http://localhost:9080/apisix/admin/routes" -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" > /dev/null 2>&1; then
            echo "✅ APISIX est synchronisé avec ETCD!"
            return 0
        fi
        echo "Tentative $i/30..."
        sleep 3
    done
    echo "⚠️  APISIX n'est pas encore synchronisé, mais on continue..."
    return 1
}

# Fonction pour créer une route avec le format correct
create_route_final() {
    local route_id=$1
    local route_config=$2
    
    echo "📝 Création de la route $route_id avec format final..."
    
    # Créer la route avec le format correct (sans double échappement)
    response=$(curl -s -X PUT "$ETCD_URL/v2/keys$APISIX_PREFIX/routes/$route_id" \
        -d "value=$route_config")
    
    echo "✅ Route $route_id créée"
}

# Supprimer les routes de test existantes
echo "🧹 Nettoyage des routes de test..."
curl -s -X DELETE "$ETCD_URL/v2/keys/apisix/routes/test1" > /dev/null 2>&1 || true
curl -s -X DELETE "$ETCD_URL/v2/keys/apisix/routes/test2" > /dev/null 2>&1 || true
curl -s -X DELETE "$ETCD_URL/v2/keys/apisix/routes/test3" > /dev/null 2>&1 || true

# Créer les routes finales avec le format correct
echo "🔧 Création des routes finales..."

# Route 1: Service Customer gRPC
create_route_final "1" '{"uri":"/ndugu.v1.CustomerService/*","plugins":{"grpc-proxy":{"socket_timeout":60}},"upstream":{"nodes":{"backend:50051":1},"scheme":"grpc"}}'

# Route 2: Service Auth gRPC
create_route_final "2" '{"uri":"/ndugu.v1.AuthService/*","plugins":{"grpc-proxy":{"socket_timeout":60}},"upstream":{"nodes":{"backend:50051":1},"scheme":"grpc"}}'

# Attendre que APISIX soit synchronisé
wait_for_apisix_sync

echo "🎉 Configuration des routes APISIX terminée!"
echo ""
echo "📋 Routes configurées:"
echo "   - Route 1: gRPC CustomerService (/ndugu.v1.CustomerService/*)"
echo "   - Route 2: gRPC AuthService (/ndugu.v1.AuthService/*)"
echo ""
echo "🔗 Endpoints disponibles:"
echo "   - APISIX Gateway: http://localhost:9080"
echo "   - ETCD: http://localhost:2379"
echo ""
echo "🧪 Tests à effectuer:"
echo "   - grpcurl -plaintext localhost:9080 ndugu.v1.CustomerService/CreateCustomer"
echo "   - grpcurl -plaintext localhost:9080 ndugu.v1.AuthService/Login"
