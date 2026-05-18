# services_utils
Utilities for my REST services

## Packages

- `api_error`: common API error type and HTTP status constructors.
- `date`: RFC3339 date/time helpers for API consistency.
- `enums`: simple indexed string enum helpers.
- `logger`: JSON logging wrapper with in-memory log list support and optional file rotation.

## Removed packages

- `misc`: removed in favor of native Go helpers such as `slices.Contains` and `slices.ContainsFunc`.
