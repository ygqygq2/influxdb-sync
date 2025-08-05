#!/usr/bin/env bash
binary=$1

echo $binary
upx $binary || echo "upx [$binary] failed"
