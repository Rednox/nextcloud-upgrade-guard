# nc-guard (nextcloud-upgrade-guard)

`nc-guard` is a standalone CLI that helps administrators precheck enabled Nextcloud apps before major core upgrades.

It is built as an external tool (not a patch to Nextcloud core/updater) so teams can adopt compatibility gating in their own operations workflows while upstream discussions evolve:
- nextcloud/server#34713
- nextcloud/updater#241
- nextcloud/updater#401

## Why standalone

This project intentionally avoids modifying Nextcloud core or updater internals. A standalone guard can be versioned, deployed, and automated independently.

## Release artifacts

Pre-built static Linux binaries are published on every `v*` tag via GitHub Actions.

### Supported architectures

| File | Architecture |
|------|-------------|
| `nc-guard_<VERSION>_linux_amd64.zip` | x86-64 (most servers) |
| `nc-guard_<VERSION>_linux_arm64.zip` | ARM64 (Raspberry Pi 4/5, AWS Graviton, …) |

Each zip contains:
- `nc-guard` — static binary (CGO_ENABLED=0)
- `preupgrade.sh` — automation script
- `policy.strict-security.yaml` — example strict policy
- `RELEASE-INSTALL.md` — detailed install/usage guide

A `SHA256SUMS` file is published alongside the zips for integrity verification.

## Install from release zip

```bash
# 1 — Download (replace VERSION and ARCH as needed)
VERSION=v0.1.0
ARCH=amd64
gh release download "${VERSION}" \
  --repo Rednox/nextcloud-upgrade-guard \
  --pattern "nc-guard_${VERSION}_linux_${ARCH}.zip" \
  --pattern SHA256SUMS

# 2 — Verify checksum
sha256sum -c SHA256SUMS --ignore-missing

# 3 — Unpack
unzip "nc-guard_${VERSION}_linux_${ARCH}.zip" -d nc-guard-release
cd nc-guard-release

# 4 — Install binary and script
sudo install -m 0755 nc-guard /usr/local/bin/nc-guard
sudo install -m 0755 preupgrade.sh /usr/local/bin/preupgrade-nc.sh

# 5 — Place policy file
sudo mkdir -p /etc/nc-guard
sudo install -m 0644 policy.strict-security.yaml /etc/nc-guard/policy.yaml

# 6 — Gate-check before upgrade
nc-guard gate \
  --target 35 \
  --policy /etc/nc-guard/policy.yaml \
  --occ-path /var/www/nextcloud/occ \
  --php-bin php

# 7 — Full preupgrade flow (gate → updater.phar → occ upgrade)
TARGET_MAJOR=35 POLICY_FILE=/etc/nc-guard/policy.yaml preupgrade-nc.sh
```

See [`docs/RELEASE-INSTALL.md`](docs/RELEASE-INSTALL.md) for full details and exit-code reference.

## Install / build from source

Requirements: Go 1.22+

```bash
go build -o ./bin/nc-guard ./cmd/nc-guard
```

## Commands

### Inspect enabled/disabled apps

```bash
./bin/nc-guard inspect --occ-path /var/www/nextcloud/occ --php-bin php --format table
```

### Check compatibility against a target major

```bash
./bin/nc-guard check --target 35 --occ-path /var/www/nextcloud/occ --php-bin php --format table
```

### Gate upgrades by policy

```bash
./bin/nc-guard gate --target 35 --policy examples/policy.strict-security.yaml
```

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Policy pass |
| 10 | Incompatible app found and policy blocks |
| 11 | Unknown app found and policy blocks |
| 12 | Critical app risk found and policy blocks |
| 20 | Operational/tooling error |

## Automation example

Use `scripts/preupgrade.sh` in cron/systemd/manual maintenance pipelines:

```bash
TARGET_MAJOR=35 POLICY_FILE=/etc/nc-guard/policy.yaml scripts/preupgrade.sh
```

When gate passes, the script runs:
- `php /var/www/nextcloud/updater/updater.phar --no-interaction`
- `php /var/www/nextcloud/occ upgrade`

## v0.1 limitations

- local resolver only (no remote App Store compatibility lookup yet)
- unresolved app metadata is marked `unknown`
- no network dependency by design
