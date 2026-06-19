# RELEASE-INSTALL — nc-guard operator guide

This file is included in every release zip alongside the binary and scripts.

## Contents of this zip

| File | Description |
|------|-------------|
| `nc-guard` | Static Linux binary (amd64 or arm64) |
| `preupgrade.sh` | Automation script: gate + upgrade |
| `policy.strict-security.yaml` | Example strict policy for production |
| `RELEASE-INSTALL.md` | This file |

---

## Requirements

- Linux (amd64 or arm64)
- PHP available as `php` (or set `PHP_BIN`)
- Nextcloud installed at `/var/www/nextcloud` (or set `NC_ROOT`)
- Root or sudo access to write to `/usr/local/bin` and `/etc/nc-guard`

---

## Install

### 1. Verify the checksum

Download `SHA256SUMS` from the same GitHub Release and verify:

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

### 2. Unpack the zip

```bash
unzip nc-guard_<VERSION>_linux_<ARCH>.zip -d nc-guard-release
cd nc-guard-release
```

### 3. Install the binary

```bash
sudo install -m 0755 nc-guard /usr/local/bin/nc-guard
nc-guard --version
```

### 4. Install the preupgrade script

```bash
sudo install -m 0755 preupgrade.sh /usr/local/bin/preupgrade-nc.sh
```

### 5. Place a policy file

```bash
sudo mkdir -p /etc/nc-guard
sudo install -m 0644 policy.strict-security.yaml /etc/nc-guard/policy.yaml
```

Edit `/etc/nc-guard/policy.yaml` to match your critical apps.

---

## Usage

### Inspect currently enabled apps

```bash
nc-guard inspect --occ-path /var/www/nextcloud/occ --php-bin php --format table
```

### Check compatibility before upgrading to a major version

```bash
nc-guard check --target 35 --occ-path /var/www/nextcloud/occ --php-bin php --format table
```

### Gate-only (exit code signals pass/fail for scripting)

```bash
nc-guard gate \
  --target 35 \
  --policy /etc/nc-guard/policy.yaml \
  --occ-path /var/www/nextcloud/occ \
  --php-bin php
```

Exit codes:

| Code | Meaning |
|------|---------|
| 0 | Policy pass — safe to upgrade |
| 10 | Incompatible app found, policy blocks |
| 11 | Unknown app found, policy blocks |
| 12 | Critical app at risk, policy blocks |
| 20 | Operational/tooling error |

### Full preupgrade automation

```bash
TARGET_MAJOR=35 \
NC_ROOT=/var/www/nextcloud \
POLICY_FILE=/etc/nc-guard/policy.yaml \
PHP_BIN=php \
preupgrade-nc.sh
```

When the gate passes, the script runs:
1. `php /var/www/nextcloud/updater/updater.phar --no-interaction`
2. `php /var/www/nextcloud/occ upgrade`

If the gate fails, the script exits immediately with the nc-guard exit code and the upgrade is **not** triggered.

---

## Customising the policy

Key fields in `policy.yaml`:

```yaml
block_on_incompatible: true   # block if any app reports incompatibility
block_on_unknown: true        # block if any app cannot be resolved
critical_apps:                # apps that must pass even when block_on_unknown is false
  - oidc_login
  - user_saml
  - twofactor_totp
  - twofactor_webauthn
block_on_critical_risk: true  # block if a critical app is unknown or incompatible
```

Set `block_on_unknown: false` in less strict environments but keep `block_on_critical_risk: true` for security-sensitive apps.

---

## Uninstall

```bash
sudo rm /usr/local/bin/nc-guard
sudo rm /usr/local/bin/preupgrade-nc.sh
sudo rm -r /etc/nc-guard          # only if you no longer need the policy
```
