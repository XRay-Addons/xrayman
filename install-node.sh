#!/usr/bin/env bash
set -euo pipefail

STEP_STATUS=()

ok()   { echo "✔ $1"; STEP_STATUS+=("✔ $1"); }
warn() { echo "⚠ $1"; STEP_STATUS+=("⚠ $1"); }
fail() { echo "✖ $1" >&2; exit 1; }

########## 0. PRE-FLIGHT ##########

[[ "$EUID" -eq 0 ]] || fail "Run as root"
ok "Running as root"

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

install_pkg() {
  local pkg="$1"

  if command_exists apt-get; then
    apt-get update -y >/dev/null
    apt-get install -y "$pkg" >/dev/null
  elif command_exists dnf; then
    dnf install -y "$pkg" >/dev/null
  elif command_exists yum; then
    yum install -y "$pkg" >/dev/null
  else
    fail "No package manager"
  fi
}

for dep in curl unzip logrotate; do
  command_exists "$dep" || install_pkg "$dep"
  ok "dependency: $dep"
done

umask 002

########## 1. USER / GROUP ##########

if ! getent group xray-node >/dev/null; then
  groupadd xray-node
fi

if ! id xray-node >/dev/null 2>&1; then
  useradd --system --no-create-home --shell /usr/sbin/nologin \
    --gid xray-node xray-node
fi

INSTALL_USER="${SUDO_USER:-$USER}"
if id "$INSTALL_USER" &>/dev/null; then
  usermod -aG xray-node "$INSTALL_USER" || true
  ok "user added to group"
fi

########## 2. PATHS ##########

BIN_DIR="/opt/xray-node"
DATA_DIR="$BIN_DIR/data"

CONFIG_DIR="/etc/xray-node"
PERSIST_DIR="/var/lib/xray-node"
LOG_DIR="/var/log/xray-node"

# USER-WORKSPACE (FIXED ROOT ISSUE)
LNK_DEFAULT="${HOME}/xraynode-links"

mkdir -p "$BIN_DIR" "$DATA_DIR" "$CONFIG_DIR" "$PERSIST_DIR" "$LOG_DIR"

########## 3. DOWNLOAD ##########

TMP_DIR="$(mktemp -d /tmp/xray-node-install-XXXX)"
trap 'rm -rf "$TMP_DIR"' EXIT

ZIP_URL="https://github.com/XRay-Addons/xrayman/releases/latest/download/xrayman.zip"
ZIP_FILE="$TMP_DIR/xray.zip"

HTTP_CODE=$(curl -L -s -w "%{http_code}" -o "$ZIP_FILE" "$ZIP_URL")

[[ "$HTTP_CODE" -eq 200 ]] || fail "download failed"
[[ -s "$ZIP_FILE" ]] || fail "empty archive"

ok "downloaded"

unzip -q "$ZIP_FILE" -d "$TMP_DIR"

########## 4. INSTALL BINARIES ##########

install -o xray-node -g xray-node -m 755 \
  "$TMP_DIR/xray-node/xray-node" \
  "$BIN_DIR/xray-node"

install -o xray-node -g xray-node -m 644 \
  "$TMP_DIR/xray-node/data/"*.dat \
  "$DATA_DIR/"

ok "binaries installed"

########## 5. OWNERSHIP MODEL ##########

chown -R xray-node:xray-node \
  "$BIN_DIR" \
  "$DATA_DIR" \
  "$CONFIG_DIR" \
  "$PERSIST_DIR" \
  "$LOG_DIR"

chmod 755 "$BIN_DIR"
chmod 2775 "$CONFIG_DIR"
chmod 2775 "$PERSIST_DIR"
chmod 2775 "$LOG_DIR"

ok "ownership fixed"

########## 6. LOGROTATE ##########

