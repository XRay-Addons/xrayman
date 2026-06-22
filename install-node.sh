#!/usr/bin/env bash
set -euo pipefail

STEP_STATUS=()

ok() {
  echo "✔ $1"
  STEP_STATUS+=("✔ $1")
}

warn() {
  echo "⚠ $1"
  STEP_STATUS+=("⚠ $1")
}

fail() {
  echo "✖ $1" >&2
  exit 1
}

########## 0. Pre-flight checks ##########

if [[ "$EUID" -ne 0 ]]; then
  fail "Script must be run with sudo or as root"
fi
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
    fail "Unsupported package manager (cannot install $pkg)"
  fi
}

for dep in unzip logrotate curl; do
  if ! command_exists "$dep"; then
    install_pkg "$dep"
    ok "Installed dependency: $dep"
  else
    ok "Dependency already present: $dep"
  fi
done

fixOwnership() {
  local path="$1"
  local user="${2:-xray-node}"
  local group="${3:-xray-node}"

  # skip if path does not exist
  if [[ ! -e "$path" ]]; then
    warn "fixOwnership: path does not exist: $path"
    return 1
  fi

  # set ownership if user/group exists
  if id "$user" &>/dev/null; then
    chown "$user:$group" "$path" 2>/dev/null || true
    ok "Ownership set: $user:$group -> $path"
  else
    warn "fixOwnership: user '$user' does not exist, skipped"
  fi
}

########## 1. Download latest release ##########
TMP_DIR="$(mktemp -d /tmp/xray-node-install-XXXX)"
trap 'rm -rf "$TMP_DIR"' EXIT

BASE_URL="https://github.com/XRay-Addons/xrayman"
ZIP_URL="$BASE_URL/releases/latest/download/xrayman.zip"

ZIP_FILE="$TMP_DIR/xray-node.zip"

# Follow redirects explicitly for clarity/debugging
HTTP_CODE=$(curl -L -s -w "%{http_code}" -o "$ZIP_FILE" "$ZIP_URL")

if [[ "$HTTP_CODE" -ne 200 ]]; then
  fail "Failed to download release archive (HTTP $HTTP_CODE)"
fi

if [[ ! -s "$ZIP_FILE" ]]; then
  fail "Downloaded file is empty"
fi

ok "Downloaded latest release archive"

########## 2. Unpack and validate archive ##########

unzip -q "$TMP_DIR/xray-node.zip" -d "$TMP_DIR"

if [[ ! -f "$TMP_DIR/xray-node/xray-node" ]] ||
   [[ ! -f "$TMP_DIR/xray-node/data/geoip.dat" ]] ||
   [[ ! -f "$TMP_DIR/xray-node/data/geosite.dat" ]]; then
  fail "Archive validation failed (missing required files)"
fi
ok "Archive structure validated"

########## 3. Install binaries and data ##########

BIN_DIR="/opt/xray-node"
DATA_DIR="$BIN_DIR/data"

mkdir -p "$DATA_DIR"

cp "$TMP_DIR/xray-node/xray-node" "$BIN_DIR/xray-node"
cp "$TMP_DIR/xray-node/data/"*.dat "$DATA_DIR/"

chmod 755 "$BIN_DIR/xray-node"
chmod 644 "$DATA_DIR/"*.dat

ok "Installed binary and data files to /opt/xray-node"

########## 4. User and group ##########

if ! getent group xray-node >/dev/null; then
  groupadd xray-node
  ok "Created group xray-node"
else
  warn "Group xray-node already exists"
fi

if ! id xray-node >/dev/null 2>&1; then
  useradd --system --no-create-home --shell /usr/sbin/nologin --gid xray-node xray-node
  ok "Created user xray-node"
else
  warn "User xray-node already exists"
fi

INSTALL_USER="${SUDO_USER:-$USER}"
if id "$INSTALL_USER" &>/dev/null; then
  usermod -aG xray-node "$INSTALL_USER"
  ok "Added user '$INSTALL_USER' to group xray-node"
  warn "You may need to re-login or run: newgrp xray-node"
else
  warn "Installer user '$INSTALL_USER' not found, skipping group add"
fi

########## 5. Standard directories ##########

CONFIG_DIR="/etc/xray-node"
PERSIST_DIR="/var/lib/xray-node"
LOG_DIR="/var/log/xray-node"

for dir in "$CONFIG_DIR" "$PERSIST_DIR" "$LOG_DIR"; do
  if [[ -d "$dir" ]]; then
    warn "Directory already exists: $dir"
  else
    mkdir -p "$dir"
    ok "Created directory: $dir"
  fi
done

chown root:xray-node "$CONFIG_DIR"
chmod 775 "$CONFIG_DIR"

chown xray-node:xray-node "$PERSIST_DIR"
chmod 770 "$PERSIST_DIR"

chown root:xray-node "$LOG_DIR"
chmod 775 "$LOG_DIR"

ok "Permissions configured for config, state and logs"

########## 6. Logrotate ##########

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

ok "Logrotate configured for xray-node logs"

########## 7. Config files ##########

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
    # uncomment for cloudflare warp outbound
    #{
    #  "tag": "warp",
    #  "protocol": "wireguard",
    #  "settings": {
    #    # paste data from generated WARP config, IPv4 or v6 depends on your server
    #    "secretKey": "#################",
    #    "DNS": ["####","####","####"],
    #    "address": ["####","########"],
    #    "peers": [
    #      {
    #        "publicKey": "#################",
    #        "endpoint": "#################",
    #      }
    #    ],
    #    "mtu": ####,
    #    "domainStrategy": "ForceIP",
    #    "kernelMode": false
    #  }
    #},
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

