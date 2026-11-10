# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

The releases below describe the original upstream project.

## [0.3.1] - 2025-01-06

### Fixed
- Fixed documentation naming: renamed `os.md` to `os_list.md` to match data source name `aeza_os_list`
- Updated provider documentation links

## [0.3.0] - 2025-01-06

### Added
- Complete English documentation for Terraform Registry
  - Provider documentation (`docs/index.md`)
  - Data sources documentation (products, services, service_types, os, recipes, service_groups)
  - Resources documentation (service, service_actions, service_prolong)
- Authentication guide with API key configuration examples
- Usage examples for all data sources and resources

## [0.2.0] - 2025-01-06

### Added
- GPG signing for all releases
- Signed SHA256SUMS checksums for security verification

### Changed
- Enabled GPG signing in release workflow
- All future releases will be cryptographically signed

## [0.1.7] - 2025-01-06

### Changed
- Moved main.go from cmd/terraform-provider-aeza/ to root directory
- Updated Makefile to build from root
- Simplified .goreleaser.yml (removed dir/main config)

## [0.1.6] - 2025-01-06 [FAILED]

### Fixed
- Use both `dir: .` and `main: ./cmd/terraform-provider-aeza` in GoReleaser config

## [0.1.5] - 2025-01-06 [FAILED]

### Fixed
- Changed `main:` to `dir:` in .goreleaser.yml (GoReleaser v2 syntax)

## [0.1.4] - 2025-01-06 [FAILED]

### Fixed
- Fixed main path in .goreleaser.yml (removed leading `./`)

## [0.1.3] - 2025-01-06 [FAILED]

### Fixed
- Simplified GoReleaser configuration to minimal working version
- Updated GoReleaser action to v6 with explicit v2 distribution
- Removed complex ignore rules and platform combinations

## [0.1.2] - 2025-01-06 [FAILED]

### Changed
- Temporarily disabled GPG signing until secrets are configured in GitHub
- Commented out GPG import step in release workflow

### Fixed
- Fixed GoReleaser release workflow to work without GPG signing

## [0.1.1] - 2025-01-06 [FAILED]

### Added
- GoReleaser configuration for multi-platform builds
- GitHub Actions workflow for automated releases
- GitHub Actions workflow for CI/CD testing
- GPG signing for releases
- Automated changelog generation
- GPG setup documentation

### Changed
- Release archives now include LICENSE and README.md files
- Improved .gitignore to exclude dist/ directory

### Fixed
- Deprecated GoReleaser configuration options

## [0.1.0] - 2025-01-05

### Added
- Initial release of Aeza Terraform Provider
- Support for VPS service creation and management
- Data sources: products, services, OS, groups, types, recipes
- Resources: service, service_actions, service_prolong
- API client with Legacy and V2 API support
- Examples for basic usage, data sources, and resources
- MPL-2.0 license
- Comprehensive README with usage examples

[Unreleased]: https://github.com/mrB4el/terraform-provider-aeza/compare/v0.1.1...HEAD
[0.3.1]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.7...v0.2.0
[0.1.7]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/scinfra-pro/terraform-provider-aeza/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/scinfra-pro/terraform-provider-aeza/releases/tag/v0.1.0