cat > /etc/logrotate.d/xray-node <<EOF
$LOG_DIR/*.log {
  size 100M
  rotate 5
  compress
  missingok
  notifempty
  copytruncate
}
EOF

ok "logrotate configured"

########## 7. CONFIGS ##########

CONFIG_SERVER_CONTENT=$(cat <<EOF
{
  "log": {
    "loglevel": "warning",
    "dnsLog": false,
    "access": "$LOG_DIR/xray-access.log",
    "error": "$LOG_DIR/xray-error.log"
  },
  "dns": {
    "servers": ["https://1.1.1.1/dns-query"],
    "queryStrategy": "UseIP"
  },
  "api": {
    "tag": "api",
    "listen": "127.0.0.1:32999",
    "services": ["HandlerService", "LoggerService", "StatsService", "ReflectionService"]
  },
  "stats": {},

  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true,
        "statsUserOnline": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true,
      "statsOutboundUplink": true,
      "statsOutboundDownlink": true
    }
  },
  "inbounds": [],
  "outbounds": [
    {
      "tag": "direct",
      "protocol": "freedom"
    },
    /* uncomment for cloudflare warp outbound
    {
      "tag": "warp",
      "protocol": "wireguard",
      "settings": {
        # paste data from generated WARP config, IPv4 or v6 depends on your server
        "secretKey": "#################",
        "DNS": ["####","####","####"],
        "address": ["####","########"],
        "peers": [
          {
            "publicKey": "#################",
            "endpoint": "#################",
          }
        ],
        "mtu": ####,
        "domainStrategy": "ForceIP",
        "kernelMode": false
      }
    },*/
    {
      "tag": "blacklist",
      "protocol": "blackhole"
    }
  ],
  "routing": {
    "rules": [],
    "domainStrategy": "IPIfNonMatch"
  }
}
EOF
)

CONFIG_CLIENT_CONTENT='[]'

create_config() {
  local file="$1"
  local content="$2"

  if [[ -f "$file" ]]; then
    read -rp "overwrite $file? YES: " ans
    [[ "$ans" == "YES" ]] || { warn "skip $file"; return; }
  fi

  printf "%s\n" "$content" > "$file"
  chown xray-node:xray-node "$file"
  chmod 664 "$file"

  ok "config: $file"
}

create_config "$CONFIG_DIR/xray_server.json" "$CONFIG_SERVER_CONTENT"
create_config "$CONFIG_DIR/xray_client.json" "$CONFIG_CLIENT_CONTENT"

########## 8. WARP ##########

read -rp "warp? (NO skip): " WARP

if [[ "${WARP^^}" != "NO" ]]; then
  WARP_BIN="$TMP_DIR/wgcf"

  curl -fsSL \
    "https://github.com/ViRb3/wgcf/releases/download/v2.2.31/wgcf_2.2.31_linux_amd64" \
    -o "$WARP_BIN"

  chmod +x "$WARP_BIN"

  yes "Yes" | timeout 60s "$WARP_BIN" register || true
  timeout 30s "$WARP_BIN" generate || true

  if [[ -f wgcf-profile.conf ]]; then
    mv wgcf-profile.conf "$CONFIG_DIR/warp.conf"
    chown xray-node:xray-node "$CONFIG_DIR/warp.conf"
    chmod 664 "$CONFIG_DIR/warp.conf"
    ok "warp generated"
  else
    warn "warp skipped"
  fi
fi

########## 9. PORT ##########

read -rp "port [5001]: " PORT
PORT="${PORT:-5001}"

[[ "$PORT" =~ ^[0-9]+$ ]] || fail "bad port"

########## 10. SYSTEMD ##########

SERVICE_FILE="/etc/systemd/system/xray-node.service"

cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=xray-node
After=network.target

[Service]
User=xray-node
Group=xray-node
ExecStart=$BIN_DIR/xray-node -a :$PORT -c $CONFIG_DIR -d $DATA_DIR -p $PERSIST_DIR
Restart=always

StandardOutput=append:$LOG_DIR/out.log
StandardError=append:$LOG_DIR/err.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable xray-node >/dev/null

ok "systemd ready"

########## 11. SYMLINKS (FIXED ROOT PROBLEM) ##########

read -rp "symlink dir (empty = skip): " LNK

if [[ -z "$LNK" ]]; then
  warn "symlinks skipped"
else

  INSTALL_USER="${SUDO_USER:-$USER}"

  # create directory as USER (not root ownership)
  mkdir -p "$LNK"
  chown "$INSTALL_USER":"$INSTALL_USER" "$LNK"
  chmod 755 "$LNK"

  run_as_user() {
    sudo -u "$INSTALL_USER" -H bash -c "$1"
  }

  run_as_user "
    ln -sfn '$BIN_DIR' '$LNK/xray-node'
    ln -sfn '$CONFIG_DIR' '$LNK/config'
    ln -sfn '$LOG_DIR' '$LNK/logs'
    ln -sfn '$PERSIST_DIR' '$LNK/persist'
    ln -sfn '$SERVICE_FILE' '$LNK/service'
  "

  ok "symlinks created as user (correct ownership)"
fi

########## DONE ##########

echo "=== SUMMARY ==="
printf '%s\n' "${STEP_STATUS[@]}"
