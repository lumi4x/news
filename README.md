# Varso API

![CI/CD](https://github.com/lumi4x/news/actions/workflows/.github/workflows/cicd.yaml/badge.svg)

## Prerequisites

1. Add `export COMPOSE_FILE="docker-compose.dev.yaml"` to your `~/.profile` file. Alternatively, add `-f` to your docker compose commands (e.g. `docker compose -f docker-compose.dev.yaml`).

2. Install `protoc`:

```
brew install protobuf
```

2. Run go generate to generate the DB schema and compile the protobuf files:

```
go generate ./...
```

3. Generate dev appropriate .env file:

```
./scripts/setup_dev_env.sh
```

4. For frontend development, enable this chrome flag: [chrome://flags/#allow-insecure-localhost](chrome://flags/#allow-insecure-localhost)

## Running the system

1. Run the system using the following command:

```
# Changes to nginx.conf require specifying --build as well.
docker compose up -d
```

2. Run the frontend

```
cd src/services/www
npm install
npm start
```

## Publishing the API NPM package

To release a new package of the API, simply increment the `package.json` version in `src/packages/varsoapi` and tag your commit. A GitHub Actions workflow will kick off, creating and pushing a new package.

## Useful commands

### Reset / clean database

If prompted for a password, enter 'devsecret'

```
docker compose stop app && dropdb -U postgres -h localhost -p 5432 app && docker compose up -d app && docker compose logs -f app
```

### Adding a new migration

```
go run -mod=mod ./src/services/auth/internal/ent/migrate/main.go add_name
```

### Adding a new service

1. Add it to all docker-compose files
2. Add it to cicd.yaml
3. Create the service in src/services/

### Running production images locally

1. Generate a GitHub access token: https://github.com/settings/tokens
2. Run `export NPM_TOKEN=<your_token>`
3. Log into GHCR by running: `echo $NPM_TOKEN | docker login ghcr.io -u USERNAME --password-stdin`
4. Use `-f docker-compose.yaml` in your `docker compose` commands to work with the production file.
5. Boot up the system with `docker compose -f docker-compose.yaml up -d`

### Setting up a production machine

Taken from: https://gist.github.com/pcsaini/178a806aab36ee6d31ac7bcb382cad79

```
sudo apt update
sudo apt upgrade
sudo apt install -y ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt install docker-ce docker-ce-cli containerd.io -y

echo "PasswordAuthentication no" >> /etc/ssh/sshd_config
sudo service ssh restart
```
