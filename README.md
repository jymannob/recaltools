# Recaltools

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
* `.goreleaser.yaml` : `buildVersion` codé en dur à `development` (utiliser `{{.Version}}`).
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
