# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes via webhook or desktop notification.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a scan interval and define allowed ports in a config file:

```bash
portwatch --interval 30s --config portwatch.yaml
```

Example `portwatch.yaml`:

```yaml
allowed_ports:
  - 22
  - 80
  - 443

alerts:
  webhook: "https://hooks.example.com/notify"
  desktop: true
```

When an unexpected port opens or closes, portwatch fires the configured alert immediately.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `60s` | How often to scan open ports |
| `--config` | `portwatch.yaml` | Path to config file |
| `--once` | `false` | Run a single scan and exit |

### Example Output

```
[portwatch] 2024/01/15 10:32:05 Scan started — watching 3 allowed ports
[portwatch] 2024/01/15 10:32:35 ALERT: unexpected port opened → :8080
[portwatch] 2024/01/15 10:33:05 ALERT: port closed → :443
```

## License

MIT © 2024 yourusername