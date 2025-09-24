# Architecture du Projet Ndugu Backend

## Vue d'ensemble

Le projet Ndugu Backend est organisé selon une architecture en couches (layered architecture) avec une séparation claire des responsabilités. Cette structure facilite la maintenance, les tests et l'évolution du code.

## Structure du Projet

```
ndugu-backend-repo/
├── api/                          # Définitions des API (protobuf, OpenAPI)
│   ├── ndugu.proto              # Définition gRPC
│   └── v1/                      # Versions des API
├── apisix/                       # Configuration API Gateway
│   ├── config.yaml
│   └── routes.yaml
├── internal/                     # Code interne de l'application
│   ├── auth/                    # Client Ory (legacy)
│   ├── common/                  # Utilitaires communs
│   │   ├── errors.go           # Gestion des erreurs
│   │   ├── logger.go           # Logging
│   │   ├── response.go         # Réponses HTTP standardisées
│   │   └── validation.go       # Validation des données
│   ├── config/                 # Configuration
│   │   └── config.go
│   ├── models/                 # Modèles de données
│   │   ├── user.go            # Modèle utilisateur
│   │   ├── auth.go            # Modèles d'authentification
│   │   └── customer.go        # Modèle client
│   ├── repository/            # Couche d'accès aux données
│   │   ├── interfaces.go      # Interfaces des repositories
│   │   ├── ory_client.go      # Client Ory
│   │   ├── kratos_client.go   # Client Kratos
│   │   ├── hydra_client.go    # Client Hydra
│   │   ├── keto_client.go     # Client Keto
│   │   └── mock_user_repository.go # Mock pour les tests
│   └── services/              # Couche métier
│       ├── auth_service.go    # Service d'authentification
│       └── auth_service_test.go # Tests du service
├── services/                   # Services de l'application
│   └── coreapi/               # Service principal
│       ├── main.go           # Point d'entrée
│       └── handler/          # Handlers HTTP
│           ├── auth_handler.go
│           └── http_server.go
├── ory/                       # Configuration Ory
│   ├── kratos/
│   ├── hydra/
│   └── keto/
├── migrations/                # Migrations de base de données
├── docker-compose.yml         # Orchestration Docker
├── Dockerfile                 # Image Docker
├── Makefile                   # Commandes de développement
└── README.md                  # Documentation principale
```

## Architecture en Couches

### 1. Couche de Présentation (Handlers)
- **Responsabilité** : Gestion des requêtes HTTP/gRPC
- **Localisation** : `services/coreapi/handler/`
- **Composants** :
  - `AuthHandler` : Gestion des endpoints d'authentification
  - `HTTPServer` : Serveur HTTP avec middleware

### 2. Couche Métier (Services)
- **Responsabilité** : Logique métier et orchestration
- **Localisation** : `internal/services/`
- **Composants** :
  - `AuthService` : Service d'authentification
  - Interfaces pour les services

### 3. Couche d'Accès aux Données (Repository)
- **Responsabilité** : Accès aux données et intégration externe
- **Localisation** : `internal/repository/`
- **Composants** :
  - `UserRepository` : Gestion des utilisateurs
  - `OryClient` : Intégration avec les services Ory
  - Implémentations mock pour les tests

### 4. Couche de Modèles
- **Responsabilité** : Définition des structures de données
- **Localisation** : `internal/models/`
- **Composants** :
  - `User` : Modèle utilisateur
  - `Auth` : Modèles d'authentification
  - `Customer` : Modèle client

### 5. Couche Commune
- **Responsabilité** : Utilitaires partagés
- **Localisation** : `internal/common/`
- **Composants** :
  - `errors.go` : Gestion centralisée des erreurs
  - `logger.go` : Système de logging
  - `response.go` : Réponses HTTP standardisées
  - `validation.go` : Validation des données

## Principes Architecturaux

### 1. Séparation des Responsabilités
- Chaque couche a une responsabilité claire et bien définie
- Les dépendances vont uniquement vers les couches inférieures
- Pas de dépendances circulaires

