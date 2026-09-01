## [NO] One-command Install xray-node

```
curl -fsSL https://raw.githubusercontent.com/XRay-Addons/xrayman/main/install-node.sh -o install-node.sh && sudo bash install-node.sh
```

## [NO] One-command Install xray-nodeman

```
curl -fsSL https://raw.githubusercontent.com/XRay-Addons/xrayman/main/install-nodeman.sh -o install-nodeman.sh && sudo bash install-nodeman.sh
```

# XRay Node

Tool to be running on node

# XRay Node Manager

Service to manage all nodes, users, subscriptions, etc...

## Install xray node manager

### .deb (Ubuntu, Debian, etc...)

#### Pre-requirements

Postgres DB required. You should install it as:

- standalone application
- docker image
- use third-party Postgres hosting

Only one thing you should provide to node manager is **postgress connection string**:

```
 postgresql://username:password@host:port/database_name
```

#### Installation

```sh
# download .deb for required architecture (amd64 or arm64)
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/xray-nodeman-amd64.deb" -o xray-nodeman.deb
# install
sudo apt install ./xray-nodeman.deb
# view setup manual
sudo apt show xray-nodeman
# start
sudo systemctl daemon-reload
sudo systemctl enable --now xray-nodeman
# or restart
sudo systemctl daemon-reload
sudo systemctl restart xray-nodeman
# view logs
sudo journalctl -u xray-nodeman -n 50 -f
```

If success, user following handlers are available:

- `localhost:8080/u` - user page
- `localhost:8080/adm` - admin page
- `localhost:8080/api/version` - unprotected status page
- `localhost:9090/metrics` - prometheus metrics

### Docker

#### Standalone

Docker container with xray-nodeman and nothing else.

Postgress DB and its conntection string should be provided as above.

```sh
# download docker-compose
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/docker-compose.standalone.tar.gz" | tar -xzf -
# create .env file
sudo nano .env
# view setup manual (add variables to .env till success)
docker compose run --rm xray-nodeman --help
# start
docker compose up -d
# or restart
docker compose down && docker compose up -d
# view logs
docker compose logs -n 50 -f xray-nodeman
```

#### All-in-one

Docker container contains everything required for nodeman:

- postgress db
- postgress db backuper
- xray node manager
- prometheus
- grafana

```sh
# download docker-compose
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/docker-compose.all-in-one.tar.gz" | tar -xzf -
# create .env file
sudo nano .env
# view setup manual (add variables to .env till success)
docker compose run --rm xray-nodeman --help
# start
docker compose up -d
# or restart
docker compose down && docker compose up -d
# view logs
docker compose logs -n 50 -f xray-nodeman
```

### Manual build

#### Requirements

- go v1.26.2
- pnpm v10.33.2
- node.js v24.4.0

#### Build

```sh
# download source code
git clone https://github.com/XRay-Addons/xrayman.git
cd xrayman
git submodule update --init --recursive
# install build tools
make tools
# build
make build_nodeman
# check
./build/xray-nodeman/xray-nodeman -h
```
