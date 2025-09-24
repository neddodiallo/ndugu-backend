# Migrations Ory - Guide d'Utilisation

Ce document explique comment utiliser les commandes de migration pour les services Ory (Kratos, Hydra, Keto) dans le projet Ndugu Backend.

## Commandes Disponibles

### Migrations Individuelles

#### Kratos (Gestion d'Identités)
```bash
make migrate-kratos
```
- Applique les migrations de base de données pour Kratos
- Crée les tables nécessaires pour la gestion d'identités
- Utilise la commande : `kratos migrate sql -e --yes`

#### Hydra (OAuth2 et OpenID Connect)
```bash
make migrate-hydra
```
- Applique les migrations de base de données pour Hydra
- Crée les tables pour OAuth2, clients, consentements, etc.
- Utilise la commande : `hydra migrate sql -e --yes`

#### Keto (Contrôle d'Accès)
```bash
make migrate-keto
```
- Applique les migrations de base de données pour Keto
- Crée les tables pour les permissions et relations
- Utilise la commande : `keto migrate up --yes -c /etc/config/keto/keto.yml`

### Migration Complète

#### Toutes les Migrations
```bash
make migrate-all
```
- Applique toutes les migrations Ory en une seule commande
- Exécute dans l'ordre : Kratos → Hydra → Keto
- Idéal pour l'initialisation complète de l'environnement

## Prérequis

1. **Base de données PostgreSQL** : Le service `db` doit être en cours d'exécution
2. **Docker Compose** : Les services Ory doivent être configurés dans `docker-compose.yml`
3. **Fichiers de configuration** : Les fichiers de config doivent être présents dans `ory/`

## Utilisation

### Initialisation d'un Nouvel Environnement

```bash
# 1. Démarrer la base de données
docker-compose up -d db

# 2. Appliquer toutes les migrations
make migrate-all

# 3. Démarrer tous les services
docker-compose up -d
```

### Après une Mise à Jour

```bash
# Appliquer les nouvelles migrations
make migrate-all

# Redémarrer les services si nécessaire
docker-compose restart kratos hydra keto
```

### Vérification des Migrations

```bash
# Vérifier l'état des migrations Kratos
docker-compose run --rm kratos migrate sql -e --help

# Vérifier l'état des migrations Hydra
docker-compose run --rm hydra migrate sql -e --help

# Vérifier l'état des migrations Keto
docker-compose run --rm keto migrate status -c /etc/config/keto/keto.yml
```

## Dépannage

### Erreur "Unable to locate the table"
- Les migrations n'ont pas été appliquées
- Solution : Exécuter `make migrate-all`

### Erreur "relation already exists"
- Les tables existent déjà mais sont dans un état incohérent
- Solution : Supprimer les tables et réappliquer les migrations

### Erreur de connexion à la base de données
- Vérifier que le service `db` est en cours d'exécution
- Vérifier les paramètres de connexion dans les fichiers de configuration

## Structure des Migrations

### Kratos
- Tables d'identités et sessions
- Schémas d'identités personnalisés
- Flows d'authentification

### Hydra
- Tables OAuth2 (clients, tokens, consentements)
- Tables JWK (clés JSON Web)
- Tables de gestion des flows

### Keto
- Tables de relations et permissions
- Tables de mapping UUID
- Tables de namespaces

## Notes Importantes

1. **Ordre des Migrations** : Toujours appliquer les migrations dans l'ordre Kratos → Hydra → Keto
2. **Sauvegarde** : Toujours sauvegarder la base de données avant d'appliquer des migrations en production
3. **Environnement de Développement** : Ces commandes sont optimisées pour l'environnement de développement
4. **Production** : Adapter les commandes pour l'environnement de production (variables d'environnement, secrets, etc.)

## Support

Pour plus d'informations sur les migrations Ory :
- [Documentation Kratos](https://www.ory.sh/docs/kratos/)
- [Documentation Hydra](https://www.ory.sh/docs/hydra/)
- [Documentation Keto](https://www.ory.sh/docs/keto/)
