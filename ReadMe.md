## [NO] One-command Install xray-node

```
curl -fsSL https://raw.githubusercontent.com/XRay-Addons/xrayman/main/install-node.sh -o install-node.sh && sudo bash install-node.sh
```

## [NO] One-command Install xray-nodeman

```
curl -fsSL https://raw.githubusercontent.com/XRay-Addons/xrayman/main/install-nodeman.sh -o install-nodeman.sh && sudo bash install-nodeman.sh
```

## Install xray node manager

### .deb (Ubuntu, Debian, etc...)

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

### docker

```

```
