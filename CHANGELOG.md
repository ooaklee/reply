# Changelog

Notable changes to GHATD are recorded here for application authors and
contributors. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
with release versions expressed using [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This record starts with the current upgrade. Earlier releases have not been
retrospectively catalogued. An unreleased entry is not a release announcement or
evidence that a feature has been deployed.

## Maintaining this file

- Update the current unreleased section in the same change as a notable feature,
  fix, security improvement, deprecation or compatibility change. Describe the
  effect on users rather than copying commit messages.
- Keep plans and unimplemented features in the project tracker. State limitations
  of partial or opt-in implementations and link to their canonical package guides.
- Mark incompatible changes with **Breaking:** and include the required migration
  action. Use public issue, PR or documentation links where useful; never include
  private project details, credentials or machine-specific paths.
- Use one `Unreleased` section until a version is selected. Then name it
  `[x.y.z] - Unreleased`. On publication, replace `Unreleased` with the actual
  `YYYY-MM-DD` release date and start a fresh unreleased section above it.
- Keep newest releases first and omit empty categories. Do not invent versions,
  dates or historical entries. Preserve published entries; clarify inaccuracies
  explicitly rather than silently changing the recorded release scope.

Before `1.0.0`, advance the minor version for features or incompatible API
changes and the patch version for fixes-only releases. Pre-release identifiers
such as `alpha`, `beta` and `rc` distinguish candidates from stable releases.
Choosing or documenting a version does not authorise creating a tag or publishing
a release.

## Changelog entry template

Use the categories that apply:

- `Added`: newly available capabilities.
- `Changed`: adjustments to existing behaviour.
- `Deprecated`: functionality scheduled for retirement.
- `Removed`: functionality no longer available.
- `Fixed`: corrected behaviour.
- `Security`: security-related changes.

The version below is a placeholder, not the next GHATD release. Replace it only
when the release version has been selected, and remove unused subsections.

```markdown
## [x.y.z] - Unreleased

### Added

- Describe the new capability and any opt-in requirements.

### Changed

- **Breaking:** Describe the compatibility impact and migration action, if any.

### Deprecated

- Identify the deprecated API and its replacement.

### Removed

- Identify the removed behaviour and migration path.

### Fixed

- Describe the corrected behaviour and affected users.

### Security

- Describe the security improvement and any required operator action.
```

---

## [2.1.0] - Unreleased

## Unreleased

### Added

- Structural `Replier.ResolveErrors` for wrapped and joined errors, using the
  premerged manifest and a bounded traversal.
- Optional request context and a privacy-safe unmapped-response observer.

### Fixed

- Update the YAML dependency used by the assertion test helpers to `v3.0.1`,
  addressing CVE-2022-28948. Reply's production package does not import YAML.
- Correctly mapped wrapped errors no longer miss their manifest entry.
- Response state is request-local, allowing shared repliers with independent,
  concurrency-safe default or custom transfer objects.
- Invalid response/prototype inputs return errors; native writer/encoder causes
  remain discoverable through `errors.Is` and `errors.As`.

### Changed

- Joined and explicit errors now have deterministic input order and unique
  identities. Incomplete trees return a generic 500 without hiding unknown
  branches. Raw fallback logging is removed; use the opt-in observer instead.
- Custom transfer factories must return independent non-nil state, not their
  prototype pointer. Nil options and attributes are ignored.

See [the integration and migration guide](UPGRADING.md) for compatibility details,
commit pinning, context propagation and custom-object ownership requirements.
