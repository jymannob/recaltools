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
* **`restore` peut paniquer** (`gamelists_backup.go`, `restoreSystem`) : `xmlquery.Query` renvoie `nil, nil` si le jeu n'est pas dans le XML → `xml.ReplaceChildNode(nil, …)`. Le `return` sur erreur abandonne aussi tout le système sans écrire le fichier.
* **XPath `contains()` sur le path** (`gameFindXpath`) : peut matcher le mauvais jeu (path inclus dans un autre) et casse si le path contient un `"`. Préférer une égalité stricte.
* **Plusieurs `RomsDir`** : `Gamelists` n'est jamais réinitialisé entre deux dossiers → les gamelists du 1er dossier sont retraitées (`backup`, `restore`, `clean` ; `clean` tente de re-déplacer un fichier déjà déplacé).
* `backupSystem` : `xmlquery.QueryAll` exécuté deux fois, et l'erreur de `utils.WriteJsonFile` n'est pas récupérée.
* Commande `clean` sans tests.
* `.goreleaser.yaml` : `buildVersion` codé en dur à `development` (utiliser `{{.Version}}`).
* `cmd/gamelist-backup` (ancien binaire `gamelistBkp`) toujours buildé par le Makefile : à supprimer ?
* Monter la version de Go (1.17) et les dépendances.