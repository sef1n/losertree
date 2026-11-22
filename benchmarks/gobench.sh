#!/bin/bash
set -euo pipefail

COUNT="${COUNT:-1}"
BENCHNAME="${BENCHNAME:-.}"
BENCHTIME="${BENCHTIME:-1x}"

if [[ $# -eq 0 ]]; then
  echo "usage: $0 <pkg>..." >&2
  exit 2
fi

if command -v benchstat >/dev/null; then
  BENCHSTAT=benchstat
elif [[ -x "$(go env GOPATH)/bin/benchstat" ]]; then
  BENCHSTAT="$(go env GOPATH)/bin/benchstat"
else
  echo "install: go install golang.org/x/perf/cmd/benchstat@latest" >&2
  exit 1
fi

results=()
for arg in "$@"; do
  n=$(basename "$arg")
  out="$n"".bench"
  results+=("$out")

  echo "=== go test -bench=. $arg → $out ===" >&2

  go test -bench="$BENCHNAME" "$arg" -benchmem -run='^$' -benchtime="$BENCHTIME" --count="$COUNT" 2>&1 | sed -u '/^pkg:/d' | tee "$out"
done

echo "=== benchstat (${#results[@]} files) ===" >&2
"$BENCHSTAT" "${results[@]}"
