# envsnap

> CLI tool to snapshot and diff environment variables across deployments

---

## Installation

```bash
go install github.com/yourusername/envsnap@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envsnap/releases).

---

## Usage

**Take a snapshot of the current environment:**

```bash
envsnap capture --output snapshot.json
```

**Diff two snapshots:**

```bash
envsnap diff snapshot-v1.json snapshot-v2.json
```

**Example output:**

```
+ NEW_FEATURE_FLAG=true
- DEPRECATED_KEY=old_value
~ DATABASE_URL  [changed]
```

**Compare against a live environment:**

```bash
envsnap diff snapshot-v1.json --live
```

---

## Commands

| Command | Description |
|---------|-------------|
| `capture` | Save current env vars to a snapshot file |
| `diff` | Compare two snapshots or a snapshot vs live env |
| `list` | Print all variables in a snapshot |

---

## License

[MIT](LICENSE) © 2024 Your Name