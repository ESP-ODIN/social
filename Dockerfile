# Dockerfile de développement pour le service API (Go)
# Objectif : lancer l'API avec hot-reload (via air) dans un conteneur,
# pour être branché tel quel dans le docker-compose global de l'archi micro-service.

FROM golang:1.27-alpine

WORKDIR /app

# air : hot-reload pour Go (rebuild + restart au changement de fichier),
# l'équivalent de `next dev` côté Node.
RUN go install github.com/air-verse/air@latest

# On copie d'abord uniquement les manifests pour profiter du cache Docker :
# tant que go.mod / go.sum ne changent pas, ce layer n'est pas rejoué.
COPY go.mod go.sum ./
RUN go mod download

# Le reste du code sera de toute façon écrasé par le bind mount en dev,
# mais on le copie pour que l'image soit utilisable seule (ex: premier build, CI).
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
