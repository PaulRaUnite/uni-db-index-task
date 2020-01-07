# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.2

### Fixed

Encoding from filter of type `[]string` led to panic.

## 1.0.1

### Fixed

- Cast error when decoding filter to a `[]string`.

## 1.0.0

### Added

* Support for `search` and `sort` query params.
* Support for nested structures
* Support for type aliases
* Default parameters
* Decode now returns typed errors with explanations of what when wrong (can be rendered to client).

### Fixed
* remove uint only restriction for page params
* pointer destinations were not really supported

## 0.1.0
### Added
* proof-of-concept implementation