#!/bin/bash
# HTML template functions for Pact gallery
# Uses shared CSS from docs/assets/css/main.css

# =============================================================================
# Main Index Page (Commit List)
# =============================================================================

# Generate main index header
render_main_header() {
  cat << 'EOF'
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sample Gallery - Pact</title>
  <link rel="stylesheet" href="../assets/css/main.css">
</head>
<body>
  <nav class="navbar">
    <div class="nav-container">
      <a href="../" class="nav-logo">
        <span class="logo-icon">&#9670;</span>
        <span class="logo-text">Pact</span>
      </a>
      <div class="nav-links">
        <a href="../specification.html" class="nav-link">Specification</a>
        <a href="./" class="nav-link active">Gallery</a>
        <a href="https://github.com/tubasasakunn/pact" class="nav-link nav-link-external" target="_blank">
          GitHub
          <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
            <path d="M3.5 3a.5.5 0 0 0 0 1h3.793L2.146 9.146a.5.5 0 0 0 .708.708L8 4.707V8.5a.5.5 0 0 0 1 0v-5a.5.5 0 0 0-.5-.5h-5z"/>
          </svg>
        </a>
      </div>
    </div>
  </nav>

  <header class="page-header">
    <div class="container">
      <h1 style="font-size: 1.75rem; margin-bottom: 0.5rem;">Sample Gallery</h1>
      <p style="color: var(--color-text-secondary);">Generated SVG diagrams from .pact sample files</p>
    </div>
  </header>

  <main style="padding: var(--space-xl);">
    <div class="container" style="max-width: 900px;">
      <ul class="commit-list">
EOF
}

# Generate a single commit item
# Usage: render_commit_item <commit_id> <class_count> <state_count> <flow_count> <sequence_count> <date> <message>
render_commit_item() {
  local commit="$1"
  local class_count="$2"
  local state_count="$3"
  local flow_count="$4"
  local sequence_count="$5"
  local date="$6"
  local message="$7"
  local total=$((class_count + state_count + flow_count + sequence_count))

  # Escape HTML special characters in commit message
  message=$(echo "$message" | sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g; s/"/\&quot;/g')

  cat << EOF
        <li class="commit-item">
          <a href="commit/$commit/" class="commit-link">
            <div class="commit-icon">
              <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="4"/>
                <path d="M12 2v6M12 16v6"/>
              </svg>
            </div>
            <div class="commit-info">
              <div class="commit-header-row">
                <div class="commit-id">$commit</div>
                <div class="commit-date">$date</div>
              </div>
              <div class="commit-message">$message</div>
              <div class="commit-badges">
                <span class="badge">class: $class_count</span>
                <span class="badge">state: $state_count</span>
                <span class="badge">flow: $flow_count</span>
                <span class="badge">sequence: $sequence_count</span>
              </div>
            </div>
            <div class="commit-arrow">&rarr;</div>
          </a>
        </li>
EOF
}

# Generate main index footer
render_main_footer() {
  cat << 'EOF'
      </ul>
    </div>
  </main>

  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <p>&copy; 2024-2026 Pact Project</p>
        <div class="footer-links">
          <a href="https://github.com/tubasasakunn/pact" target="_blank">GitHub</a>
          <a href="../specification.html">Specification</a>
          <a href="./">Gallery</a>
        </div>
      </div>
    </div>
  </footer>
</body>
</html>
EOF
}

# =============================================================================
# Commit Detail Page (SVG Gallery)
# =============================================================================

# Generate commit page header
# Usage: render_commit_header <commit_id>
render_commit_header() {
  local commit="$1"

  cat << EOF
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>$commit - Pact Sample Gallery</title>
  <link rel="stylesheet" href="../../../assets/css/main.css">
</head>
<body>
  <nav class="navbar">
    <div class="nav-container">
      <a href="../../../" class="nav-logo">
        <span class="logo-icon">&#9670;</span>
        <span class="logo-text">Pact</span>
      </a>
      <div class="nav-links">
        <a href="../../../specification.html" class="nav-link">Specification</a>
        <a href="../../" class="nav-link active">Gallery</a>
        <a href="https://github.com/tubasasakunn/pact" class="nav-link nav-link-external" target="_blank">
          GitHub
          <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
            <path d="M3.5 3a.5.5 0 0 0 0 1h3.793L2.146 9.146a.5.5 0 0 0 .708.708L8 4.707V8.5a.5.5 0 0 0 1 0v-5a.5.5 0 0 0-.5-.5h-5z"/>
          </svg>
        </a>
      </div>
    </div>
  </nav>

  <header class="page-header">
    <div class="page-header-content">
      <a href="../../" class="back-link">&larr;</a>
      <h1 class="page-title">Commit <code>$commit</code></h1>
    </div>
  </header>

  <div class="tabs-container">
    <div class="tabs">
EOF
}

# Generate tab buttons
# Usage: render_tabs <class_count> <state_count> <flow_count> <sequence_count>
render_tabs() {
  local class_count="$1"
  local state_count="$2"
  local flow_count="$3"
  local sequence_count="$4"

  cat << EOF
      <button class="tab active" data-target="class">Class <span class="count">$class_count</span></button>
      <button class="tab" data-target="state">State <span class="count">$state_count</span></button>
      <button class="tab" data-target="flow">Flow <span class="count">$flow_count</span></button>
      <button class="tab" data-target="sequence">Sequence <span class="count">$sequence_count</span></button>
    </div>
  </div>

  <main style="padding: var(--space-xl);">
    <div class="container">
EOF
}

# Generate gallery section start
# Usage: render_gallery_start <category> [is_active]
render_gallery_start() {
  local category="$1"
  local is_active="$2"
  local active_class=""
  [ "$is_active" = "true" ] && active_class=" active"

  cat << EOF
      <section id="$category" class="gallery$active_class">
        <div class="gallery-grid">
EOF
}

# Generate SVG card
# Usage: render_svg_card <category> <filename>
render_svg_card() {
  local category="$1"
  local filename="$2"
  local alt="${filename%.svg}"
  alt="${alt//_/ }"

  cat << EOF
          <div class="svg-card">
            <div class="svg-card-header">$filename</div>
            <div class="svg-card-body"><img src="$category/$filename" alt="$alt" loading="lazy"></div>
            <div class="svg-card-footer"><a href="$category/$filename" target="_blank">Open in new tab</a></div>
          </div>
EOF
}

# Generate gallery section end
render_gallery_end() {
  cat << 'EOF'
        </div>
      </section>
EOF
}

# Generate commit page footer with JavaScript
render_commit_footer() {
  cat << 'EOF'
    </div>
  </main>

  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <p>&copy; 2024-2026 Pact Project</p>
        <div class="footer-links">
          <a href="https://github.com/tubasasakunn/pact" target="_blank">GitHub</a>
          <a href="../../../specification.html">Specification</a>
          <a href="../../">Gallery</a>
        </div>
      </div>
    </div>
  </footer>

  <script>
    document.querySelectorAll('.tab').forEach(tab => {
      tab.addEventListener('click', () => {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.gallery').forEach(g => g.classList.remove('active'));
        tab.classList.add('active');
        document.getElementById(tab.dataset.target).classList.add('active');
      });
    });
  </script>
</body>
</html>
EOF
}
