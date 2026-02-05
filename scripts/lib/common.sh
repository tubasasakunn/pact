#!/bin/bash
# Common utility functions for Pact gallery scripts

# =============================================================================
# Configuration
# =============================================================================

readonly DOCS_DIR="docs/sample"
readonly COMMIT_DIR="$DOCS_DIR/commit"
readonly SAMPLE_PACT_DIR="sample/pact"
readonly CATEGORIES=("class" "state" "flow" "sequence")

# =============================================================================
# Utility Functions
# =============================================================================

# Count SVG files in a directory
# Usage: count_svgs <directory>
count_svgs() {
  local dir="$1"
  if [ -d "$dir" ]; then
    find "$dir" -maxdepth 1 -name "*.svg" -type f 2>/dev/null | wc -l | tr -d ' '
  else
    echo "0"
  fi
}

# Get list of commit directories (sorted by newest first)
# Usage: get_commit_dirs
get_commit_dirs() {
  if [ -d "$COMMIT_DIR" ]; then
    find "$COMMIT_DIR" -mindepth 1 -maxdepth 1 -type d -printf '%T@ %p\n' 2>/dev/null | \
      sort -rn | cut -d' ' -f2- | xargs -I{} basename {}
  fi
}

# Ensure directory exists
# Usage: ensure_dir <directory>
ensure_dir() {
  local dir="$1"
  [ -d "$dir" ] || mkdir -p "$dir"
}

# Log message with timestamp
# Usage: log <message>
log() {
  echo "[$(date '+%H:%M:%S')] $*"
}

# Log error message
# Usage: log_error <message>
log_error() {
  echo "[$(date '+%H:%M:%S')] ERROR: $*" >&2
}

# Log success message
# Usage: log_success <message>
log_success() {
  echo "[$(date '+%H:%M:%S')] ✓ $*"
}
