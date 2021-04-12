# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.1] - 2021-04-06

### Changed

- Updated Dockerfiles to pull base images from Artifactory instead of DTR.

## [1.4.0] - 2021-02-02

### Changed

- Updated to MIT License in all files

## [1.3.1] - 2021-01-22

### Changed

- Added User-Agent header to all outbound HTTP requests.

## [1.3.0] - 2021-01-14

### Changed

- Updated license file.


## [1.2.2] - 2020-10-29

### Security

- CASMHMS-4148 - Update go module vendor code for security fix.

## [1.2.1] - 2020-10-20

### Security

- CASMHMS-4105 - Updated base Golang Alpine image to resolve libcrypto vulnerability.
- CASMHMS-4090 - Vendor library code.

## [1.2.0] - 2020-08-20

### Changed

- Now uses a TLS cert-aware HTTP client pair for all NWP transactions.

## [1.1.1] - 2020-04-28

### Changed

- CASMHMS-2969 - Updated hms-bmc-networkprotocol to use trusted baseOS.

## [1.1.0] - 2020-03-03

### Added

- Added a new init function: InitWithNWP() which supports SSH keys and boot order.  Existing Init() function remains in place.

## [1.0.2] - 2020-02-25

### Changed

- CASMHMS-3012 - no longer set NTP server directly

## [1.0.1] - 2019-12-20

### Changed

- Improved NW protocol parameter parsing
- Increased Redfish timeout from 10 to 17 seconds.
- Fixed HTTP client leak

## [1.0.0] - 2019-09-17

### Added

- This is the initial release of the `hms-bmc-networkprotocol` repo. It contains common code used by REDS and MEDS but could be used by any service that wants to set Mountain controller network protocol stuff (NTP, syslog to start with).

### Changed

### Deprecated

### Removed

### Fixed

### Security
