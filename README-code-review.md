# Code Review et Points d'Amélioration

Ce document présente une auto review du projet, je précise que le but du projet était plus de se familiariser avec
le framework Gin/Gorm/Chi plus que de faire quelque chose de parfaitement clean, ni de faire quelque chose de
fonctionnellement complet et ou pertinent.

Cette review n'a pas la prétention d'être exhaustive, mais de mettre en avant quelques points d'améliorations auquels j'ai pensé
rapidement.

**Nommage**

Les nommage des paramètre/variable ne sont pas top je trouve (une ou deux lettre souvent), je ne trouve pas ça très propre
mais je vois beaucoup d'exemple en go qui sont comme ça j'ai l'impression que ça fait partie de la "philosophie" go
mais peut être est-ce une fausse impression

**Gestion des erreurs**
Pas d'exception handler global (on a des RestExceptionHandler en spring je n'ai pas trouvé d'équivalent simple pour le moment)

**Validation métier manquante**


**Securité**
- Aucune proteciton par rôle des path api

**Observabilité**

- Absence de contexte métier dans les logs
- Pas de metrics (otel ou autre), ainsi que d'actuator (health check / liveness / readiness)


**Messaging**
- Gestion simpliste des erreurs de publication
- Absence de Dead Letter Queue pour les messages en échec
- Pas de retry avec backoff exponentiel
- Un seul consumer pas très scalable (en java spring l'autoconf pour avoir plusieurs consumer est quasi automatique)

**Tests**
- Coverage largement insuffisant
- Particuliérement niveau IT sur le messaging (j'ai l'habitude d'utiliser testcontainers pour ça)


**API**
- Pas de documentation OpenAPI/Swagger
- Pas de pagination pour le get all des playlist
