test_devlxd() {
  if [ -n "$TRAVIS_PULL_REQUEST" ]; then
    echo "SKIPPING"
    return
  fi

  if ! lxc image alias list | grep -q "^| testimage\s*|.*$"; then
    if [ -e "$LXD_TEST_IMAGE" ]; then
        lxc image import $LXD_TEST_IMAGE --alias testimage
    else
        ../scripts/lxd-images import busybox --alias testimage
    fi
  fi

  go build -tags netgo -a -installsuffix devlxd devlxd-client.go
  lxc launch testimage devlxd

  lxc file push devlxd-client devlxd/bin/
  lxc exec devlxd chmod +x /bin/devlxd-client

  lxc config set devlxd user.foo bar
  lxc exec devlxd devlxd-client user.foo | grep bar

  lxc stop devlxd --force
}
