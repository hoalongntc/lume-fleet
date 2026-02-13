# lume-fleet

`lume-fleet` manages multiple [Lume](https://github.com/trycua/cua/tree/main/libs/lume) VMs from a single declarative `fleet.yml`.

All VM operations are executed through the `lume` CLI (no direct HTTP API calls from `lume-fleet`).

## Prerequisites

- Go 1.24+
- `lume` installed and reachable in `PATH`
- `lume serve` running locally (required by `lume` commands)

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
- `image`: macOS IPSW path/`latest` or Linux ISO path
- `vnc-port`: integer `0-65535`
- `storage`: named storage location
- `shared-dir`: host directory to share when running
- `tags`: list of tags for filtering
- `autostart`: set `false` to keep VM created/stopped on `up`

### `vnc-port` behavior

`vnc-port` is applied only during VM **creation** (`lume create --vnc-port ...`) and not during VM run/start.

- Use `0` for auto-assigned VNC port.
- Use a fixed port when you need deterministic unattended setup behavior.

### `image` behavior

For Linux VMs, `image` is mounted as ISO only on the start immediately after creation (`up` create flow). It is not mounted for later `up` runs on existing VMs.

For macOS VMs, `image` is sent as `ipsw` during create. If omitted, `lume-fleet` uses `latest`.

## Example

```yaml
defaults:
  os: macos
  cpu: 4
  memory: 8GB
  disk-size: 50GB
  unattended: tahoe
  image: latest
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
    image: ~/Downloads/ubuntu-25.10-desktop-arm64.iso
    tags: [ci, ephemeral]
```

## Notes

- Keep local overrides in `fleet.yml`; commit `fleet.yml.example` for team defaults.
- For macOS guests, `lume-fleet up` enforces the 2-VM concurrent limit.
