#!/bin/bash

echo "🚀 Démarrage du projet Ndugu Backend"
echo "===================================="

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Fonction pour afficher les résultats
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Étape 1: Nettoyer les conteneurs existants
echo "🧹 Nettoyage des conteneurs existants..."
docker-compose down --remove-orphans 2>/dev/null || true

# Étape 2: Construire et démarrer les services
echo "🔨 Construction et démarrage des services..."
if docker-compose up -d --build; then
    print_status 0 "Services démarrés avec succès"
else
    print_status 1 "Erreur lors du démarrage des services"
    exit 1
fi

# Étape 3: Attendre que les services soient prêts
echo "⏳ Attente que les services soient prêts..."
sleep 10

# Étape 4: Vérifier l'état des services
echo "🔍 Vérification de l'état des services..."
docker-compose ps

# Étape 5: Configurer les routes APISIX
echo "🛣️  Configuration des routes APISIX..."
if [ -f "scripts/setup-apisix-routes-final.sh" ]; then
    chmod +x scripts/setup-apisix-routes-final.sh
    ./scripts/setup-apisix-routes-final.sh
else
    echo -e "${YELLOW}⚠️  Script de configuration APISIX non trouvé${NC}"
fi

# Étape 6: Afficher les informations de test
echo ""
echo "🎯 Projet démarré avec succès!"
echo ""
echo "📋 Services disponibles:"
echo "   - Backend gRPC: localhost:50051"
echo "   - APISIX Gateway: localhost:9080"
echo "   - ETCD: localhost:2379"
echo "   - PostgreSQL: localhost:5432"
echo ""
echo "🔧 Services Ory:"
echo "   - Kratos: localhost:4433 (public), localhost:4434 (admin)"
echo "   - Hydra: localhost:4444 (public), localhost:4445 (admin)"
echo "   - Keto: localhost:4466 (read), localhost:4467 (write)"
echo ""
echo "🧪 Tests gRPC:"
echo "   grpcurl -plaintext localhost:9080 ndugu.v1.AuthService/CreateUser"
echo "   grpcurl -plaintext localhost:9080 ndugu.v1.CustomerService/CreateCustomer"
echo ""
echo "📊 Logs des services:"
echo "   docker-compose logs -f backend"
echo "   docker-compose logs -f apisix"
echo ""
echo "🛑 Pour arrêter:"
echo "   docker-compose down"
