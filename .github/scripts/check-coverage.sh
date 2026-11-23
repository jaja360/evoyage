#!/usr/bin/env bash
# Fail if total Go coverage is below the value in the repository .coverage-floor.
#
# Safe to run from any directory inside the repository (CI runs it from server/).
#
# Usage: ./.github/scripts/check-coverage.sh [coverage.out]
# Environment:
#   COVERAGE_FLOOR       override the floor value (CI reads .coverage-floor and exports it)
#   COVERAGE_FLOOR_FILE  path to the floor file (default: <repo-root>/.coverage-floor)
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
profile="${1:-coverage.out}"
floor_file="${COVERAGE_FLOOR_FILE:-.coverage-floor}"
[[ "$floor_file" = /* ]] || floor_file="${repo_root}/${floor_file}"
floor="${COVERAGE_FLOOR:-}"

if [[ -z "$floor" ]]; then
  if [[ ! -f "$floor_file" ]]; then
    echo "::error::missing coverage floor file: $floor_file"
    exit 1
  fi
  floor="$(tr -d '[:space:]' < "$floor_file")"
fi

if [[ ! -f "$profile" ]]; then
  echo "::error::missing coverage profile: $profile"
  echo "::error::run 'go test ./... -coverprofile=$profile' first"
  exit 1
fi

if ! [[ "$floor" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  echo "::error::$floor_file must contain a single number, got: '$floor'"
  exit 1
fi

if ! total="$(go tool cover -func="$profile" | awk '/^total:/ { sub(/%/, "", $NF); print $NF }')"; then
  echo "::error::'go tool cover -func=$profile' failed; run it from the module root"
  exit 1
fi
if ! [[ "$total" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  echo "::error::could not parse total coverage from $profile"
  exit 1
fi

echo "Measured total coverage: ${total}%"
echo "Required floor:          ${floor}%"

if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  {
    echo "### Go coverage"
    echo ""
    echo "| measured | floor |"
    echo "| ---: | ---: |"
    echo "| ${total}% | ${floor}% |"
    echo ""
  } >> "$GITHUB_STEP_SUMMARY"
fi

if awk -v measured="$total" -v floor="$floor" 'BEGIN { exit !(measured + 0 < floor + 0) }'; then
  echo "::error::coverage ${total}% is below the ${floor}% floor"
  echo "Add tests, or raise the floor only if coverage improved. Do not lower it."
  exit 1
fi

echo "Coverage gate passed."
