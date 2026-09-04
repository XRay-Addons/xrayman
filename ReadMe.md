# XRay Node

Tool to be running on node

## Install xray node

### .deb (Ubuntu, Debian, etc...)

```sh
# download .deb for required architecture (amd64 or arm64)
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/xray-node-amd64.deb" -o xray-node.deb
# install
sudo apt install ./xray-node.deb
# view setup manual
sudo apt show xray-node
# start
sudo systemctl daemon-reload
sudo systemctl enable --now xray-node
# or restart
sudo systemctl daemon-reload
sudo systemctl restart xray-node
# view logs
sudo journalctl -u xray-node -n 50 -f
```

If success, user following handlers are available:

- `${ENDPOINT}/api` - node api
- `${ENDPOINT}/api/version` - unprotected status handler

### Manual build

#### Requirements

- go v1.26.2

#### Build

```sh
# download source code
git clone https://github.com/XRay-Addons/xrayman.git
cd xrayman
git submodule update --init --recursive
# install build tools
make tools
# build
make build_node
# check
./build/xray-node/xray-node -h
```

# XRay Node Manager

Service to manage all nodes, users, subscriptions, etc...

## Install xray nodeman

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

- `${ENDPOINT}/u` - user page
- `${ENDPOINT}adm` - admin page
- `${ENDPOINT}/api/version` - unprotected status handler
- `${METRICS_ENDPOINT}/metrics` - prometheus metrics

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

## Pretty logs view

To view pretty human-readable logs use `jq` preset:

### Install JQ

JQ install page: https://jqlang.org/download/

### Add alias for go zap logs viewer

#### OSX

Open `~/.zshrc`

```sh
sudo nano ~/.zshrc
```

Add `zl` - alias for formatting zap logs

```sh
# >>> zl - jq setup for go zap log formatting >>>
zl() {
jq -r '
  def level_color:
    if . == "error" then "\u001b[31m"
    elif . == "warn" then "\u001b[33m"
    elif . == "info" then "\u001b[32m"
    elif . == "debug" then "\u001b[36m"
    else "\u001b[0m"
    end;

  . as $log
  |
  "\($log.ts) \($log.level | level_color)\($log.level | ascii_upcase)\u001b[0m \($log.msg)",
  ($log
    | to_entries[]
    | select(.key != "ts" and .key != "level" and .key != "msg")
    | if (.value | type) == "string" and (.value | contains("\n")) then
        "- \u001b[36m\(.key)\u001b[0m:\n\(.value | split("\n") | map("    " + .) | join("\n"))"
      else
        "- \u001b[36m\(.key)\u001b[0m: \(.value)"
      end
  )
'
}
# <<< zlog initializing <<<
```

Apply modifications

```sh
source ~/.zshrc
```

### Use `zl` alias

Just add `| zl` to any command returning logs, like

```sh
docker compose logs -n 50 -f xray-nodeman | zl
```
