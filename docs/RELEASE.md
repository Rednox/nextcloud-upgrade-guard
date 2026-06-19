# Release guide — nc-guard

## Cutting a release

### Prerequisites

- You have push access to the repository.
- All changes you want in the release are merged to the default branch.
- Your local clone is up to date: `git pull`.

### Create and push a tag

```bash
# Test release (dry-run, real workflow triggers)
git tag v0.1.0-rc1
git push origin v0.1.0-rc1

# Production release
git tag v0.1.0
git push origin v0.1.0
```

Tags must match the pattern `v*` to trigger the release workflow.  
**Never rewrite a published tag** — cut a new patch tag instead (see Rollback below).

---

## Workflow behaviour

File: `.github/workflows/release.yml`

| Step | Details |
|------|---------|
| Trigger | Push of any tag matching `v*` |
| Go version | 1.22 |
| Test gate | `go test ./...` and `go vet ./...` — workflow fails if either fails |
| Build matrix | `GOOS=linux`, `GOARCH`: `amd64`, `arm64`, `CGO_ENABLED=0`, `-trimpath -ldflags="-s -w"` |
| Zip contents | `nc-guard`, `preupgrade.sh`, `policy.strict-security.yaml`, `RELEASE-INSTALL.md` |
| Artifact names | `nc-guard_<VERSION>_linux_amd64.zip`, `nc-guard_<VERSION>_linux_arm64.zip` |
| Checksum file | `SHA256SUMS` (SHA-256 of all zips) |
| Release | GitHub Release created with auto-generated notes + install header; all zips and SHA256SUMS attached |

The `publish` job runs only after both arch builds succeed.

---

## Verifying artifacts

After the release is published:

```bash
# Download artifacts (example with gh CLI)
gh release download v0.1.0 --repo Rednox/nextcloud-upgrade-guard

# Verify checksums
sha256sum -c SHA256SUMS

# Inspect zip contents
unzip -l nc-guard_v0.1.0_linux_amd64.zip
```

Expected zip contents:

```
nc-guard
preupgrade.sh
policy.strict-security.yaml
RELEASE-INSTALL.md
```

---

## Rollback / fix strategy

**Never delete or move a published tag.** Consumers may have already downloaded the artifact.

If a release contains a bug:

1. Fix the bug on the default branch.
2. Cut a new patch tag:
   ```bash
   git tag v0.1.1
   git push origin v0.1.1
   ```
3. The new release supersedes the previous one. Mark the bad release as a pre-release or add a warning to its description via the GitHub UI if needed.

---

## Required GitHub settings

The release workflow uses `softprops/action-gh-release@v2`, which creates the release using the built-in `GITHUB_TOKEN`.  
No additional secrets or repository settings are required beyond the default Actions permissions.

Confirm once:

- **Settings → Actions → General → Workflow permissions** is set to  
  _"Read repository contents and packages permissions"_ (default) — the job overrides to `contents: write` as needed.
