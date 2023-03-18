#!/bin/bash

dir_path=$(cd "$(dirname "$0")" && pwd)

cd "$dir_path/../client" || exit

echo "Building bot..."
start=$(date +%s)
go build -ldflags '-w -s' -o chatgpt-bot .
end=$(date +%s)
take=$(( end - start ))
echo "Done in ${take} seconds."
