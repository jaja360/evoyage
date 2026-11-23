# Contributing to EVoyage

Thanks for helping. EVoyage is a monorepo with a **C++23 routing/guidance engine** and a
**Go + HTMX web front-end**, connected by a **JSON-over-Unix-domain-socket contract**.

> **Architecture in one line:** the C++23 engine owns graph ingestion, waiting-time-aware
> routing and guidance; the Go + HTMX server owns HTTP/UI and region manifests; the two
> speak JSON over a Unix domain socket, always.

## Development environment

| Tool | Version | Needed for |
| --- | --- | --- |
| C++ compiler | GCC 14+ or Clang 18+ | engine |
| CMake | 3.28+ | engine build (presets) |
| Ninja | any recent | presets default generator |
| Go | as pinned in `server/go.mod` | web front-end |
| `clang-format` | 23.x (`python -m pip install clang-format==23.1.5`) | formatting gate |
| `clang-tidy` | 23.x (`python -m pip install clang-tidy==23.1.4`) | static-analysis gate |
| `gofumpt` | `go install mvdan.cc/gofumpt@latest` | Go formatting gate |
| `golangci-lint` | v2.x (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`) | Go lint gate |
| `pre-commit` | 3.7+ | local hook runner |
| `yamllint` | 1.35+ | YAML gate |

After cloning:

```bash
pre-commit install --install-hooks            # pre-commit + pre-commit hooks
pre-commit install --hook-type commit-msg     # commitlint
```

The commit-msg hook needs `@commitlint/cli` available through pre-commit's mirror; no global
Node.js install is required.

## Build and test

The repository is a monorepo: the C++ engine lives in `engine/` and the Go server in
`server/`. Build and test each half from its own directory.

```bash
# C++ engine
cd engine
cmake --preset debug             # also: release, clang-debug, asan, ubsan, tsan, fuzz
cmake --build --preset debug
ctest --preset debug

# Go front-end
cd ../server
go build ./...
go vet ./...
go test ./... -race -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -1
```

`-Werror` is part of the build definition: a warning is a build failure, not a review
comment.

## The quality-gate contract

Every pull request must pass **all** of the following. They are required status checks and
are blocking — there is no "merge now, fix later" path.

### 1. Formatting — C++

Run from the repository root:

```bash
clang-format --dry-run --Werror $(git ls-files 'engine/include/*' 'engine/src/*' 'engine/tests/*' |
  grep -E '\.(c|cc|cpp|cxx|h|hh|hpp|hxx|ipp|tpp)$')
```

The engine's `engine/.clang-format` is authoritative. To fix, drop `--dry-run --Werror` (or
run `pre-commit run clang-format --all-files`).

### 2. Formatting — Go

Run from `server/`:

```bash
gofumpt -l -d .        # non-empty output = failure
gofumpt -w .           # fix
```

### 3. Static analysis — clang-tidy

`engine/.clang-tidy` is authoritative; the script configures the engine and reads its
compilation database.

```bash
cd engine
cmake --preset debug                # needs compile_commands.json
../../.github/scripts/run-clang-tidy.sh
```

From the repository root the same script can be run directly as
`./.github/scripts/run-clang-tidy.sh`; it locates `engine/` itself.

### 4. Static analysis — Go

Run from `server/`:

```bash
golangci-lint run    # configuration: server/.golangci.yml
go vet ./...
```

### 5. YAML

Run from the repository root:

```bash
yamllint -c .github/yamllint.yml \
  $(find .github -type f \( -name '*.yml' -o -name '*.yaml' \)) \
  $(find . -maxdepth 1 -type f \( -name '*.yml' -o -name '*.yaml' \)) \
  $(find engine server -maxdepth 1 -type f \( -name '*.yml' -o -name '*.yaml' \))
```

### 6. Tests

```bash
cd engine && ctest --preset debug
cd server && go test ./... -race
```

New behaviour needs tests. Bug fixes need a regression test that fails before the fix.

### 7. Sanitizers

```bash
cd engine
for p in asan ubsan tsan; do
  cmake --preset "$p" && cmake --build --preset "$p" && ctest --preset "$p"
done
```

### 8. Coverage ratchet

`.coverage-floor` holds a single integer: the minimum total Go coverage, in percent, that
CI accepts. The check is:

```bash
cd server
go test ./... -coverprofile=coverage.out
../.github/scripts/check-coverage.sh coverage.out   # fails below the floor
```

The script resolves `.coverage-floor` from the repository root, so it works from any
directory in the tree.

Rules:

- CI fails when measured coverage is **below** the floor.
- A pull request may **raise** the floor but must never **lower** it. `.github/scripts/check-floor.sh`
  enforces this on PRs by comparing against the base branch.
- If your PR changes behaviour, add tests. Lowering the floor to get green is not an
  option; discuss it with the maintainer first if a refactor genuinely cannot keep
  coverage.

### 9. Conventional Commits

Commit messages and PR titles follow [Conventional Commits](https://www.conventionalcommits.org/):
`<type>(<scope>): <subject>`, for example `feat(planner): fold queue time into edge costs`.
Allowed types: `build`, `chore`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`,
`style`, `test`. Configuration: `.commitlintrc.yml` (`@commitlint/config-conventional`).

Squash-merge keeps `main` readable; the PR title becomes the commit subject.

## Copyright headers

Every source file starts with an SPDX identifier line:

```
SPDX-License-Identifier: AGPL-3.0-or-later
```

followed by one or more copyright lines. The policy:

- The **SPDX line is required** in every source file (C++, Go, shell, CMake, YAML).
  CI and review treat a missing SPDX line as a defect.
- The **copyright line** stays with the original author on files they create.
- A contributor creating a **new** file adds their own copyright line beneath the existing
  ones, in the form `Copyright (C) <year> <name or handle>`.
- **Ordinary edits do not modify headers.** Git history is the source of truth for who
  changed what; churning copyright lines on every touch produces merge conflicts and noise.
- A **substantial rewrite** of a file may add a copyright line or extend an existing year
  span (e.g. `2026` → `2026-2027`). Use judgement: rewriting most of a file qualifies,
  fixing a bug does not.

Signing off commits ([Developer Certificate of Origin](https://developercertificate.org/),
`git commit -s`) is optional today and encouraged for contributors sending substantial
patches. It may become a merge requirement in the future; adding a `Signed-off-by` trailer
now costs nothing and is never removed.

## Review expectations

- Keep changes scoped. Unrelated refactors belong in their own PR.
- Do not add dependencies without discussing them first. Every dependency is a supply-chain
  and licence decision: EVoyage is AGPL-3.0-or-later and must remain redistributable.
- Never commit datasets. Map extracts, compiled graphs, tiles and station databases are
  per-region build artefacts (see `.gitignore`). Providers and attribution belong in the
  region manifest and README.
- Update `README.md` attribution when you add a region or a charging-data provider.
- Keep citations in `CITATION.cff`: it is the single source of truth (GitHub generates the
  BibTeX export and the "Cite this repository" button from it). Update it when the paper
  metadata or the software author list changes.

## Maintainer checklist — branch protection (manual, one-time)

GitHub settings cannot be committed, so the maintainer configures this after the repository
is created. `main` must be protected with:

- [ ] Require a pull request before merging (at least 1 approval, dismiss stale approvals).
- [ ] Require review from Code Owners (`.github/CODEOWNERS`).
- [ ] Require status checks to pass, with these **required** and "up to date" enforced:
  - [ ] C++ CI — `gcc / Debug`, `gcc / Release`, `clang / Debug`, `clang / Release`
  - [ ] Go CI — `lint`, `test`
  - [ ] Lint — `clang-format`, `clang-tidy`, `gofumpt`, `yamllint`
  - [ ] Sanitizers — `asan`, `ubsan`, `tsan`
  - [ ] Commitlint — `pr-title`, `commits`
- [ ] Require conversation resolution before merging.
- [ ] Require linear history (squash or rebase merges only).
- [ ] Block force pushes and deletions on `main`.
- [ ] Enable GitHub private vulnerability reporting (Settings → Code security).
- [ ] Allow Renovate to automerge only after checks are green; do not exempt it from the
  checks above.

Only after these are enabled should the "required status checks" claim in this document be
considered true.

## Reporting problems

- Security issues: **never** a public issue — see [`SECURITY.md`](SECURITY.md).
- Behaviour and community: [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
