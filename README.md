# Social

API Go minimale, sans dependance externe. Prerequis : Go 1.27.0 ou plus recent.

## Demarrage local (PowerShell)

```powershell
go run ./cmd/api
```

Dans un second terminal :

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
```

La route retourne HTTP 200 et `{"status":"ok"}`. Arreter avec Ctrl+C.
Cette route verifie que le serveur HTTP repond ; aucune base de donnees n'est encore connectee.

## Configuration

Par defaut, le serveur ecoute sur `127.0.0.1:8080`. Pour changer le port :

```powershell
$env:HTTP_ADDR = "127.0.0.1:8081"
go run ./cmd/api
```

`.env.exemple` documente les variables. Les fichiers `.env` ne sont pas charges automatiquement.
La configuration est lue dans l'environnement du processus.

## Verification et compilation

```powershell
go test ./...
go vet ./...
go build -o bin/ ./cmd/...
```

Le binaire de l'API sous Windows est `bin/api.exe`.
Si Make est installe, les commandes `make run`, `make test`, `make vet`,
`make build` et `make fmt` sont egalement disponibles.

## Structure

- `cmd/api` : serveur HTTP et routes.
- `config` : lecture et validation de la configuration.
- `db` : emplacement reserve a l'acces aux donnees.
- `cmd/migrate` : commande reservee aux migrations ; elle retourne une erreur explicite tant que la base n'est pas configuree.
