# Test .deb package and systemd service installation

1. Install latest (main) version of `nfpm` from go repo (apt and other versions is not good enough, nfpm manual and config examples are valid only for main version) and docker.

2. Open xray-node-test.sh or xray-nodeman-test.sh and edit list of commands or list of files available on package creation, including
   `dpkg --install, --remove, --purge`, editing config and env files

3. Run from root test process and view list of files, users etc after each command:

```sh
./packaging/deb-test/xray-node-test.sh
or
./packaging/deb-test/xray-nodeman-test.sh
```

4. Analyze output and check everything created and deleted as desired.
