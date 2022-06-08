# Recaltools

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
