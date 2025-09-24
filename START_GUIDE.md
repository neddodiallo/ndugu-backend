# 🚀 Guide de Démarrage - Ndugu Backend

## 📋 Prérequis

- Docker et Docker Compose installés
- Go 1.23+ (pour les tests locaux)
- grpcurl (pour les tests gRPC)

## 🔧 Installation de grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## 🚀 Démarrage du Projet

### Option 1: Script Automatique
```bash
chmod +x run.sh
./run.sh
```

### Option 2: Makefile
```bash
make start
```

### Option 3: Docker Compose Manuel
```bash
# Nettoyer les conteneurs existants
docker-compose down --remove-orphans

# Construire et démarrer
docker-compose up -d --build

# Attendre que les services soient prêts
sleep 15

# Configurer les routes APISIX
chmod +x scripts/setup-apisix-routes-final.sh
./scripts/setup-apisix-routes-final.sh
```

## 🔍 Vérification des Services

### État des Services
```bash
docker-compose ps
```

### Logs des Services
```bash
# Backend
docker-compose logs -f backend

# APISIX
docker-compose logs -f apisix

# Tous les services
docker-compose logs -f
```

## 🧪 Tests des Endpoints gRPC

### Test de Connectivité
```bash
# Vérifier que APISIX est accessible
curl http://localhost:9080

# Lister les services gRPC disponibles
grpcurl -plaintext localhost:9080 list
```

### Tests des Services
```bash
# Service Auth
grpcurl -plaintext localhost:9080 list ndugu.v1.AuthService
grpcurl -plaintext localhost:9080 describe ndugu.v1.AuthService

# Service Customer
grpcurl -plaintext localhost:9080 list ndugu.v1.CustomerService
grpcurl -plaintext localhost:9080 describe ndugu.v1.CustomerService
```

### Tests des Méthodes
```bash
# Créer un utilisateur
grpcurl -plaintext localhost:9080 ndugu.v1.AuthService/CreateUser

# Créer un client
grpcurl -plaintext localhost:9080 ndugu.v1.CustomerService/CreateCustomer
```

## 📊 Services Disponibles

### Backend gRPC
- **Port**: 50051
- **Services**: AuthService, CustomerService
- **Protocole**: gRPC uniquement

### APISIX Gateway
- **Port**: 9080 (Gateway)
- **Port**: 9180 (Admin API)
- **Fonction**: Routage gRPC

### ETCD
- **Port**: 2379
- **Fonction**: Configuration APISIX

### PostgreSQL
- **Port**: 5432
- **Base**: ndugu
- **Utilisateur**: user
- **Mot de passe**: password

### Services Ory
- **Kratos**: 4433 (public), 4434 (admin)
- **Hydra**: 4444 (public), 4445 (admin)
- **Keto**: 4466 (read), 4467 (write)

## 🛠️ Dépannage

### Problèmes Courants

#### 1. Services ne démarrent pas
```bash
# Vérifier les logs
docker-compose logs

# Redémarrer les services
docker-compose restart
```

#### 2. APISIX ne route pas
```bash
# Vérifier la configuration ETCD
curl http://localhost:2379/v2/keys/apisix/routes

# Reconfigurer les routes
./scripts/setup-apisix-routes-final.sh
```

#### 3. Backend ne compile pas
```bash
# Vérifier les dépendances
go mod tidy

# Compiler localement
go build ./services/coreapi/
```

### Commandes Utiles

```bash
# Nettoyer tout
docker-compose down --volumes --rmi all

# Reconstruire
docker-compose up -d --build --force-recreate

# Vérifier l'état
docker-compose ps
docker-compose logs --tail=50
```

## 🎯 Tests de Validation

### Script de Test Automatique
```bash
chmod +x test-grpc.sh
./test-grpc.sh
```

### Tests Manuels
```bash
# Test de santé
curl http://localhost:9080/health

# Test gRPC
grpcurl -plaintext localhost:9080 list
```

## 🛑 Arrêt du Projet

```bash
# Arrêter les services
docker-compose down

# Arrêter et supprimer les volumes
docker-compose down --volumes
```

## 📝 Notes

- Le projet utilise uniquement gRPC (pas de HTTP REST)
- APISIX fait le routage gRPC vers le backend
- Les services Ory sont configurés mais pas encore intégrés
- Les routes sont configurées dans ETCD via APISIX

## 🔗 Liens Utiles

- [Documentation APISIX](https://apisix.apache.org/docs/)
- [Documentation gRPC](https://grpc.io/docs/)
- [Documentation Ory](https://www.ory.sh/docs/)

