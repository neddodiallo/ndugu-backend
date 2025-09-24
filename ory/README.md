# Configuration Ory pour Ndugu Backend

Ce dossier contient la configuration pour les services Ory utilisés dans le projet Ndugu Backend.

## Services Ory

### Ory Kratos - Gestion d'identité
- **Port public**: 4433
- **Port admin**: 4434
- **Configuration**: `kratos/kratos.yml`
- **Schéma d'identité**: `kratos/identity.schema.json`

Kratos gère l'authentification et l'inscription des utilisateurs.

### Ory Hydra - OAuth2 et OpenID Connect
- **Port public**: 4444
- **Port admin**: 4445
- **Configuration**: `hydra/hydra.yml`

Hydra gère l'autorisation OAuth2 et OpenID Connect.

### Ory Keto - Contrôle d'accès
- **Port lecture**: 4466
- **Port écriture**: 4467
- **Configuration**: `keto/keto.yml`

Keto gère les permissions et le contrôle d'accès basé sur les attributs.

## Démarrage des services

```bash
# Démarrer tous les services
docker-compose up -d

# Vérifier le statut des services
docker-compose ps

# Voir les logs
docker-compose logs kratos
docker-compose logs hydra
docker-compose logs keto
```

## URLs importantes

- **Kratos Public API**: http://localhost:4433
- **Kratos Admin API**: http://localhost:4434
- **Hydra Public API**: http://localhost:4444
- **Hydra Admin API**: http://localhost:4445
- **Keto Read API**: http://localhost:4466
- **Keto Write API**: http://localhost:4467

## Configuration de base de données

Tous les services Ory utilisent la même base de données PostgreSQL configurée dans `docker-compose.yml`. Les tables sont créées automatiquement au premier démarrage.

## Sécurité

⚠️ **Important**: Les secrets par défaut dans les fichiers de configuration sont à des fins de développement uniquement. Pour la production, changez tous les secrets :

- `SECRETS_SYSTEM` dans Hydra
- `secrets.cookie` et `secrets.cipher` dans Kratos
- `salt` dans la configuration OIDC de Hydra

## Intégration avec l'application

Le module `internal/auth/auth.go` fournit une interface Go pour interagir avec tous les services Ory. Voir les exemples d'utilisation dans le code.

