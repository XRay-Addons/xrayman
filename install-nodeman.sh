#!/usr/bin/env bash
set -euo pipefail

APP=xray-nodeman
BIN_DIR=/opt
BIN_PATH=$BIN_DIR/$APP
ENV_DIR=/etc/$APP
PACKAGE_BIN_PATH="$APP/$APP"
EXEC_CMD=$(cat <<EOF
$BIN_PATH \
--log-lvl info
EOF
)

# default env content shipped with this installer
DEFAULT_ENV_CONTENT=$(cat <<'EOF'
ENDPOINT=127.0.0.1:5002
DBCONN=postgresql://user:password@host:5432/dbname
JWT_SECRET=Secret***
ADMIN_PASSWORD=****
EOF
)

ENV_FILE=$ENV_DIR/$APP.env
SERVICE_FILE=/etc/systemd/system/$APP.service


# detect OS and map it to the naming used in release assets
case "$(uname -s)" in
  Linux)  OS=linux ;;
  Darwin) OS=darwin ;;
  *)
    echo "ERROR: unsupported OS $(uname -s)" >&2
    exit 1
    ;;
esac
 
# detect architecture and map it to the naming used in release assets
case "$(uname -m)" in
  x86_64)  ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
  arm64)   ARCH=arm64 ;;
  armv7l)  ARCH=armv7 ;;
  *)
    echo "ERROR: unsupported architecture $(uname -m)" >&2
    exit 1
    ;;
esac
 
DOWNLOAD_URL="https://github.com/XRay-Addons/xrayman/releases/latest/download/xrayman-nodeman-$OS-$ARCH.tar.gz"

log() {
  echo "==> $*"
}

# --- temporary working directory, cleaned up automatically on exit ---
WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT

# --- 1. create user/group if they don't exist yet ---
if id -u "$APP" &>/dev/null; then
  log "User '$APP' already exists, skipping"
else
  useradd --system --no-create-home --shell /usr/sbin/nologin "$APP"
  log "User '$APP' created"
fi

# --- 2. download tar.gz archive into temp folder ---
log "Downloading $DOWNLOAD_URL"
curl -sL "$DOWNLOAD_URL" -o "$WORKDIR/$APP.tar.gz"
log "Download complete"

# --- 3. extract into temp folder (not into BIN_DIR directly) ---
log "Extracting archive"
mkdir -p "$WORKDIR/extracted"
tar -xzf "$WORKDIR/$APP.tar.gz" -C "$WORKDIR/extracted"

NEW_BIN="$WORKDIR/extracted/$PACKAGE_BIN_PATH"

if [ ! -f "$NEW_BIN" ]; then
  echo "ERROR: binary not found at $NEW_BIN after extraction" >&2
  echo "Archive contents:" >&2
  find "$WORKDIR/extracted" -type f >&2
  exit 1
fi

chmod +x "$NEW_BIN"

# --- 4. get version and help output from the extracted (not yet installed) binary ---
# doing this in one place, before install, so we don't need to call the
# binary again later from its final location
VERSION_OUTPUT=$("$NEW_BIN" --version 2>/dev/null || echo "unknown")
HELP_OUTPUT=$("$NEW_BIN" --help 2>&1 || true)
log "Downloaded binary version: $VERSION_OUTPUT"

# --- 5. stop the service before touching the binary ---
# we don't check whether it exists or is running - stopping a
# nonexistent/already-stopped service is harmless with '|| true'
systemctl stop "$APP" 2>/dev/null || true
log "Service '$APP' stopped (if it was running)"

# --- 6. install new binary, always replacing whatever was there before ---
install -m 755 -o root -g root "$NEW_BIN" "$BIN_PATH"
log "Binary installed to $BIN_PATH (version $VERSION_OUTPUT)"

# --- 7. build the env file ---
mkdir -p "$ENV_DIR"

ENV_EXISTED=false
OLD_ENV_CONTENT=""
if [ -f "$ENV_FILE" ]; then
  ENV_EXISTED=true
  OLD_ENV_CONTENT=$(cat "$ENV_FILE")
  log "Existing env file found at $ENV_FILE, keeping its content untouched"
else
  log "No existing env file found, creating a new one"
fi

{
  # keep old content exactly as-is, uncommented, at the top
  if [ "$ENV_EXISTED" = true ]; then
    echo "$OLD_ENV_CONTENT"
    echo ""
  fi

  echo "# ============================================================"
  echo "# ============================================================"
  echo "# ============================================================"
  echo "# Updated by install.sh"
  echo "# Installed at: $(date -u +'%Y-%m-%d %H:%M:%S UTC')"
  echo "#"
  echo "# --- Application version output (--version) ---"
  echo "$VERSION_OUTPUT" | sed 's/^/# /'
  echo "#"
  echo "# ============================================================"
  echo "#"
  echo "# --- Application help output (--help) ---"
  echo "$HELP_OUTPUT" | sed 's/^/# /'
  echo "#"
  echo "# ============================================================"
  echo "#"
  echo "# --- Systemd service command ---"
  echo "$EXEC_CMD" | sed 's/^/# /'
  echo "#"
  echo "# --- Env options example ---"
  echo "$DEFAULT_ENV_CONTENT" | sed 's/^/# /'
} > "$ENV_FILE"

chown root:"$APP" "$ENV_FILE"
chmod 640 "$ENV_FILE"
log "Env file written to $ENV_FILE"

# --- 8. write systemd unit file ---
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=$APP service
After=network.target

[Service]
Type=simple
User=$APP
Group=$APP
EnvironmentFile=$ENV_FILE
ExecStart=$EXEC_CMD
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
log "Systemd unit file written to $SERVICE_FILE"

systemctl daemon-reload
log "systemd daemon reloaded"

systemctl enable "$APP"
log "Service '$APP' enabled (autostart on boot)"

# --- 9. ask for a folder to hold symlinks to the key files ---
read -rp "Enter folder path to store shortcut symlinks (optional): " LINKS_DIR

if [ -n "$LINKS_DIR" ]; then
  mkdir -p "$LINKS_DIR"

  ln -sf "$BIN_PATH" "$LINKS_DIR/$APP"
  ln -sf "$ENV_FILE" "$LINKS_DIR/$APP.env"
  ln -sf "$SERVICE_FILE" "$LINKS_DIR/$APP.service"

  log "Symlinks created in $LINKS_DIR:"
  log "  $LINKS_DIR/$APP         -> $BIN_PATH"
  log "  $LINKS_DIR/$APP.env     -> $ENV_FILE"
  log "  $LINKS_DIR/$APP.service -> $SERVICE_FILE"
fi

echo ""
echo "Done: $APP installed:"
echo "$VERSION_OUTPUT"
echo "Edit $ENV_FILE config file and run service"
echo "	via 'sudo systemctl restart $APP'"
echo "View service logs"
echo "	via 'journalctl -u $APP -f -n 50'"