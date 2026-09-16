#!/bin/sh
set -e
 
CONFIG_DIR="/etc/xray-node"
ENV_FILE="${CONFIG_DIR}/xray-node.env"
 
# Владелец/права на директорию и persistent-каталог.
# Владелец самого env-файла уже выставлен через file_info в nfpm.yaml —
# здесь его трогать не нужно (и не надо, чтобы не сбить noreplace-логику
# при апгрейде, если пользователь поменял права руками).
if [ -d "$CONFIG_DIR" ]; then
    chown root:root "$CONFIG_DIR"
    chmod 755 "$CONFIG_DIR"
fi
 
if [ -d /var/lib/xray-node/persistent ]; then
    chown xray-node:xray-node /var/lib/xray-node/persistent
    chmod 750 /var/lib/xray-node/persistent
fi
 
echo "Reload systemd units."
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
    # Включаем автозапуск, но не стартуем сразу — пользователь может
    # захотеть сначала отредактировать конфиг. Уберите строку ниже,
    # если хотите стартовать сервис сразу же после установки.
    systemctl enable xray-node.service || true
fi

echo
echo "╔═════════════════════════════════════════════════════════════╗"
echo "║             XRayMan Node installed successfully!            ║"
echo "╠═════════════════════════════════════════════════════════════╣"
echo "║ [!] Edit .env file before starting the service.             ║"
echo "║   sudo nano ${ENV_FILE}                                      "
echo "║ [!] Setup app configuration [!]                             ║"
echo "║   [!] User manual available via sudo apt show xray-node     ║"
echo "║   [!] User manual available via sudo apt show xray-node     ║"
echo "║   [!] User manual available via sudo apt show xray-node     ║"
echo "║                                                             ║"
echo "║ Then start (or restart) the service:                        ║"
echo "║   sudo systemctl enable --now xray-node                     ║"
echo "║   or: sudo systemctl restart xray-node.                     ║"
echo "║                                                             ║"
echo "╚═════════════════════════════════════════════════════════════╝"

exit 0
