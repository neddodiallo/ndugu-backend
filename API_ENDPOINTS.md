# Endpoints API Ndugu Backend

Ce document décrit tous les endpoints disponibles dans l'API Ndugu Backend avec les services Ory intégrés.

## 🌐 Accès aux services

### API Gateway (APISIX)
- **URL** : http://localhost:9080
- **Fonction** : Point d'entrée principal pour toutes les requêtes

### Services directs
- **gRPC** : localhost:50051
- **HTTP** : localhost:8080

## 📡 Endpoints gRPC

### Service AuthService

#### CreateUser
- **Méthode** : `ndugu.v1.AuthService/CreateUser`
- **Description** : Crée un nouvel utilisateur via Ory Kratos
- **Request** :
  ```json
  {
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
  ```
- **Response** :
  ```json
  {
    "userId": "uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "createdAt": "2024-01-01T00:00:00Z"
  }
  ```

#### GetUser
- **Méthode** : `ndugu.v1.AuthService/GetUser`
- **Description** : Récupère un utilisateur par son ID
- **Request** :
  ```json
  {
    "userId": "uuid"
  }
  ```

#### ValidateSession
- **Méthode** : `ndugu.v1.AuthService/ValidateSession`
- **Description** : Valide une session Kratos
- **Request** :
  ```json
  {
    "sessionToken": "session_token_here"
  }
  ```

#### CreateOAuth2Client
- **Méthode** : `ndugu.v1.AuthService/CreateOAuth2Client`
- **Description** : Crée un client OAuth2 via Ory Hydra (temporairement désactivé)
- **Status** : ⚠️ En développement

#### CreatePermission
- **Méthode** : `ndugu.v1.AuthService/CreatePermission`
- **Description** : Crée une permission via Ory Keto (temporairement désactivé)
- **Status** : ⚠️ En développement

#### CheckPermission
- **Méthode** : `ndugu.v1.AuthService/CheckPermission`
- **Description** : Vérifie une permission via Ory Keto (temporairement désactivé)
- **Status** : ⚠️ En développement

## 🌐 Endpoints HTTP REST

### Utilisateurs

#### Créer un utilisateur
- **URL** : `POST /user`
- **Headers** : `Content-Type: application/json`
- **Body** :
  ```json
  {
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
  ```

#### Récupérer un utilisateur
- **URL** : `GET /user?id={userId}`
- **Response** :
  ```json
  {
    "id": "uuid",
    "email": "user@example.com",
    "name": {
      "first": "John",
      "last": "Doe"
    },
    "traits": {...},
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
  ```

### Santé des services

#### Vérifier l'état
- **URL** : `GET /health`
- **Response** : `Services Ory opérationnels`

## 🔧 Services Ory

### Ory Kratos (Gestion d'identité)
- **Public API** : http://localhost:4433
- **Admin API** : http://localhost:4434
- **Fonctionnalités** : Authentification, inscription, gestion des sessions

### Ory Hydra (OAuth2/OpenID Connect)
- **Public API** : http://localhost:4444
- **Admin API** : http://localhost:4445
- **Fonctionnalités** : OAuth2, OpenID Connect (en développement)

### Ory Keto (Contrôle d'accès)
- **Read API** : http://localhost:4466
- **Write API** : http://localhost:4467
- **Fonctionnalités** : Permissions, contrôle d'accès (en développement)

## 🚀 Exemples d'utilisation

### Via API Gateway (APISIX)

```bash
# Créer un utilisateur via APISIX
curl -X POST http://localhost:9080/api/user \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User"
  }'

# Vérifier la santé via APISIX
curl http://localhost:9080/health
```

### Via HTTP direct

```bash
# Créer un utilisateur
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User"
  }'

# Récupérer un utilisateur
curl "http://localhost:8080/user?id=USER_ID_HERE"
```

### Via gRPC (avec grpcurl)

```bash
# Installer grpcurl si nécessaire
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Créer un utilisateur via gRPC
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "firstName": "Test",
  "lastName": "User"
}' localhost:50051 ndugu.v1.AuthService/CreateUser
```

## 📝 Notes importantes

1. **Services temporairement désactivés** : Hydra et Keto sont configurés mais leurs fonctionnalités sont temporairement désactivées en attendant la résolution des problèmes de versions.

2. **CORS** : L'API Gateway APISIX est configuré avec CORS pour permettre les requêtes cross-origin.

3. **Sécurité** : Les secrets par défaut sont uniquement pour le développement. Changez-les pour la production.

4. **Base de données** : Tous les services Ory utilisent la même base de données PostgreSQL.

## 🔄 Prochaines étapes

1. Résoudre les problèmes de versions pour Hydra et Keto
2. Implémenter les fonctionnalités OAuth2 complètes
3. Ajouter le contrôle d'accès avec Keto
4. Tests d'intégration complets
5. Documentation des flux d'authentification
