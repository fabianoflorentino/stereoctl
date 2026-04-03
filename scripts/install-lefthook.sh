#!/usr/bin/env bash
set -euo pipefail

# Instala/ativa lefthook para este repositório.
# - tenta usar `lefthook` se já estiver no PATH
# - tenta `go install github.com/evilmartians/lefthook/cmd/lefthook@latest` se `go` disponível
# - executa `lefthook install` ao final

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Installing lefthook for repo at $ROOT_DIR"

if command -v lefthook >/dev/null 2>&1; then
  echo "lefthook found at: $(command -v lefthook)"
  echo "Attempting: lefthook install"
  if lefthook install; then
    echo "lefthook hooks installed."
    exit 0
  fi
  echo "Warning: 'lefthook install' failed. Will attempt to reinstall lefthook via 'go install' if possible."
fi

if command -v go >/dev/null 2>&1; then
  echo "Installing lefthook via 'go install'..."
  GO_PKG="github.com/evilmartians/lefthook/cmd/lefthook@latest"
  if go install "$GO_PKG"; then
    echo "go install succeeded"
    if command -v lefthook >/dev/null 2>&1; then
      lefthook install
      echo "lefthook hooks installed."
      exit 0
    fi

    # try common bin locations
    if [ -n "${GOBIN:-}" ] && [ -x "${GOBIN}/lefthook" ]; then
      echo "Found lefthook at $GOBIN/lefthook (not on PATH). Running install using it."
      PATH="$GOBIN:$PATH" lefthook install
      echo "lefthook hooks installed."
      exit 0
    fi

    GOPATH=$(go env GOPATH 2>/dev/null || echo "")
    if [ -n "$GOPATH" ] && [ -x "$GOPATH/bin/lefthook" ]; then
      echo "Found lefthook at $GOPATH/bin/lefthook. Running install."
      PATH="$GOPATH/bin:$PATH" lefthook install
      echo "lefthook hooks installed."
      exit 0
    fi

    echo "lefthook binary installed but not on PATH. Add '$GOPATH/bin' or '${GOBIN:-}' to PATH and rerun '$0' or run 'lefthook install' manually."
    exit 0
  else
    echo "go install failed. Trying Homebrew if available..."
    if command -v brew >/dev/null 2>&1; then
      echo "Installing lefthook via Homebrew..."
      if brew install evilmartians/lefthook/lefthook || brew install lefthook; then
        echo "Homebrew install succeeded. Running lefthook install."
        BREW_PREFIX=$(brew --prefix 2>/dev/null || echo "")
        if [ -n "$BREW_PREFIX" ] && [ -x "$BREW_PREFIX/bin/lefthook" ]; then
          echo "Running $BREW_PREFIX/bin/lefthook install"
          "$BREW_PREFIX/bin/lefthook" install || echo "lefthook install failed after brew install; run '$BREW_PREFIX/bin/lefthook install' or 'lefthook install' manually."
          exit 0
        fi
        lefthook install || echo "lefthook install failed after brew install; run 'lefthook install' manually."
        exit 0
      fi
    fi

    echo "Please install lefthook manually. See: https://github.com/evilmartians/lefthook"
    exit 1
  fi
fi

echo "Neither 'lefthook' nor 'go' were found. Install lefthook manually: https://github.com/evilmartians/lefthook"
exit 1
