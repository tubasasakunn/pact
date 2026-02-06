#!/bin/bash
# Generate GitHub Pages gallery HTML for sample SVGs
#
# Usage:
#   ./scripts/generate-gallery.sh
#
# This script generates:
#   - docs/sample/index.html       : Commit list page
#   - docs/sample/commit/*/index.html : Individual commit gallery pages

set -e

# =============================================================================
# Setup
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck source=lib/html-templates.sh
source "$SCRIPT_DIR/lib/html-templates.sh"

# =============================================================================
# Generate Main Index (Commit List)
# =============================================================================

generate_main_index() {
  log "Generating main index..."

  ensure_dir "$DOCS_DIR"

  # Start HTML
  render_main_header > "$DOCS_DIR/index.html"

  # Add commit items
  local commits
  commits=$(get_commit_dirs)

  if [ -z "$commits" ]; then
    log "No commits found in $COMMIT_DIR"
  else
    for commit in $commits; do
      local commit_path="$COMMIT_DIR/$commit"
      local class_count state_count flow_count sequence_count
      local commit_date commit_message

      class_count=$(count_svgs "$commit_path/class")
      state_count=$(count_svgs "$commit_path/state")
      flow_count=$(count_svgs "$commit_path/flow")
      sequence_count=$(count_svgs "$commit_path/sequence")
      commit_date=$(get_commit_date "$commit")
      commit_message=$(get_commit_message "$commit")

      render_commit_item "$commit" "$class_count" "$state_count" "$flow_count" "$sequence_count" "$commit_date" "$commit_message" >> "$DOCS_DIR/index.html"
    done
  fi

  # Close HTML
  render_main_footer >> "$DOCS_DIR/index.html"

  log_success "Generated $DOCS_DIR/index.html"
}

# =============================================================================
# Generate Commit Index (SVG Gallery)
# =============================================================================

generate_commit_index() {
  local commit="$1"
  local commit_path="$COMMIT_DIR/$commit"
  local output_file="$commit_path/index.html"

  log "Generating gallery for commit: $commit"

  # Count SVGs
  local class_count state_count flow_count sequence_count
  class_count=$(count_svgs "$commit_path/class")
  state_count=$(count_svgs "$commit_path/state")
  flow_count=$(count_svgs "$commit_path/flow")
  sequence_count=$(count_svgs "$commit_path/sequence")

  # Start HTML
  render_commit_header "$commit" > "$output_file"
  render_tabs "$class_count" "$state_count" "$flow_count" "$sequence_count" >> "$output_file"

  # Generate gallery sections
  local first=true
  for category in "${CATEGORIES[@]}"; do
    local is_active="false"
    [ "$first" = "true" ] && is_active="true"
    first=false

    render_gallery_start "$category" "$is_active" >> "$output_file"

    if [ -d "$commit_path/$category" ]; then
      for svg in "$commit_path/$category"/*.svg; do
        [ -f "$svg" ] || continue
        local filename
        filename=$(basename "$svg")
        render_svg_card "$category" "$filename" >> "$output_file"
      done
    fi

    render_gallery_end >> "$output_file"
  done

  # Close HTML
  render_commit_footer >> "$output_file"

  log_success "Generated $output_file"
}

# =============================================================================
# Main
# =============================================================================

main() {
  log "=== Pact Gallery Generator ==="

  # Generate main index
  generate_main_index

  # Generate individual commit pages
  local commits
  commits=$(get_commit_dirs)

  for commit in $commits; do
    generate_commit_index "$commit"
  done

  log "=== Done ==="
}

main "$@"
