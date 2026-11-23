# EVoyage

[![CI — C++](https://github.com/jaja360/evoyage/actions/workflows/ci-cpp.yml/badge.svg)](https://github.com/jaja360/evoyage/actions/workflows/ci-cpp.yml)
[![CI — Go](https://github.com/jaja360/evoyage/actions/workflows/ci-go.yml/badge.svg)](https://github.com/jaja360/evoyage/actions/workflows/ci-go.yml)
[![Lint](https://github.com/jaja360/evoyage/actions/workflows/lint.yml/badge.svg)](https://github.com/jaja360/evoyage/actions/workflows/lint.yml)
[![Sanitizers](https://github.com/jaja360/evoyage/actions/workflows/sanitize.yml/badge.svg)](https://github.com/jaja360/evoyage/actions/workflows/sanitize.yml)
[![License: AGPL-3.0-or-later](https://img.shields.io/badge/license-AGPL--3.0--or--later-blue.svg)](LICENSE)
[![C++23](https://img.shields.io/badge/C%2B%2B-23-blue.svg)](https://en.cppreference.com/w/cpp/23)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev)

**EVoyage is an electric-vehicle route planner that treats waiting time at charging
stations as a first-class routing cost.**

Most EV routers answer "where do I plug in?" and ignore the queue at the plug. EVoyage
folds day×hour occupancy priors — the `f·g` charging-demand model from the original
research — directly into the shortest-path search, so a route can prefer a station that is
likely to be free over one that is closer but busy. The result is a plan that accounts for
the time actually spent charging, not just the kWh.

EVoyage is a ground-up rewrite of the UQAM master's project
([Champagne Gareau et al., SIGSPATIAL 2019](https://doi.org/10.1145/3347146.3359064)).
It is **region-agnostic** (OpenStreetMap-based), **low-memory** by design, and
**quality-gated from commit #1**.

## How it is built

| Layer | Stack | Role |
| --- | --- | --- |
| Routing & guidance engine | C++23, CMake presets, [doctest](https://github.com/doctest/doctest) | Graph ingestion, waiting-time-aware planner, manoeuvre guidance |
| Web front-end | Go + [HTMX](https://htmx.org/) | HTTP API, map UI, region manifests |
| Contract | JSON over a Unix domain socket | The Go front-end talks to the C++ engine; no in-process bridging |

The engine is **self-contained**: EVoyage no longer depends on an external routing engine
(OSRM is replaced by its own guidance and map-matching module) and ships only the data a
region actually needs.

## Project status

> **Early scaffold.** The repository layout, CI, licensing and contributor tooling are in
> place; the engine and web layers are under active construction. Nothing here is ready
> for production routing yet. APIs, presets and data formats will change without notice
> until the first tagged release.

## Quick start (development)

Prerequisites: a C++23 compiler (GCC 14+ or Clang 18+), CMake 3.28+, Ninja, Go (version
pinned by `go.mod`), and `clang-format`/`gofumpt` for the local hooks.

```bash
# C++ engine — configure, build, test (presets are defined in CMakePresets.json)
cmake --preset debug
cmake --build --preset debug
ctest --preset debug

# Clang and Release lanes
cmake --preset clang-debug && cmake --build --preset clang-debug && ctest --preset clang-debug
cmake --preset release && cmake --build --preset release && ctest --preset release

# Sanitizer builds
cmake --preset asan && cmake --build --preset asan && ctest --preset asan

# Go front-end
go build ./cmd/...
go test ./... -race
go run ./cmd/evoyage-server

# Everything the CI checks, locally
pre-commit run --all-files
```

Build artefacts land in `build/`, which is ignored by Git. Do not commit datasets: map
extracts (`*.osm.pbf`), compiled graphs (`graph.bin`), tiles and station databases are
generated per region and are intentionally excluded.

## Roadmap

1. **Ingestion** — OSM extract → compact regional graph, charging-station registry,
   day×hour occupancy priors from historical sessions.
2. **Planner core** — waiting-time-aware `f·g` cost model, multi-objective search
   (time/energy/queue), memory-bounded graph layout.
3. **Guidance** — self-contained turn-by-turn manoeuvre generation and map matching
   (retiring the OSRM dependency).
4. **Web** — Go + HTMX UI, region manifests, Unix-socket JSON API.
5. **Releases** — reproducible region builds, benchmarks, versioned data formats.

## Attribution

### Map data

Map data © [OpenStreetMap](https://www.openstreetmap.org/copyright) contributors, available
under the [Open Database License (ODbL) 1.0](https://opendatacommons.org/licenses/odbl/1-0/).
Derived graphs, tiles and region manifests produced by EVoyage from OSM data remain
subject to the ODbL.

### Charging-station data

Charging-station attributes come from per-region providers. Each region manifest records
the providers it uses, and the UI credits them. Providers currently planned or supported:

| Region | Provider | Attribution / licence |
| --- | --- | --- |
| _to be filled per region_ | _provider name_ | _provider attribution text and licence link_ |

Add a row here — and the matching entry in the region manifest — before enabling a region.
Never mix providers in one region without listing each of them.

### Prior work

EVoyage is a rewrite of research published as:

> Jaël Champagne Gareau, Éric Beaudry, and Vladimir Makarenkov. *An Efficient Electric
> Vehicle Path-Planner That Considers the Waiting Time.* In Proceedings of the 27th ACM
> SIGSPATIAL International Conference on Advances in Geographic Information Systems
> (SIGSPATIAL '19), 2019. DOI: [10.1145/3347146.3359064](https://doi.org/10.1145/3347146.3359064).

If you use EVoyage in academic work, please cite it via [`CITATION.cff`](CITATION.cff).

## License

- **Source code:** [GNU Affero General Public License v3.0 or later](LICENSE)
  (AGPL-3.0-or-later). Because EVoyage is network software, running a modified version as
  a service obliges you to offer the corresponding source to its users.
- **Derived map data:** the ODbL applies to databases derived from OpenStreetMap, including
  the graphs and region manifests EVoyage generates. ODbL and AGPL are separate grants:
  the code is AGPL, the OSM-derived data is ODbL.
- **Route output:** a route returned by EVoyage is a *Produced Work* in ODbL terms (§4.3).
  You may publish and use it, provided you attribute the map data as described above.
  Route output does not inherit the AGPL of the engine.

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the contribution and quality-gate contract.
