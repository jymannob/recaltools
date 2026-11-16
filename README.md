# Recaltools

> **⚠️ Projet archivé — n'est plus maintenu.**
> Depuis **Recalbox 10**, EmulationStation enregistre les données utilisateur (favoris, playcount, lastplayed, hidden…) dans un fichier séparé `gamelist-userdata.ini` au lieu de `gamelist.xml` : un re-scrap ne les efface plus, ce qui rend `backup`/`restore` obsolètes.
> L'outil reste utilisable pour les versions antérieures à Recalbox 10. Les derniers correctifs (bugs, écritures atomiques, Go 1.24) sont sur la branche `fix-known-issues`, non fusionnée.

Set of tools for recalbox
* `backup` save gamelists user metadatas (favorite, playcount, ...)
* `restore` apply metadatas saved by `backup` command to gamelists
* (**todo**) `clean` delete all scraping data and rename all `gamelist.xml`

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
