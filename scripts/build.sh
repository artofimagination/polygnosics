#!/bin/bash
set -e

go build -i -race -o "$1"
chmod +x "$1"