taskkill /F /IM node.exe
taskkill /F /IM wails3.exe
taskkill /F /IM SunnyNetTools.exe
set GOEXPERIMENT=nodwarf5
set CGO_ENABLED=1
rd /s /q frontend\node_modules\.vite
wails3 dev --port 9265
