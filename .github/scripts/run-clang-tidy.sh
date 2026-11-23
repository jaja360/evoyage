#!/usr/bin/env bash
# Run clang-tidy over first-party C++ sources using the engine .clang-tidy.
#
# Expects engine/CMakePresets.json to define a configure preset that exports
# compile_commands.json (CLANG_TIDY_PRESET, default: debug).
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
engine_dir="${repo_root}/engine"
cd "$engine_dir"

preset="${CLANG_TIDY_PRESET:-debug}"
jobs="${CLANG_TIDY_JOBS:-$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 2)}"

compile_db="$(find . -maxdepth 4 -name compile_commands.json -not -path './build/_deps/*' -print -quit)"

if [[ -z "$compile_db" ]]; then
  if [[ ! -f CMakePresets.json ]]; then
    echo "::error::no compile_commands.json and no CMakePresets.json to generate one"
    exit 1
  fi
  echo "Generating compile_commands.json via preset '${preset}'"
  cmake --preset "$preset"
  compile_db="$(find . -maxdepth 4 -name compile_commands.json -not -path './build/_deps/*' -print -quit)"
fi

if [[ -z "$compile_db" ]]; then
  echo "::error::compile_commands.json not found after configuring; is CMAKE_EXPORT_COMPILE_COMMANDS ON?"
  exit 1
fi

build_dir="$(dirname "$compile_db")"
echo "Using compilation database: ${compile_db}"

# Only first-party translation units; skip generated code and third-party deps.
mapfile -t files < <(
  python3 - "$compile_db" "$engine_dir" <<'PY'
import json
import os
import sys

with open(sys.argv[1], encoding="utf-8") as handle:
    database = json.load(handle)

root = os.path.realpath(sys.argv[2])
prefixes = tuple(os.path.join(root, part) + os.sep for part in ("include", "src", "tests"))

seen = set()
for entry in database:
    path = entry.get("file", "")
    if not path:
        continue
    real = os.path.realpath(path)
    if real in seen or not real.startswith(prefixes):
        continue
    seen.add(real)
    print(path)
PY
)

if [[ ${#files[@]} -eq 0 ]]; then
  echo "::error::no include/, src/ or tests/ translation units in the compilation database"
  exit 1
fi

echo "Running clang-tidy on ${#files[@]} file(s) with ${jobs} job(s)"
printf '%s\0' "${files[@]}" |
  xargs -0 -n 1 -P "$jobs" clang-tidy -p "$build_dir" --quiet --warnings-as-errors='*'

echo "clang-tidy passed."
