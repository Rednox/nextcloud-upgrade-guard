#!/usr/bin/env bash
set -euo pipefail

TARGET_MAJOR="${TARGET_MAJOR:-35}"
NC_ROOT="${NC_ROOT:-/var/www/nextcloud}"
POLICY_FILE="${POLICY_FILE:-/etc/nc-guard/policy.yaml}"
PHP_BIN="${PHP_BIN:-php}"

nc-guard gate --target "$TARGET_MAJOR" --policy "$POLICY_FILE" --occ-path "$NC_ROOT/occ" --php-bin "$PHP_BIN"
"$PHP_BIN" "$NC_ROOT/updater/updater.phar" --no-interaction
"$PHP_BIN" "$NC_ROOT/occ" upgrade
