# lume-fleet

`lume-fleet` manages multiple [Lume](https://github.com/trycua/cua/tree/main/libs/lume) VMs from a single declarative `fleet.yml`.

## Prerequisites

- Go 1.24+
- `lume` installed and reachable in `PATH`
- Lume API server running locally:

```bash
lume serve
```

`lume-fleet` talks to `http://localhost:7777`.

## Install

Build locally:

```bash
make build
```

Install to Go bin:

```bash
make install
```

## Quick Start

1. Copy the example config:

```bash
cp fleet.yml.example fleet.yml
```

2. Edit `fleet.yml` for your machines.

3. Bring up VMs:

```bash
./lume-fleet up
```

4. Check status:

```bash
./lume-fleet status
```

## Commands

- `lume-fleet up [vm1 vm2 ...] [--tag <tag>]`
  - Creates missing VMs and starts stopped ones.
- `lume-fleet down [vm1 vm2 ...] [--tag <tag>]`
  - Stops running VMs.
- `lume-fleet destroy [vm1 vm2 ...] [--tag <tag>] [--force]`
  - Deletes VMs (`--force` required to execute).
- `lume-fleet status [--tag <tag>] [--json]`
  - Shows fleet status table or JSON.
- `lume-fleet version`
  - Prints CLI version.

Global flag:

- `--config <path>` (default: `fleet.yml`)

## Config Schema

Top-level keys:

- `defaults`: values inherited by VMs
- `vms`: map of VM name -> spec

Supported fields:

- `os`: `macos` or `linux`
- `cpu`: integer CPU count
- `memory`: size string (e.g. `4GB`, `512MB`)
- `disk-size`: size string (e.g. `50GB`)
- `unattended`: macOS unattended preset/name
- `vnc-port`: integer `0-65535`
- `storage`: named storage location
- `shared-dir`: host directory to share when running
- `tags`: list of tags for filtering
- `autostart`: set `false` to keep VM created/stopped on `up`

### `vnc-port` behavior

`vnc-port` is sent only during VM **creation** (`POST /lume/vms`) and not during VM run/start (`POST /lume/vms/:name/run`).

- Use `0` for auto-assigned VNC port.
- Use a fixed port when you need deterministic unattended setup behavior.

## Example

```yaml
defaults:
  os: macos
  cpu: 4
  memory: 8GB
  disk-size: 50GB
  unattended: tahoe
  vnc-port: 0

vms:
  dev-main:
    cpu: 8
    memory: 16GB
    vnc-port: 5901
    shared-dir: ~/Projects
    tags: [dev]

  ci-runner-1:
    os: linux
    cpu: 4
    memory: 4GB
    tags: [ci, ephemeral]
```

## Notes

- Keep local overrides in `fleet.yml`; commit `fleet.yml.example` for team defaults.
- For macOS guests, `lume-fleet up` enforces the 2-VM concurrent limit.
