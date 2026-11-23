# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2026 Jaël Champagne Gareau
# Third-party dependencies, fetched from source so that every checkout builds identical code.
#
# No package manager (vcpkg/Conan) yet: doctest is header-only and trivially pinned by tag,
# so the extra machinery buys nothing. A package manager will be introduced together with
# the first real external dependency (e.g. a routing engine or a serialization library).

include(FetchContent)

FetchContent_Declare(doctest
    SYSTEM
    GIT_REPOSITORY https://github.com/doctest/doctest.git
    GIT_TAG v2.4.12
    GIT_SHALLOW TRUE)

set(DOCTEST_NO_INSTALL ON CACHE BOOL "Do not install doctest with EVoyage" FORCE)
set(DOCTEST_WITH_TESTS OFF CACHE BOOL "Do not build doctest's own tests" FORCE)
set(DOCTEST_WITH_EXAMPLES OFF CACHE BOOL "Do not build doctest's own examples" FORCE)
FetchContent_MakeAvailable(doctest)

# Make doctest_discover_tests() available regardless of how doctest was brought in.
if(NOT COMMAND doctest_discover_tests)
    list(APPEND CMAKE_MODULE_PATH "${doctest_SOURCE_DIR}/scripts/cmake")
    include(doctest)
endif()
