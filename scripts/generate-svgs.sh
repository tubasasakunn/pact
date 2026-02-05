#!/bin/bash
# Generate SVG diagrams from .pact files
#
# Usage:
#   ./scripts/generate-svgs.sh <commit_id>
#   ./scripts/generate-svgs.sh <commit_id> <output_base_dir>
#
# Examples:
#   ./scripts/generate-svgs.sh abc1234
#   ./scripts/generate-svgs.sh abc1234 docs/sample/commit

set -e

# =============================================================================
# Setup
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

# =============================================================================
# Configuration
# =============================================================================

COMMIT_ID="${1:-$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')}"
OUTPUT_BASE="${2:-$COMMIT_DIR}"
OUTPUT_DIR="$OUTPUT_BASE/$COMMIT_ID"
PACT_BIN="${PACT_BIN:-$PROJECT_ROOT/bin/pact}"

# =============================================================================
# Validation
# =============================================================================

if [ ! -x "$PACT_BIN" ]; then
  log_error "Pact binary not found at: $PACT_BIN"
  log_error "Run 'make build' first."
  exit 1
fi

if [ ! -d "$PROJECT_ROOT/$SAMPLE_PACT_DIR" ]; then
  log_error "Sample directory not found: $PROJECT_ROOT/$SAMPLE_PACT_DIR"
  exit 1
fi

# =============================================================================
# Main
# =============================================================================

log "Generating SVGs for commit: $COMMIT_ID"
log "Output directory: $OUTPUT_DIR"

# Create output directories
for category in "${CATEGORIES[@]}"; do
  ensure_dir "$OUTPUT_DIR/$category"
done

# Generate SVGs for each category
for category in "${CATEGORIES[@]}"; do
  pact_dir="$PROJECT_ROOT/$SAMPLE_PACT_DIR/$category"

  if [ ! -d "$pact_dir" ]; then
    log "Skipping $category (no source directory)"
    continue
  fi

  file_count=0
  for pact_file in "$pact_dir"/*.pact; do
    [ -f "$pact_file" ] || continue

    if "$PACT_BIN" generate -o "$OUTPUT_DIR/$category/" -t "$category" "$pact_file" 2>/dev/null; then
      ((file_count++)) || true
    else
      log_error "Failed to generate: $pact_file"
    fi
  done

  log "  $category: $file_count files processed"
done

# Summary
total=0
for category in "${CATEGORIES[@]}"; do
  count=$(count_svgs "$OUTPUT_DIR/$category")
  total=$((total + count))
done

log_success "Generated $total SVG files in $OUTPUT_DIR"
