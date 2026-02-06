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

# Get list of commit directories (sorted by newest first using git commit date)
# Usage: get_commit_dirs
get_commit_dirs() {
  if [ -d "$COMMIT_DIR" ]; then
    local entries=""
    for dir in "$COMMIT_DIR"/*/; do
      [ -d "$dir" ] || continue
      local commit_hash
      commit_hash=$(basename "$dir")
      local timestamp
      # Try to get commit timestamp from git
      timestamp=$(git log -1 --format='%ct' "$commit_hash" 2>/dev/null)
      if [ -z "$timestamp" ]; then
        # Fallback to directory modification time
        timestamp=$(stat -c '%Y' "$dir" 2>/dev/null || echo "0")
      fi
      entries="${entries}${timestamp} ${commit_hash}"$'\n'
    done
    echo "$entries" | grep -v '^$' | sort -rn | awk '{print $2}'
  fi
}

# Get commit date (YYYY-MM-DD format)
# Usage: get_commit_date <commit_hash>
get_commit_date() {
  local commit_hash="$1"
  local date_str
  date_str=$(git log -1 --format='%ci' "$commit_hash" 2>/dev/null | cut -d' ' -f1)
  if [ -n "$date_str" ]; then
    echo "$date_str"
  else
    echo ""
  fi
}

# Get commit message (first line only)
# Usage: get_commit_message <commit_hash>
get_commit_message() {
  local commit_hash="$1"
  local msg
  msg=$(git log -1 --format='%s' "$commit_hash" 2>/dev/null)
  if [ -n "$msg" ]; then
    echo "$msg"
  else
    echo ""
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
