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
- Support for feature, release, hotfix, and support branches
- Configurable merge strategies (merge, rebase, squash)
- Three-layer configuration precedence (defaults → git config → CLI flags)
- State machine for finish command with conflict recovery
- Automatic child branch updates with `autoUpdate` configuration
- Import configuration from git-flow-avh

[Unreleased]: https://github.com/gittower/git-flow-next/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/gittower/git-flow-next/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/gittower/git-flow-next/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/gittower/git-flow-next/releases/tag/v0.1.0
