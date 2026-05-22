#!/usr/bin/env bash
# 开发入口（macOS / Linux），等价于 dev.bat

cd "$(dirname "$0")"

killall node 2>/dev/null || true
killall wails3 2>/dev/null || true
killall SunnyNetTools 2>/dev/null || true

export runDebug=true
export CGO_ENABLED=1

# bindings 目录若属主为 root（曾 sudo wails dev），需先：sudo chown -R "$(whoami)" frontend/bindings
if [ -d frontend/bindings/changeme ] && [ ! -w frontend/bindings/changeme ]; then
  echo "frontend/bindings 不可写（常为 root 属主）。请先执行："
  echo "  sudo chown -R \"\$(whoami)\" frontend/bindings"
  exit 1
fi
if [ ! -f frontend/bindings/changeme/Service/appmain.js ]; then
  echo "缺少 Wails bindings，正在生成…"
  wails3 generate bindings -clean=true || exit 1
fi

rm -rf frontend/node_modules/.vite frontend/dist 2>/dev/null || true

# node_modules 从 Windows 拷贝时 .bin 常缺 +x，会触发 vite Permission denied (exit 126)
if [ -d frontend/node_modules/.bin ]; then
  chmod +x frontend/node_modules/.bin/* 2>/dev/null || true
fi

exec wails3 dev -config ./build/config.yml -port 9245
