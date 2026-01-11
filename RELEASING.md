# Releasing

This document describes the release process for git-flow-next.

## Version Locations

Update version in **both** files (must match):

- `version/version.go` — Core version constant
- `cmd/version.go` — Version variable for CLI

## Semantic Versioning

We follow [Semantic Versioning](https://semver.org/):

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| Breaking changes | Major | 1.0.0 → 2.0.0 |
| New features (`feat:`) | Minor | 1.0.0 → 1.1.0 |
| Bug fixes (`fix:`) | Patch | 1.0.0 → 1.0.1 |

## Changelog Convention

We follow [Keep a Changelog](https://keepachangelog.com/). Only **user-facing changes** go in CHANGELOG.md:

| Commit Type | Changelog Section | Include? |
|-------------|-------------------|----------|
| `feat:` | Added | ✅ Yes |
| `fix:` | Fixed | ✅ Yes |
| `BREAKING CHANGE` | Changed | ✅ Yes |
| `perf:` | Changed | ✅ Yes |
| `build:` | Changed | ⚠️ If user-facing |
| `docs:` | — | ❌ No |
| `refactor:` | — | ❌ No |
| `test:` | — | ❌ No |
| `chore:` | — | ❌ No |
| `ci:` | — | ❌ No |
| `style:` | — | ❌ No |

## Release Process

### 1. Determine Version Bump

Review commits since last release:

```bash
git log $(git describe --tags --abbrev=0)..HEAD --oneline
```

- Any `feat:` → **minor** bump
- Only `fix:` → **patch** bump
- Any `BREAKING CHANGE` or `!` suffix → **major** bump

### 2. Update Changelog

Edit `CHANGELOG.md`:

1. Move items from `[Unreleased]` to new version section
2. Add release date
3. Categorize changes: Added, Changed, Fixed, Removed, Security

### 3. Update Version

Update version in both files:

```bash
# Edit version/version.go
# Edit cmd/version.go
```

Verify:

```bash
go build -o /tmp/git-flow main.go && /tmp/git-flow version
```

### 4. Commit Version Bump

```bash
git add CHANGELOG.md version/version.go cmd/version.go
git commit -m "chore: Bump version to X.Y.Z"
```

### 5. Push and Tag

```bash
git push origin main
git tag vX.Y.Z
git push origin vX.Y.Z
```

### 6. GitHub Actions

The `.github/workflows/release.yml` workflow automatically:

1. Builds binaries for all platforms (darwin, linux, windows)
2. Creates GitHub release with artifacts
3. Generates checksums
4. Marks as prerelease if tag contains `-alpha`, `-beta`, or `-rc`

### 7. Update Homebrew Tap

After the GitHub release is published, update the Homebrew formula:

```bash
cd /path/to/homebrew-tap
./update_formula.rb
git add Formula/
git commit -m "Update git-flow-next to vX.Y.Z"
git push
```

Repository: https://github.com/gittower/homebrew-tap

### 8. Update Website Documentation

If commands or options changed, update the documentation website:

1. Review changes to `docs/*.md` files
2. Update corresponding pages in the website repo
3. Rebuild and deploy

Repository: https://github.com/gittower/git-flow-next-website

## Preview Releases

For preview releases, use suffixes:

- Alpha: `v1.0.0-alpha.1`
- Beta: `v1.0.0-beta.1`
- Release candidate: `v1.0.0-rc.1`

These are automatically marked as prereleases on GitHub.

## Checklist

- [ ] Reviewed commits since last release
- [ ] Determined correct version bump (major/minor/patch)
- [ ] Updated CHANGELOG.md with user-facing changes
- [ ] Updated version in `version/version.go`
- [ ] Updated version in `cmd/version.go`
- [ ] Verified build: `go build && ./git-flow version`
- [ ] Committed: `chore: Bump version to X.Y.Z`
- [ ] Pushed to main
- [ ] Created and pushed tag
- [ ] Verified GitHub release was created
- [ ] Updated Homebrew tap (`./update_formula.rb`)
- [ ] Updated website documentation (if commands changed)
