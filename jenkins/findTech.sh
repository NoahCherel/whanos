#!/bin/bash

path="$1"
RECOGNISED=()

declare -A LANG=(
    ["c"]="Makefile"
    ["java"]="app/pom.xml"
    ["javascript"]="package.json"
    ["python"]="requirements.txt"
    ["befunge"]="app/main.bf"
)

for i in "${!LANG[@]}"; do
    file="$path/${LANG[$i]}"

    if test -f "$file"; then
        RECOGNISED+=("$i")
    fi
done

if [ ${#RECOGNISED[@]} -eq 0 ]; then
    >&2 echo "Not found"
    exit 1
elif [ ${#RECOGNISED[@]} -gt 1 ]; then
    >&2 echo "found"
    exit 1
else
    echo "${RECOGNISED[0]}"
    exit 0
fi
