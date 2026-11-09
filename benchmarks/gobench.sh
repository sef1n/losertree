#!/bin/bash
set -euo pipefail

COUNT="${COUNT:-1}"
BENCHNAME="${BENCHNAME:-.}"

if [[ $# -eq 0 ]]; then
  echo "usage: $0 <pkg>..." >&2
  exit 2
fi

results=()
for arg in "$@"; do
  n=$(basename "$arg")
  out="$n"".bench"
  results+=("$out")

  echo "=== go test -bench=. $arg → $out ===" >&2

  go test -bench="$BENCHNAME" "$arg" -benchmem -run='^$' --count="$COUNT" 2>&1 | sed -u '/^pkg:/d' | tee "$out"
done

echo "=== benchstat (${#results[@]} files) ===" >&2
benchstat "${results[@]}"
