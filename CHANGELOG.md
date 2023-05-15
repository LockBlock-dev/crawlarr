# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2023-05-15

### Added

-   `SAME_ORIGIN` match type, see the [README](/README.md#matching-types)

### Fixed

-   Concurrent write/read to maps
-   Match type now defaults to `SAME_BASE`

### Changed

-   Delay is now applied as a rate limit

## [1.0.1] - 2023-05-09

### Changed

-   Project now follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## [1.0.0] - 2023-05-08

### Added

-   Initial release
