## What changed

<!-- One or two sentences. Link the issue: "Closes #123". -->

## Why

<!-- The problem this solves. Motivation over mechanism. -->

## How it was tested

<!-- Commands run, cases covered, datasets/regions used. -->

## Checklist

- [ ] PR title follows Conventional Commits (`feat(scope): subject`).
- [ ] Tests added or updated for the new behaviour (bug fixes include a regression test).
- [ ] Coverage floor respected: `./.github/scripts/check-coverage.sh` passes and
      `.coverage-floor` was not lowered.
- [ ] Linters run locally: `clang-format`, `clang-tidy`, `gofumpt`, `golangci-lint`,
      `yamllint` (or `pre-commit run --all-files`).
- [ ] Sanitizers pass for the affected C++ code (`asan`, `ubsan`, `tsan` presets).
- [ ] No new dependency without prior discussion (licence + supply chain reviewed).
- [ ] Docs updated (README attribution for new regions/data providers, CONTRIBUTING if a
      gate changed).
- [ ] No generated data committed (`*.osm.pbf`, `graph.bin`, tiles, station DBs).
