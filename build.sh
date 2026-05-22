#!/bin/bash
set -euo pipefail
source ~/.bash_profile 2>/dev/null || true

# 与 Taskfile.yml 中 APP_NAME 保持一致
APP_NAME="SunnyNetTools"

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
PRODUCT_BASENAME="${APP_NAME}-MacOS"
APP_DIR="$PROJECT_DIR/bin"
APP_PATH="$APP_DIR/$APP_NAME.app"
DMG_TMP_DIR="$APP_DIR/dmg-tmp"
DMG_PATH="$APP_DIR/dist/${PRODUCT_BASENAME}.dmg"
ICON_PNG="$PROJECT_DIR/build/appicon.png"
ICON_ICNS="$PROJECT_DIR/build/darwin/icons.icns"

# 使用 macOS 自带 sips + iconutil 从 appicon.png 生成 icns（比仅依赖 wails 更可靠）
generate_macos_icns_from_png() {
  local png="$1"
  local out_icns="$2"
  local iconset="${out_icns%.icns}.iconset"

  rm -rf "$iconset" "$out_icns" "$PROJECT_DIR/build/darwin/Assets.car"
  mkdir -p "$iconset"

  sips -z 16 16 "$png" --out "$iconset/icon_16x16.png" >/dev/null
  sips -z 32 32 "$png" --out "$iconset/icon_16x16@2x.png" >/dev/null
  sips -z 32 32 "$png" --out "$iconset/icon_32x32.png" >/dev/null
  sips -z 64 64 "$png" --out "$iconset/icon_32x32@2x.png" >/dev/null
  sips -z 128 128 "$png" --out "$iconset/icon_128x128.png" >/dev/null
  sips -z 256 256 "$png" --out "$iconset/icon_128x128@2x.png" >/dev/null
  sips -z 256 256 "$png" --out "$iconset/icon_256x256.png" >/dev/null
  sips -z 512 512 "$png" --out "$iconset/icon_256x256@2x.png" >/dev/null
  sips -z 512 512 "$png" --out "$iconset/icon_512x512.png" >/dev/null
  sips -z 1024 1024 "$png" --out "$iconset/icon_512x512@2x.png" >/dev/null

  iconutil -c icns "$iconset" -o "$out_icns"
  rm -rf "$iconset"
}

if [ ! -f "$ICON_PNG" ]; then
  echo "❌ 未找到图标源文件: $ICON_PNG"
  exit 1
fi

echo "🎨 Generating icons.icns from build/appicon.png (sips + iconutil)..."
generate_macos_icns_from_png "$ICON_PNG" "$ICON_ICNS"

# Windows ico 仍走 wails（macOS 打包不依赖，失败可忽略）
(cd "$PROJECT_DIR/build" && wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico) 2>/dev/null || true

echo "🔨 Building Wails app ($APP_NAME)..."
cd "$PROJECT_DIR"
wails3 package

if [ ! -d "$APP_PATH" ]; then
  echo "❌ 未找到 $APP_PATH，请确认 wails3 package 已成功完成。"
  exit 1
fi

# 打包后再覆盖一次图标并刷新 bundle（避免 task 缓存旧 icns）
mkdir -p "$APP_PATH/Contents/Resources"
cp -f "$ICON_ICNS" "$APP_PATH/Contents/Resources/icons.icns"
touch "$APP_PATH/Contents/Resources/icons.icns" "$APP_PATH"

MACOS_BIN="$APP_PATH/Contents/MacOS/$APP_NAME"
if [ ! -f "$MACOS_BIN" ]; then
  echo "❌ 应用包内缺少可执行文件: $MACOS_BIN"
  exit 1
fi

echo "🚫 Removing quarantine flag..."
xattr -rd com.apple.quarantine "$APP_PATH" 2>/dev/null || true

echo "🔐 Ad-hoc signing app..."
codesign --deep --force --sign - "$APP_PATH"

rm -rf "$APP_DIR/dist"
mkdir -p "$APP_DIR/dist"
rm -rf "$DMG_TMP_DIR"
mkdir -p "$DMG_TMP_DIR"

# 仅放入 .app；「应用程序」快捷方式由 create-dmg 的 --app-drop-link 创建（勿再手动 ln）
cp -R "$APP_PATH" "$DMG_TMP_DIR/"

echo "💽 Creating DMG (headless, no Finder window)..."
# --skip-jenkins：不运行 AppleScript 美化，打包时不会弹出安装窗口
create-dmg \
  --skip-jenkins \
  --hdiutil-quiet \
  --volname "SunnyNetTools Installer" \
  --app-drop-link 500 185 \
  "$DMG_PATH" \
  "$DMG_TMP_DIR"

rm -rf "$DMG_TMP_DIR"

echo "✅ Done!"
echo "   运行应用: open \"$APP_PATH\""
echo "   DMG: $DMG_PATH"
echo "   产物名: ${PRODUCT_BASENAME}.dmg"
echo "   若 Finder 仍显示旧图标: rm -rf bin/$APP_NAME.app && ./build.sh && killall Finder"
