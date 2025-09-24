# Ndugu Backend

Backend pour l'application Ndugu utilisant les services Ory pour l'authentification et l'autorisation.

## Services Ory intégrés

Ce projet utilise trois services Ory pour la gestion complète de l'identité et des autorisations :

- **Ory Kratos** : Gestion d'identité et authentification
- **Ory Hydra** : OAuth2 et OpenID Connect
- **Ory Keto** : Contrôle d'accès basé sur les attributs

## Démarrage rapide

1. **Cloner le projet**
   ```bash
   git clone <repository-url>
   cd ndugu-backend-repo
   ```

2. **Démarrer les services**
   ```bash
   docker-compose up -d
   ```

3. **Installer les dépendances Go**
   ```bash
   go mod tidy
   ```

4. **Démarrer l'API**
   ```bash
   go run services/coreapi/main.go
   ```

## Architecture

```
├── api/                    # Définitions gRPC/Protobuf
├── apisix/                 # Configuration API Gateway
├── internal/
│   └── auth/              # Module d'authentification Ory
├── ory/                   # Configuration des services Ory
│   ├── kratos/           # Configuration Kratos
│   ├── hydra/            # Configuration Hydra
│   └── keto/             # Configuration Keto
├── services/
│   └── coreapi/          # Service API principal
└── docker-compose.yml    # Orchestration des services
```

## Services disponibles

- **Backend API** : http://localhost:50051 (gRPC)
- **API Gateway (APISIX)** : http://localhost:9080
- **Kratos Public** : http://localhost:4433
- **Kratos Admin** : http://localhost:4434
- **Hydra Public** : http://localhost:4444
- **Hydra Admin** : http://localhost:4445
- **Keto Read** : http://localhost:4466
- **Keto Write** : http://localhost:4467
- **PostgreSQL** : localhost:5432

## Utilisation des services Ory

Voir le fichier `ory/README.md` pour la documentation détaillée des services Ory.

### Exemple d'utilisation

```go
// Créer un client Ory
oryClient := auth.NewOryClient()

// Créer un utilisateur
user, err := oryClient.CreateUser(ctx, "user@example.com", "John", "Doe")

// Créer un client OAuth2
client, err := oryClient.CreateOAuth2Client(ctx, "client-id", "Client Name", "http://localhost:3000/callback")

// Créer une permission
err := oryClient.CreatePermission(ctx, "files", "document1", "read", "user:123")

// Vérifier une permission
hasPermission, err := oryClient.CheckPermission(ctx, "files", "document1", "read", "user:123")
```

## Développement

### Prérequis
- Docker et Docker Compose
- Go 1.21+
- Git

### Tests
```bash
# Tester les services Ory
curl http://localhost:4433/health/ready  # Kratos
curl http://localhost:4444/health/ready  # Hydra
curl http://localhost:4466/health/ready  # Keto
```

### Configuration

Les fichiers de configuration Ory se trouvent dans le dossier `ory/`. Pour la production, assurez-vous de changer tous les secrets par défaut.

## Sécurité

⚠️ **Important** : Les secrets par défaut sont uniquement pour le développement. Changez tous les secrets pour la production.

```bash
# Images which grpc dependencies
docker run -it -v $(pwd):/src dscale.azurecr.io/library/dscale-go-grpctools-go1.23:v1.0.0 bash
```