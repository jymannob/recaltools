# Recaltools

> **⚠️ Projet archivé — n'est plus maintenu.**
> Depuis **Recalbox 10**, EmulationStation enregistre les données utilisateur (favoris, playcount, lastplayed, hidden…) dans un fichier séparé `gamelist-userdata.ini` au lieu de `gamelist.xml` : un re-scrap ne les efface plus, ce qui rend `backup`/`restore` obsolètes.
> L'outil reste utilisable pour les versions antérieures à Recalbox 10. Les derniers correctifs (bugs, écritures atomiques, Go 1.24) sont sur la branche `fix-known-issues`, non fusionnée.

Set of tools for recalbox
* `backup` save gamelists user metadatas (favorite, playcount, ...)
* `restore` apply metadatas saved by `backup` command to gamelists
* `clean` delete all scraping data and rename all `gamelist.xml`
* (**todo**) remove only media data referenced by `gamelist.xml`

## developement Usage

For backup gamelists metadatas
```bash
make tool backup <path_to_roms_directory>...
```

For restore gamelists metadatas
```bash
make tool restore <path_to_roms_directory>...
```

Show help
```bash
make tool
```

## Build

For build binary
```bash
PROJECT_VERSION=<Version> make build
```
## TODO 
* un package parent ne doit pas être dépendant d'un sous package
* passer les packages non publique en `internal`
* gérer l'execution des goroutines de preference dans des fonctions annonymes
* limiter le nombre de goroutines concurentes (channels ?, workers ?, semaphores ?)
* integrer du perf profiling (pprof ?, "runtime/trace")
* utiliser un bufferio pour ecrir les logs ou dans les fichhiers

### Problèmes repérés (à traiter)
* **`restore` peut paniquer** (`gamelists_backup.go`, `restoreSystem`) : `xmlquery.Query` renvoie `nil, nil` si le jeu n'est pas dans le XML → `xml.ReplaceChildNode(nil, …)`. Le `return` sur erreur abandonne aussi tout le système sans écrire le fichier.
* **XPath `contains()` sur le path** (`gameFindXpath`) : peut matcher le mauvais jeu (path inclus dans un autre) et casse si le path contient un `"`. Préférer une égalité stricte.
* **Plusieurs `RomsDir`** : `Gamelists` n'est jamais réinitialisé entre deux dossiers → les gamelists du 1er dossier sont retraitées (`backup`, `restore`, `clean` ; `clean` tente de re-déplacer un fichier déjà déplacé).
* `backupSystem` : `xmlquery.QueryAll` exécuté deux fois, et l'erreur de `utils.WriteJsonFile` n'est pas récupérée.
* Commande `clean` sans tests.
* `.goreleaser.yaml` : `buildVersion` codé en dur à `development` (utiliser `{{.Version}}`).
* `cmd/gamelist-backup` (ancien binaire `gamelistBkp`) toujours buildé par le Makefile : à supprimer ?
* Monter la version de Go (1.17) et les dépendances.
## Roadmap

Objectif : une GUI utilisable dans 3 contextes (sur un PC avec `share` monté en réseau, sur la Recalbox via un navigateur, en TUI via SSH), sans multiplier les interfaces.

```
                ┌────────────────────────────────┐
                │ core (Go, sans I/O d'affichage)│  backup / restore / clean
                │ résultats structurés, events   │  dry-run, progression
                └──────┬──────────────┬──────────┘
                       │              │
             ┌─────────┴────┐    ┌────┴───────┐
             │ API HTTP/JSON│    │ CLI + TUI  │  (appels directs au core)
             └──┬────────┬──┘    └────────────┘
                │        │
   Web UI (go:embed)   Desktop (Wails : même front web dans une fenêtre native)
```

1. **Corriger les bugs connus** (section « Problèmes repérés »).
2. **Extraire le cœur en bibliothèque** (`internal/core`) : résultats structurés et erreurs remontées au lieu de `log.Println`, événements de progression, mode dry-run pour `clean` et `restore`. Le CLI devient un client du core.
3. **API HTTP/JSON + front web** embarqué via `go:embed`, servi par le binaire (Recalbox ou PC).
4. **TUI** (bubbletea) branchée directement sur le core, puis éventuellement packaging desktop avec **Wails**.

Points d'attention :
* Sécurité de l'API web sur la Recalbox (écoute limitée, jeton d'accès) : elle peut supprimer des fichiers.
* Accès concurrents avec EmulationStation qui réécrit les gamelists.
* Un seul binaire statique sans dépendance, compatible Raspberry Pi.
