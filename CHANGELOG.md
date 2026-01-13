# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-01-11

### Added

- Add `track` command for topic branches to track existing remote branches
- Add `--fetch` option for topic branch finish command

### Fixed

- Create base branches in dependency order during init
- Fix docs link in README.md

### Changed

- Optimize binary size (~50% reduction)

## [0.1.1] - 2025-09-24

### Fixed

- Minor bug fixes and improvements

## [0.1.0] - 2025-09-16

### Added

- Initial release of git-flow-next
- Support for feature, release, hotfix, and support branch workflows
- Fully customizable base and topic branches with configurable prefixes and relationships
- Configurable merge strategies: merge, rebase, or squash when finishing branches
- Flexible configuration via git config or command-line flags
- Conflict recovery: resolve conflicts and continue where you left off
- Automatic updates to child branches (e.g., develop syncs from main)
- Compatibility with existing git-flow-avh repositories

[Unreleased]: https://github.com/gittower/git-flow-next/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/gittower/git-flow-next/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/gittower/git-flow-next/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/gittower/git-flow-next/releases/tag/v0.1.0
