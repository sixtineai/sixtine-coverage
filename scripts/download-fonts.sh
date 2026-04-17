#!/usr/bin/env bash
#
# Downloads Inter, JetBrains Mono and Material Icons Round woff2 files
# from Google Fonts into src/assets/fonts/ so no external requests are
# needed at runtime (works behind strict CSP).
#
# Usage: ./scripts/download-fonts.sh
#
set -euo pipefail

FONT_DIR="$(cd "$(dirname "$0")/.." && pwd)/src/assets/fonts"
mkdir -p "$FONT_DIR"

UA="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

fetch_latin_url() {
  local css_url="$1"
  curl -sS -H "User-Agent: $UA" "$css_url" \
    | awk '/^\/\* latin \*\/$/{found=1} found && /url\(/{print; found=0}' \
    | grep -oP 'url\(\K[^)]+\.woff2' \
    | head -1
}

echo "Fetching Inter latin woff2 URL..."
INTER_URL=$(fetch_latin_url "https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap")
echo "  -> $INTER_URL"
curl -sS -o "$FONT_DIR/inter-latin.woff2" "$INTER_URL"

echo "Fetching JetBrains Mono latin woff2 URL..."
JBMONO_URL=$(fetch_latin_url "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap")
echo "  -> $JBMONO_URL"
curl -sS -o "$FONT_DIR/jetbrains-mono-latin.woff2" "$JBMONO_URL"

echo "Fetching Material Icons Round woff2 URL..."
ICONS_URL=$(curl -sS -H "User-Agent: $UA" "https://fonts.googleapis.com/icon?family=Material+Icons+Round" \
  | grep -oP 'url\(\K[^)]+\.woff2' | head -1)
echo "  -> $ICONS_URL"
curl -sS -o "$FONT_DIR/material-icons-round.woff2" "$ICONS_URL"

echo "Done. Fonts saved to $FONT_DIR"
ls -lh "$FONT_DIR"
