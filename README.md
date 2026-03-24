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
    # Enable all AIP rules
    - AIP
    # Or enable specific categories/rules
    # - AIP_CORE
    # - AIP_0203
    # - AIP_0140_LOWER_SNAKE
  except:
    # Disable specific rules
    - AIP_0191_JAVA_PACKAGE
plugins:
  - plugin: buf-plugin-gapi-lint
```

### Listing available rules

To see all rules and categories provided by the plugin:

```bash
buf config ls-lint-rules
```

### Rule naming

Each rule from the api-linter is exposed as a buf lint rule. The group prefix (`core::`, `client-libraries::`) is dropped and all rules are prefixed with `AIP_`:

| api-linter rule | buf rule ID |
|---|---|
| `core::0203::field-behavior-required` | `AIP_0203_FIELD_BEHAVIOR_REQUIRED` |
| `core::0140::lower-snake` | `AIP_0140_LOWER_SNAKE` |
| `client-libraries::4232::repeated-fields` | `AIP_4232_REPEATED_FIELDS` |

### Categories

Rules are organized into several category levels:

| Category | Description |
|---|---|
| `AIP` | All rules (catch-all) |
| `AIP_CORE` | All `core::` rules |
| `AIP_CLIENT_LIBRARIES` | All `client-libraries::` rules (AIP-4232) |
| `AIP_0203`, `AIP_0140`, etc. | Rules for a specific AIP number |

### Defaults

Rules in the `core::` group are enabled by default. `client-libraries::` rules are disabled by default, matching the api-linter's own defaults.

## Versioning

Release versions mirror the upstream api-linter version. For example, `v2.3.1` of this plugin wraps api-linter `v2.3.1`.

## License

See [COPYING](COPYING) for details.
