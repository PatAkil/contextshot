#!/usr/bin/env bash
# Build Contextshot.app — a minimal macOS app bundle that runs the menu-bar
# daemon as an agent (no Dock icon) via LSUIElement.
set -euo pipefail

APP="dist/Contextshot.app"
MACOS_DIR="$APP/Contents/MacOS"

rm -rf "$APP"
mkdir -p "$MACOS_DIR"

go build -o "$MACOS_DIR/contextshot" ./cmd/...

cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>Contextshot</string>
	<key>CFBundleIdentifier</key>
	<string>com.patakil.contextshot</string>
	<key>CFBundleExecutable</key>
	<string>contextshot</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>0.1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSUIElement</key>
	<true/>
	<key>LSMinimumSystemVersion</key>
	<string>13.0</string>
</dict>
</plist>
PLIST

# Ad-hoc sign with a stable identifier matching the bundle id, so the code
# signature and LaunchServices identity agree (required for TCC permissions
# like Accessibility and Screen Recording to actually apply).
codesign --force --sign - --identifier com.patakil.contextshot "$MACOS_DIR/contextshot"

echo "Built $APP"
