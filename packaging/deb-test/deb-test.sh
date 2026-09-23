#!/bin/sh
set -eu

ARCH="${ARCH:-arm64}"

deb_test() {
    SERVICE=$1
    NFPM=$2
    DUMMY_FILES=$3
    WATCH_DIRS=$4
    OPERATIONS=$5

    TMP_DIR="$(mktemp -d)"
    PACKAGE="$TMP_DIR/package.deb"

    cleanup() {
        rm -rf "$TMP_DIR"
    }
    trap cleanup EXIT

    echo "==> Preparing temporary package tree"

    # Copy nfpm files.
    mkdir -p "$TMP_DIR/packaging/$SERVICE"
    cp -R "packaging/$SERVICE/deb" "$TMP_DIR/packaging/$SERVICE/"

    # Create dummy source files required by nfpm.
    for file in $DUMMY_FILES; do
        path="$TMP_DIR/$file"

        mkdir -p "$(dirname "$path")"
        touch "$path"
        chmod 755 "$path"

        echo "    dummy: $file"
    done

    echo
    echo "==> Checking nfpm files"
    echo "==> Building package"

    (
        cd "$TMP_DIR"

        ARCH="$ARCH" nfpm package \
            --config "$NFPM" \
            --packager deb \
            --target "$PACKAGE"
    )

    echo
    echo "==> Package: $PACKAGE"
    echo "==> Testing package in Debian"

    docker run --rm -i \
        -v "$PACKAGE:/build/package.deb:ro" \
        debian:12 \
        sh -s -- "$OPERATIONS" "$WATCH_DIRS" <<'CONTAINER'

set -u

OPERATIONS=$1
WATCH_DIRS=$2

export DEBIAN_FRONTEND=noninteractive

echo "==> Installing test dependencies"

apt-get update >/dev/null
apt-get install -y passwd >/dev/null

# systemctl stub. The test container does not run systemd.
cat > /usr/local/bin/systemctl <<'EOF'
#!/bin/sh
echo "[fake systemctl] $*"
exit 0
EOF
chmod +x /usr/local/bin/systemctl

snapshot() {
    for dir in $WATCH_DIRS; do
        if [ -e "$dir" ]; then
            find "$dir" \
                -printf '%p\t%M\t%u\t%g\n' \
                2>/dev/null
        fi
    done | sort
}

snapshot_users() {
    getent passwd | sort
}

snapshot_groups() {
    getent group | sort
}

show_file_changes() {
    before=$1
    after=$2

    awk -F '\t' '
        NR == FNR {
            before[$1] = $0
            next
        }

        {
            after[$1] = $0
        }

        END {
            for (path in after) {
                if (!(path in before)) {
                    split(after[path], f, "\t")
                    printf "+ %-60s %s %s:%s\n", f[1], f[2], f[3], f[4]
                }
            }

            for (path in before) {
                if (!(path in after)) {
                    split(before[path], f, "\t")
                    printf "- %-60s %s %s:%s\n", f[1], f[2], f[3], f[4]
                }
            }
        }
    ' "$before" "$after" |
    sort
}

show_user_changes() {
    before=$1
    after=$2

    echo "Added:"
    comm -13 "$before" "$after" || true

    echo "Removed:"
    comm -23 "$before" "$after" || true
}

show_group_changes() {
    before=$1
    after=$2

    echo "Added:"
    comm -13 "$before" "$after" || true

    echo "Removed:"
    comm -23 "$before" "$after" || true
}

# Initial state is captured once and never changed.
initial=/tmp/files.initial
current=/tmp/files.current

users_initial=/tmp/users.initial
users_current=/tmp/users.current

groups_initial=/tmp/groups.initial
groups_current=/tmp/groups.current

snapshot > "$initial"
snapshot_users > "$users_initial"
snapshot_groups > "$groups_initial"

i=0

printf '%s\n' "$OPERATIONS" |
while IFS= read -r operation; do
    [ -z "$operation" ] && continue

    i=$((i + 1))

    echo
    echo "============================================================"
    echo "==> Operation $i"
    echo "    $operation"
    echo "============================================================"

    set +e
    sh -c "$operation"
    rc=$?
    set -e

    echo
    echo "--- exit code: $rc ---"

    snapshot > "$current"
    snapshot_users > "$users_current"
    snapshot_groups > "$groups_current"

    echo
    echo "--- filesystem changes from initial state ---"

    changes=$(show_file_changes "$initial" "$current")

    if [ -n "$changes" ]; then
        printf '%s\n' "$changes"
    else
        echo "No changes."
    fi

    echo
    echo "--- users changes from initial state ---"
    show_user_changes "$users_initial" "$users_current"

    echo
    echo "--- groups changes from initial state ---"
    show_group_changes "$groups_initial" "$groups_current"
done
CONTAINER
}