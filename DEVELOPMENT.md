# Development Guide

## Changing the `README.md`

Our `README.md` is generated.
Manual modifications in the `README.md` file will be overwritten.

If you want to change the `README.md` file, please make your changes to `assets/README.template`.

## Extending the *yml* format in `/games`

Game files in `/games` may use either the `.yml` or `.yaml` extension; both are loaded by the tooling.

If you aim to modify (add, change, or remove) a YAML property in `/games/*.yml` (or `/games/*.yaml`), please ensure that you make this change in all YAML files.

Our tooling in `/app` also needs adjustment.
Mainly in

* `/app/cmd/types.go`: Adjusting the type structure
* `app/cmd/convertYamlToJson.go` -> `mergeGameInformation`: Adjusting the YAML to JSON merge logic