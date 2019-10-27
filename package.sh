#!/usr/bin/env bash
rm -rf dist
mkdir -p dist/ui

GOOS=linux go build -o dist/saltshakers_linux github.com/alittlebrighter/saltshakers
GOOS=darwin go build -o dist/saltshakers_macos github.com/alittlebrighter/saltshakers
GOOS=windows go build -o dist/saltshakers_win.exe github.com/alittlebrighter/saltshakers

(cd ui/spa; npm run build)
cp -r ui/spa/dist/* dist/ui/
