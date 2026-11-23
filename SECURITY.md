# Security Policy

## Supported versions

EVoyage is pre-release. Only the tip of `main` is supported; there are no maintained
release branches yet.

| Version | Supported |
| --- | --- |
| `main` | :white_check_mark: |
| Tagged releases | :x: (none yet) |
| Everything else | :x: |

## Reporting a vulnerability

**Do not open a public issue, discussion, or pull request for a security problem.**

Use GitHub's private vulnerability reporting:

1. Go to <https://github.com/jaja360/evoyage/security/advisories/new>.
2. Describe the issue: affected component (C++ engine, Go front-end, ingestion, CI),
   version/commit, reproduction steps, and impact.
3. If GitHub private reporting is unavailable, contact the maintainer
   ([@jaja360](https://github.com/jaja360)) privately and ask for a secure channel before
   sending details.

Please include a minimal reproduction where possible. Do not include real user data,
third-party API keys, or personal location histories.

## What to expect

- Acknowledgement of the report within 7 days.
- An assessment and remediation plan within 30 days, or an explanation of why the report is
  not actionable.
- Credit in the advisory once a fix is available, unless you prefer to remain anonymous.

## Scope

In scope:

- Memory-safety and denial-of-service bugs in the C++ engine (parsing of OSM/region data,
  routing, guidance).
- Injection, path traversal, XSS, SSRF or auth bypass in the Go + HTMX front-end.
- Unsafe handling of the JSON-over-Unix-socket contract between the two processes.
- Supply-chain issues in the build, CI workflows, or container images.

Out of scope:

- Vulnerabilities in third-party map or charging-station data itself (report those to the
  data provider).
- Findings that require a already-compromised host or a modified local build.
- Missing hardening headers on a development-only server.
- Automated-scanner output without a demonstrated impact.

## Disclosure

Fixes ship on `main` and are documented in the GitHub security advisory. Please allow a
reasonable embargo before publishing details.
