copy() {
    mkdir -p "$(dirname "$2")"
    cp -r "$1" "$2"
    if [[ -n "$3" ]]; then
        chmod "$3" "$2"
    fi
}

replace() {
    local src="$1"
    local dst="$2"
    shift 2

    local args=()
    while (($#)); do
        args+=("-e" "s/$1/$2/")
        shift 2
    done

    mkdir -p "$(dirname "$dst")"
    sed "${args[@]}" "$src" > "$dst"
}

createServiceUser() {
    local SERVICE_USER="$1"
    local SERVICE_GROUP="$2"

    if ! getent group "$SERVICE_GROUP" >/dev/null; then
        echo "Create system group '$SERVICE_GROUP'"
        groupadd --system "$SERVICE_GROUP"
    else
        echo "System group '$SERVICE_GROUP' already exists, skipping"
    fi

    if ! getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
        echo "Create system user '$SERVICE_USER'"
        useradd --system --gid "$SERVICE_GROUP" --no-create-home --shell /usr/sbin/nologin "$SERVICE_USER"
    else
        echo "System user '$SERVICE_USER' already exists, skipping"
    fi
}

deleteServiceUser() {
    local SERVICE_USER="$1"
    local SERVICE_GROUP="$2"

    if getent passwd "$SERVICE_USER" >/dev/null 2>&1; then
        echo "Removing system user '$SERVICE_USER'"
        userdel "$SERVICE_USER"
    else
        echo "System user '$SERVICE_USER' does not exist, skipping"
    fi

    if getent group "$SERVICE_GROUP" >/dev/null 2>&1; then
        echo "Removing system group '$SERVICE_GROUP'"
        groupdel "$SERVICE_GROUP"
    else
        echo "System group '$SERVICE_GROUP' does not exist, skipping"
    fi
}

installSystemdService() {
    local SERVICE="$1"

    DPKG_GARBAGE_FILE="/etc/$SERVICE/$SERVICE.env.dpkg-dist"
    if [ -f "$DPKG_GARBAGE_FILE" ]; then
        echo "Removing dpkg garbage '$DPKG_GARBAGE_FILE'"
        rm -f "$DPKG_GARBAGE_FILE"
    fi

    echo "Reload systemd units."
    systemctl daemon-reload

    echo "Enable $SERVICE service."
    systemctl enable $SERVICE


    printf '\n'
    printf '╔════════════════════════════════════════════════════════════════╗\n'
    printf '║     XRayMan service successfully installed: %-16s   ║\n' "$SERVICE"
    printf '╠════════════════════════════════════════════════════════════════╣\n'
    printf '║ [!] Setup app configuration [!]                                ║\n'
    printf '║   [!] User manual available via sudo apt show %-16s ║\n' "$SERVICE"
    printf '║                                                                ║\n'
    printf '║ Then start/restart the service:                                ║\n'
    printf '║   sudo systemctl restart %-16s                      ║\n' "$SERVICE"
    printf '║                                                                ║\n'
    printf '╚════════════════════════════════════════════════════════════════╝\n'
    printf '\n'
}

removeSystemdService() {
    local SERVICE="$1"

    if systemctl is-active --quiet $SERVICE.service; then
        echo "Stopping $SERVICE service"
        systemctl stop $SERVICE.service
    else
        echo "$SERVICE service is already stopped, skipping"
    fi

    if systemctl is-enabled --quiet "$SERVICE.service"; then
        echo "Disabling $SERVICE service"
        systemctl disable $SERVICE.service
    else
        echo "$SERVICE service is already disabled, skipping"
    fi
}