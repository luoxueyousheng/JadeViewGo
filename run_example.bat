@echo off
chcp 65001 >nul
rem 一键运行 example：自动把 MinGW gcc 加进 PATH 并开启 cgo
setlocal
set "PATH=C:\msys64\mingw64\bin;%PATH%"
set "CGO_ENABLED=1"
cd /d "%~dp0"

echo [info] gcc:
where gcc
echo [info] 正在构建并运行 example...
go run ./example

echo.
echo [info] 程序已退出（exit=%errorlevel%）
pause
