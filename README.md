# moc [![](https://images.microbadger.com/badges/image/ctfl/moc.svg)](https://hub.docker.com/r/ctfl/moc "DockerHub Image")

The Message Operation Center stores messages up to 160 characters long and delivers them to "clients" who deliver them to subscribers via services such as Twitter, IRC or e-mail.

## API

The complete API specifications are documented in the `swagger.yaml`. This can also be used to generate a client SDK with swagger. A simple import into Postman is also possible.

```bash
# Get messages

curl https://moc.example.com/messages
```

```bash
# Add message

curl -X POST \
	--header "Authorization: Bearer <operatorToken>" \
	--header "Content-Type: application/json" \
	--data '{"message": "Meine Nachricht"}' \
	https://moc.example.com/messages
```

## Configuration

### SQLite3 Config

```bash
DATABASE_DRIVER=sqlite3
DATABASE_PATH=test.sqlite3
```

### Operator Token

```bash
OPERATOR_TOKEN=1234
```

The shared secret with an operator for this microservice. Used to verify requests have been proxied through the operator and the payload values can be trusted.

## Usage

```
moc --migrate --seed
```

## Docker Compose

```yaml
version: '3'
services:
  moc:
    image: ctfl/moc
    env_file:
      - .env
    ports:
      - 80:80
```

## ToDo:

- Get Limit Messages & start date
- Prometheus Endpunkt