### 2. Inversion de Dépendance
- Les services dépendent d'interfaces, pas d'implémentations
- Facilite les tests et la maintenance
- Permet de changer d'implémentation sans impact

### 3. Testabilité
- Chaque couche peut être testée indépendamment
- Mocks disponibles pour les tests unitaires
- Tests d'intégration séparés

### 4. Configuration Externalisée
- Configuration via variables d'environnement
- Valeurs par défaut sensées
- Support de différents environnements

## Flux de Données

### Création d'Utilisateur
1. **Handler** reçoit la requête HTTP POST `/api/v1/users`
2. **Service** valide les données et orchestre la création
3. **Repository** sauvegarde en base de données locale
4. **OryClient** crée l'utilisateur dans Kratos
5. **Response** retourne l'utilisateur créé

### Validation de Session
1. **Handler** reçoit la requête POST `/api/v1/auth/validate`
2. **Service** délègue à OryClient
3. **KratosClient** valide la session avec Kratos
4. **Response** retourne le statut de validation

## Gestion des Erreurs

### Types d'Erreurs
- **Erreurs de validation** : Données d'entrée invalides
- **Erreurs métier** : Règles métier violées
- **Erreurs techniques** : Problèmes d'infrastructure
- **Erreurs Ory** : Problèmes avec les services Ory

### Standardisation
- Toutes les erreurs utilisent `AppError`
- Codes d'erreur standardisés
- Messages d'erreur cohérents
- Codes HTTP appropriés

## Tests

### Tests Unitaires
- **Modèles** : Validation des structures de données
- **Services** : Logique métier avec mocks
- **Handlers** : Gestion des requêtes HTTP
- **Utilitaires** : Fonctions communes

### Tests d'Intégration
- **Endpoints** : Tests complets des API
- **Services Ory** : Intégration avec Kratos/Hydra/Keto
- **Base de données** : Tests avec vraie DB

### Outils de Test
- **Makefile** : Commandes de test standardisées
- **Coverage** : Rapport de couverture de code
- **Mocks** : Implémentations mock pour les tests

## Déploiement

### Docker
- **Multi-stage build** : Optimisation de l'image
- **Docker Compose** : Orchestration des services
- **Variables d'environnement** : Configuration flexible

### API Gateway
- **APISIX** : Routage et gestion des API
- **Load balancing** : Distribution de charge
- **Rate limiting** : Protection contre les abus

## Évolutivité

### Ajout de Nouvelles Fonctionnalités
1. Créer le modèle dans `internal/models/`
2. Définir l'interface du repository
3. Implémenter le service métier
4. Créer le handler HTTP
5. Ajouter les tests

### Ajout de Nouveaux Services
1. Créer le package dans `internal/services/`
2. Définir l'interface du service
3. Implémenter la logique métier
4. Créer les tests unitaires
5. Intégrer dans le handler

## Bonnes Pratiques

### Code
- **Noms explicites** : Variables et fonctions claires
- **Documentation** : Commentaires pour la logique complexe
- **Formatage** : Code formaté avec `go fmt`
- **Linting** : Vérification avec `golangci-lint`

### Tests
- **Couverture** : Minimum 80% de couverture
- **Tests rapides** : Tests unitaires < 100ms
- **Tests isolés** : Pas de dépendances externes
- **Assertions claires** : Messages d'erreur explicites

### Git
- **Commits atomiques** : Un commit = une fonctionnalité
- **Messages clairs** : Description de ce qui a changé
- **Branches** : Feature branches pour les nouvelles fonctionnalités
- **Reviews** : Code review obligatoire

## Outils de Développement

### Makefile
- `make test` : Exécuter les tests
- `make build` : Compiler l'application
- `make run` : Démarrer l'application
- `make docker-up` : Démarrer les services Docker

### Dépendances
- **Go 1.23+** : Version minimale requise
- **Docker** : Pour l'orchestration des services
- **golangci-lint** : Linter Go
- **grpcurl** : Client gRPC pour les tests

Cette architecture garantit un code maintenable, testable et évolutif, tout en respectant les bonnes pratiques de développement Go.
