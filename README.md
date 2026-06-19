# nc-guard (nextcloud-upgrade-guard)

`nc-guard` is a standalone CLI that helps administrators precheck enabled Nextcloud apps before major core upgrades.

It is built as an external tool (not a patch to Nextcloud core/updater) so teams can adopt compatibility gating in their own operations workflows while upstream discussions evolve:
- nextcloud/server#34713
- nextcloud/updater#241
- nextcloud/updater#401

## Why standalone

This project intentionally avoids modifying Nextcloud core or updater internals. A standalone guard can be versioned, deployed, and automated independently.

## Install / build

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
