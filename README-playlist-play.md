# Fonctionnalité de Lecture de Playlist

Cette documentation explique comment tester la nouvelle fonctionnalité de lecture de playlist qui envoie des événements RabbitMQ.

## Architecture

La fonctionnalité suit le pattern suivant :
1. **HTTP POST** `/playlists/{id}/play` - Déclenche la lecture d'une playlist
2. **Publisher RabbitMQ** - Envoie un événement `TrackPlayedEvent` pour chaque track de la playlist
3. **Consumer RabbitMQ** - Consomme les événements et les persiste dans la table `TrackPlay`
4. **Base de données** - Stocke l'historique des lectures pour les statistiques futures

## Démarrage des Services

### 1. Lancer RabbitMQ avec Docker Compose

```bash
docker-compose up -d rabbitmq
```

L'interface de gestion RabbitMQ sera accessible sur : http://localhost:15672
- Username: `radioking`  
- Password: `radioking123`

### 2. Lancer l'application

```bash
go run cmd/main.go
```

L'application :
- Se connecte à RabbitMQ
- Démarre le consumer en arrière-plan
- Lance le serveur HTTP sur le port 8080

## Test de la Fonctionnalité

### 1. Créer une playlist avec des tracks

```bash
curl -X POST http://localhost:8080/playlists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${JWT}"\
  -d '{
    "name": "Ma Playlist Test",
    "description": "Playlist pour tester la fonctionnalité de lecture",
    "tracks": [
      {
        "title": "Bohemian Rhapsody",
        "artist": "Queen"
      },
      {
        "title": "Hotel California", 
        "artist": "Eagles"
      },
      {
        "title": "Stairway to Heaven",
        "artist": "Led Zeppelin"
      }
    ]
  }'
```

Il faut egalement ajouter le bearer token d'authentification keycloak dans le header Authorization

### 2. Jouer la playlist

```bash
curl -X POST http://localhost:8080/playlists/1/play
```

Response attendue :
```json
{
  "message": "Playlist is being played",
  "playlist_id": 1,
  "tracks_count": 3
}
```


### 4. Vérifier dans RabbitMQ Management UI

1. Aller sur http://localhost:15672
2. Se connecter avec `radioking` / `radioking123`
3. Vérifier que l'exchange `playlist_events` existe
4. Vérifier que la queue `track_played` existe et traite les messages

## Structure des Données

### Événement RabbitMQ (TrackPlayedEvent)
```json
{
  "playlist_id": 1,
  "track_id": 1,
  "track_title": "Bohemian Rhapsody",
  "artist": "Queen", 
  "position": 0,
  "played_at": "2024-01-15T10:30:00Z",
  "event_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

## Configuration

La configuration RabbitMQ se trouve dans `config.yaml` :
```yaml
messaging:
  rabbitmq:
    url: "amqp://radioking:radioking123@localhost:5672"
    exchange: "playlist_events"
    queue: "track_played"
    routing_key: "track.played"
```

Elle peut être surchargée par les variables d'environnement :
- `MESSAGING_RABBITMQ_URL`
- `MESSAGING_RABBITMQ_EXCHANGE` 
- `MESSAGING_RABBITMQ_QUEUE`
- `MESSAGING_RABBITMQ_ROUTING_KEY`


## Authentification

L'auth ce fait via un jwt keycloak. La configuration du realm keycloak n'est pas poussé en l'état (je vais essayer
d'ajouter cela dans la semaine si j'ai le temps).


## DB

Je n'ai mis que du sqllite pour la db pour l'instant si j'ai le temps je mettreai un mariadb dans la semaine