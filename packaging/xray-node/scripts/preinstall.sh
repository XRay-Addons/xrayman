#!/bin/bash
set -e

if ! id -u xray-node &>/dev/null; then
  useradd --system --no-create-home --shell /usr/sbin/nologin xray-node
fi