CONFIG_CLIENT_CONTENT=$(cat <<EOF
[]
EOF
)

create_config() {
  local file="$1"
  local name="$2"
  local content="$3"

  if [[ -f "$file" ]]; then
    read -rp "Config $name exists. Type YES to overwrite: " confirm
    if [[ "$confirm" != "YES" ]]; then
      warn "Kept existing config: $name"
      return
    fi
  fi

  echo "$content" > "$file"
  fixOwnership "$file"

  ok "Created config file: $name"
}

create_config "$CONFIG_DIR/xray_server.json" "xray_server.json" "$CONFIG_SERVER_CONTENT"
create_config "$CONFIG_DIR/xray_client.json" "xray_client.json" "$CONFIG_CLIENT_CONTENT"

########## 9. Warp config ############

read -rp "Generate Warp config via wgcf? (press ENTER = YES, type NO = skip): " WARP_CONFIRM

if [[ "${WARP_CONFIRM^^}" == "NO" ]]; then
  warn "Warp generation skipped by user"
else

  WARP_BIN_URL="https://github.com/ViRb3/wgcf/releases/download/v2.2.31/wgcf_2.2.31_linux_amd64"
  WARP_BIN="$TMP_DIR/wgcf"

  if curl -fsSL "$WARP_BIN_URL" -o "$WARP_BIN"; then
    chmod +x "$WARP_BIN"
    ok "Downloaded wgcf binary"

    echo "NOTE: auto-accepting wgcf terms (Yes)"

    set +e
    yes "Yes" | timeout 60s "$WARP_BIN" register
    REGISTER_STATUS=$?
    set -e

    if [[ $REGISTER_STATUS -eq 0 ]]; then
      ok "wgcf register completed"
    else
      warn "wgcf register failed or interrupted, continuing"
    fi

    set +e
    timeout 30s "$WARP_BIN" generate
    GENERATE_STATUS=$?
    set -e

    if [[ $GENERATE_STATUS -ne 0 ]]; then
      warn "wgcf generate failed, warp config wont be created but continuing"
    else
      if [[ -f "wgcf-profile.conf" ]]; then
        mv wgcf-profile.conf "$CONFIG_DIR/warp-generated.conf"
        fixOwnership "$CONFIG_DIR/warp-generated.conf"
        ok "Warp config generated and saved"
      else
        warn "wgcf did not produce wgcf-profile.conf"
      fi
    fi

  else
    warn "Failed to download wgcf binary, warp config wont be generated but continuing"
  fi

fi

########## 8. Endpoint port ##########

read -rp "Enter endpoint port [5001]: " ENDPOINT_PORT
ENDPOINT_PORT="${ENDPOINT_PORT:-5001}"

if ! [[ "$ENDPOINT_PORT" =~ ^[0-9]+$ ]] || (( ENDPOINT_PORT < 1 || ENDPOINT_PORT > 65535 )); then
  fail "Invalid endpoint port"
fi
ok "Endpoint port set to $ENDPOINT_PORT"

########## 9. systemd service ##########

SERVICE_FILE="/etc/systemd/system/xray-node.service"

cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Xray Node Service
After=network.target

[Service]
User=xray-node
Group=xray-node
WorkingDirectory=$BIN_DIR
ExecStart=$BIN_DIR/xray-node \\
  -a :$ENDPOINT_PORT \\
  -c $CONFIG_DIR \\
  -d $DATA_DIR \\
  -p $PERSIST_DIR
Restart=always
RestartSec=5
StandardOutput=append:$LOG_DIR/xray-node.log
StandardError=append:$LOG_DIR/xray-node.error.log

[Install]
WantedBy=multi-user.target
EOF

ok "systemd service file created"

########## 10. Enable service ##########

systemctl daemon-reload
systemctl enable xray-node >/dev/null

ok "systemd service enabled (not started)"

########## 11. Symlinks ##########
read -rp "Enter directory for symlinks (leave empty to skip): " LINK_DIR

if [[ -z "$LINK_DIR" ]]; then
  warn "Symlinks creation skipped"
else
  mkdir -p "$LINK_DIR"

  link_safe() {
    local target="$1"
    local link="$2"

    if [[ -e "$link" ]]; then
      warn "Symlink already exists: $link"
    else
      ln -s "$target" "$link"
      ok "Created symlink: $link → $target"
    fi
  }

  link_safe "$BIN_DIR" "$LINK_DIR/xray-node"
  link_safe "$CONFIG_DIR" "$LINK_DIR/config"
  link_safe "$LOG_DIR" "$LINK_DIR/logs"
  link_safe "$PERSIST_DIR" "$LINK_DIR/persist"
  link_safe "$SERVICE_FILE" "$LINK_DIR/systemd-service"

  ok "Symlinks setup completed"
fi

########## 12. Final status ##########

echo
echo "========== INSTALLATION SUMMARY =========="
for s in "${STEP_STATUS[@]}"; do
  echo "$s"
done
echo "=========================================="
echo
echo "Next steps:"
echo "  1. Edit configs in $CONFIG_DIR"
echo "  2. Start service: sudo systemctl start xray-node"
echo
