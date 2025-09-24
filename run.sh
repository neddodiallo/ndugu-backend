#!/bin/bash

echo "🚀 Lancement du projet Ndugu Backend"
echo "===================================="

# Vérifier si Docker est en cours d'exécution
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker n'est pas en cours d'exécution"
    echo "   Veuillez démarrer Docker Desktop"
    exit 1
fi

# Nettoyer les conteneurs existants
echo "🧹 Nettoyage des conteneurs existants..."
docker-compose down --remove-orphans 2>/dev/null || true

# Construire et démarrer les services
echo "🔨 Construction et démarrage des services..."
docker-compose up -d --build

# Attendre que les services soient prêts
echo "⏳ Attente que les services soient prêts..."
sleep 15

# Vérifier l'état des services
echo "🔍 État des services:"
docker-compose ps

# Configurer les routes APISIX
echo "🛣️  Configuration des routes APISIX..."
if [ -f "scripts/setup-apisix-routes-final.sh" ]; then
    chmod +x scripts/setup-apisix-routes-final.sh
    ./scripts/setup-apisix-routes-final.sh
fi

echo ""
echo "✅ Projet lancé avec succès!"
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
echo "   grpcurl -plaintext localhost:9080 list"
echo "   grpcurl -plaintext localhost:9080 ndugu.v1.AuthService/CreateUser"
echo ""
echo "📊 Logs:"
echo "   docker-compose logs -f backend"
echo "   docker-compose logs -f apisix"
echo ""
echo "🛑 Arrêt:"
echo "   docker-compose down"

