// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau
// Minimal command-line entry point for the EVoyage planner (placeholder for the real CLI).

#include <exception>
#include <string_view>

#include <cstdio>
#include <print>

#ifndef EVOYAGE_VERSION
#error "EVOYAGE_VERSION must be provided by the build system"
#endif

#ifndef EVOYAGE_GIT_COMMIT
#error "EVOYAGE_GIT_COMMIT must be provided by the build system"
#endif

namespace {

constexpr const char* kProgramName = "evoyage-planner";

void print_usage(std::FILE* stream) {
    std::println(stream, "Usage: {} [OPTION]", kProgramName);
    std::println(stream);
    std::println(stream, "EVoyage - an EV route planner that accounts for charging stops.");
    std::println(stream, "This is an early stub; routing is not implemented yet.");
    std::println(stream);
    std::println(stream, "Options:");
    std::println(stream, "  -h, --help     Print this help message and exit.");
    std::println(stream, "  -V, --version  Print the version and the git commit, then exit.");
}

void print_version() {
    std::println("{} {} (commit {})", kProgramName, EVOYAGE_VERSION, EVOYAGE_GIT_COMMIT);
}

}  // namespace

int main(int argc, char** argv) {
    try {
        if (argc < 2) {
            print_usage(stderr);
            return 1;
        }
        if (argc > 2) {
            std::println(stderr, "{}: too many arguments", kProgramName);
            std::println(stderr);
            print_usage(stderr);
            return 1;
        }

        const std::string_view argument{argv[1]};
        if (argument == "--help" || argument == "-h") {
            print_usage(stdout);
            return 0;
        }
        if (argument == "--version" || argument == "-V") {
            print_version();
            return 0;
        }

        std::println(stderr, "{}: unknown argument '{}'", kProgramName, argument);
        std::println(stderr);
        print_usage(stderr);
        return 1;
    } catch (const std::exception& error) {
        // Best-effort reporting with non-throwing C stdio so no exception can escape main.
        // modernize-use-std-print is suppressed here on purpose: std::println itself throws,
        // which bugprone-exception-escape forbids inside this catch handler.
        // NOLINTBEGIN(modernize-use-std-print)
        std::fprintf(stderr, "%s: internal error: %s\n", kProgramName, error.what());
        // NOLINTEND(modernize-use-std-print)
        return 1;
    }
}
