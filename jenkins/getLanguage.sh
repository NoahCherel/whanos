#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 [folder]"
    echo "folder: folder path (default: current folder)"
    exit 1
fi

declare -A EXTENSIONS=(
    ["c"]="Makefile"
    ["java"]="app/pom.xml"
    ["javascript"]="package.json"
    ["python"]="requirements.txt"
    ["befunge"]="app/main.bf"
)

folder_path="$1"
FOUND=()

for key in "${!EXTENSIONS[@]}"; do
    file_path="$folder_path/${EXTENSIONS[$key]}"

    if test -f "$file_path"; then
        FOUND+=("$key")
    fi
done

if [ ${#FOUND[@]} -eq 0 ]; then
    >&2 echo "Not found"
    exit 1
elif [ ${#FOUND[@]} -gt 1 ]; then
    >&2 echo "found"
    exit 1
else
    echo "${FOUND[0]}"
    exit 0
fi
