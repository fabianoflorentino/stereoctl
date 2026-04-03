# stereoctl

stereoctl is a small CLI to detect and fix audio/video/container incompatibilities (for example, to improve compatibility with DaVinci Resolve Free).

Key features

- `convert`: convert/remux audio to AAC stereo
- `check`: analyze a file and suggest actions based on a profile
- `fix`: apply a profile (e.g. `resolve-free`) and convert/remux when necessary

Quick installation

- Requirements: Go (for development), `ffmpeg` and `ffprobe` (for runtime and integration tests).
- To install local git hooks (optional):

```bash
bash scripts/install-lefthook.sh
```

Installing system dependencies (Ubuntu/Debian):

```bash
sudo apt-get update
sudo apt-get install -y ffmpeg
```

Usage examples

- `convert` — convert/remux a file:

```bash
stereoctl convert movie.mkv
stereoctl convert movie.mkv --output fixed.mp4 --bitrate 256k
```

- `check` — analyze only and suggest actions:

```bash
stereoctl check movie.mkv
```

- `fix` — apply the `resolve-free` profile (default):

```bash
# normal mode: convert/remux
stereoctl fix movie.mkv

# preview mode: show the ffmpeg command without running it
stereoctl fix --preview movie.mkv

# batch mode: accept a directory or glob and process multiple files
stereoctl fix --batch "*.mkv"
stereoctl fix --batch /path/to/videos
```

Important flags

- `--output, -o`: specify output path (defaults to same name with `.mp4`)
- `--bitrate, -b`: audio bitrate for conversion (`convert`)
	- Default: `192k`
- `--profile, -p`: profile to apply (`fix`)
- `--preview, -n`: print the `ffmpeg` command without executing it (`fix`)
- `--batch, -B`: treat the argument as a directory/glob and process multiple files (`fix`)

Tests

Unit and integration tests (integration requires `ffmpeg` on PATH):

```bash
go test ./... -v
```

CI

There is a workflow at `.github/workflows/integration.yml` that runs tests on `ubuntu-latest` and installs `ffmpeg` before executing integration tests.

Troubleshooting

- `ffmpeg` or `ffprobe` not found: install via your package manager (apt, brew, etc.) and verify with `ffmpeg -version`.
- `lefthook` and asdf shim: if `lefthook install` fails due to an asdf shim (`No version is set ...`), run the included installer script which attempts `go install` and falls back to Homebrew; as a last resort, run the installed binary directly, for example:

```bash
# example when Homebrew installed the binary at /home/linuxbrew/.linuxbrew/bin
/home/linuxbrew/.linuxbrew/bin/lefthook install
```

Contributing

- Open issues/PRs for bugs or enhancements. Follow the commit message conventions and run local hooks before pushing.

Suggested next steps
- Add `make` targets for `hooks-install`, `test`, `build`, and packaging binaries.
- Document profiles and decision heuristics (see `internal/profiles`).

Developer commands

- Install local hooks:

```bash
make hooks-install
```

- Build and package release artifacts:

```bash
make release
```
