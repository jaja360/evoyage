#!/usr/bin/env bash
# Coverage-ratchet guard: on a pull request, the floor may rise but never fall.
#
# Safe to run from any directory inside the repository (CI runs it from server/).
#
# Usage: BASE_REF=origin/main ./.github/scripts/check-floor.sh
# Environment:
#   BASE_REF   git ref of the base branch (default: origin/main)
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
floor_file="${COVERAGE_FLOOR_FILE:-.coverage-floor}"
[[ "$floor_file" = /* ]] || floor_file="${repo_root}/${floor_file}"
floor_rel="${floor_file#"${repo_root}"/}"
base_ref="${BASE_REF:-origin/main}"

read_floor() {
  tr -d '[:space:]' < "$1"
}

if [[ ! -f "$floor_file" ]]; then
  echo "::error::missing $floor_file on the head branch"
  exit 1
fi

new_floor="$(read_floor "$floor_file")"
if ! [[ "$new_floor" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  echo "::error::$floor_file must contain a single number, got: '$new_floor'"
  exit 1
fi

if ! old_content="$(git show "${base_ref}:${floor_rel}" 2>/dev/null)"; then
  echo "No $floor_file on ${base_ref}; this PR introduces the ratchet baseline (${new_floor}%)."
  exit 0
fi

old_floor="$(printf '%s' "$old_content" | tr -d '[:space:]')"
if ! [[ "$old_floor" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  echo "::warning::could not parse the base floor ('$old_floor'); skipping the ratchet comparison"
  exit 0
fi

echo "Base floor (${base_ref}): ${old_floor}%"
echo "Head floor:               ${new_floor}%"

if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  {
    echo "### Coverage floor ratchet"
    echo ""
    echo "| base | head |"
    echo "| ---: | ---: |"
    echo "| ${old_floor}% | ${new_floor}% |"
    echo ""
  } >> "$GITHUB_STEP_SUMMARY"
fi

if awk -v old="$old_floor" -v new="$new_floor" 'BEGIN { exit !(new + 0 < old + 0) }'; then
  echo "::error::the coverage floor may not be lowered (${old_floor}% -> ${new_floor}%)"
  echo "Restore $floor_file; raise it only when coverage actually improved."
  exit 1
fi

echo "Floor ratchet passed."
