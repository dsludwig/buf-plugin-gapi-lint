# buf-plugin-gapi-lint

A [buf](https://buf.build) lint plugin that wraps the [Google API Linter](https://github.com/googleapis/api-linter), exposing all of its rules through the buf plugin interface. This gives you Google's API design review guidelines enforced directly in your `buf lint` workflow, with per-rule configurability via `buf.yaml`.

## Installation

### With mise

```bash
mise use github:dsludwig/buf-plugin-gapi-lint
```

### From GitHub releases

Download the binary for your platform from the [releases page](https://github.com/dsludwig/buf-plugin-gapi-lint/releases), extract it, and place it on your `PATH`.

For example on macOS (Apple Silicon):

```bash
curl -sSL https://github.com/dsludwig/buf-plugin-gapi-lint/releases/latest/download/buf-plugin-gapi-lint_Darwin_arm64.tar.gz | tar xz
sudo mv buf-plugin-gapi-lint /usr/local/bin/
```

### From source

```bash
go install github.com/dsludwig/buf-plugin-gapi-lint@latest
```

## Configuration

Add the plugin to your `buf.yaml` and reference its rules in the `lint` section:

```yaml
version: v2
lint:
  use:
    - STANDARD
    # Enable specific api-linter rules or categories
    - AIP_0203
    - CORE_0140_LOWER_SNAKE
  except:
    # Disable specific rules
    - CORE_0191_JAVA_PACKAGE
plugins:
  - plugin: buf-plugin-gapi-lint
```

### Listing available rules

To see all rules and categories provided by the plugin:

```bash
buf config ls-lint-rules
```

### Rule naming

Each rule from the api-linter is exposed as a buf lint rule. The rule IDs are derived from the api-linter rule names by replacing `::` and `-` with `_` and uppercasing:

| api-linter rule | buf rule ID |
|---|---|
| `core::0203::field-behavior-required` | `CORE_0203_FIELD_BEHAVIOR_REQUIRED` |
| `core::0140::lower-snake` | `CORE_0140_LOWER_SNAKE` |
| `core::0126::unspecified` | `CORE_0126_UNSPECIFIED` |

### Categories

Rules are grouped into categories by AIP number. For example, `AIP_0203` includes all rules from [AIP-203](https://google.aip.dev/203).

### Defaults

Rules in the `core::` group are enabled by default. Rules in other groups (e.g., `cloud::`) are disabled by default, matching the api-linter's own defaults.

## Versioning

Releases follow the format `v{plugin}-{linter}`, e.g., `v1.0-2.3.1`:

- `1.0` is the plugin version, bumped for changes to the plugin itself
- `2.3.1` is the upstream api-linter version the build is based on

## License

See [LICENSE](LICENSE) for details.
