// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau
// Minimal command-line entry point for the EVoyage planner (placeholder for the real CLI).

#include <iostream>
#include <string_view>

#ifndef EVOYAGE_VERSION
#error "EVOYAGE_VERSION must be provided by the build system"
#endif

#ifndef EVOYAGE_GIT_COMMIT
#error "EVOYAGE_GIT_COMMIT must be provided by the build system"
#endif

namespace {

constexpr std::string_view kProgramName = "evoyage-planner";

void print_usage(std::ostream& out) {
    out << "Usage: " << kProgramName << " [OPTION]\n"
        << "\n"
        << "EVoyage - an EV route planner that accounts for charging stops.\n"
        << "This is an early stub; routing is not implemented yet.\n"
        << "\n"
        << "Options:\n"
        << "  -h, --help     Print this help message and exit.\n"
        << "  -V, --version  Print the version and the git commit, then exit.\n";
}

void print_version() {
    std::cout << kProgramName << ' ' << EVOYAGE_VERSION << " (commit " << EVOYAGE_GIT_COMMIT
              << ")\n";
}

}  // namespace

int main(int argc, char** argv) {
    if (argc < 2) {
        print_usage(std::cerr);
        return 1;
    }
    if (argc > 2) {
        std::cerr << kProgramName << ": too many arguments\n\n";
        print_usage(std::cerr);
        return 1;
    }

    const std::string_view argument{argv[1]};
    if (argument == "--help" || argument == "-h") {
        print_usage(std::cout);
        return 0;
    }
    if (argument == "--version" || argument == "-V") {
        print_version();
        return 0;
    }

    std::cerr << kProgramName << ": unknown argument '" << argument << "'\n\n";
    print_usage(std::cerr);
    return 1;
}
