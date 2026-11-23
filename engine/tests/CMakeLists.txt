# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2026 Jaël Champagne Gareau
# Unit tests for the EVoyage C++ core.

add_executable(evoyage_tests
    main.cpp
    test_geo.cpp)

target_link_libraries(evoyage_tests PRIVATE evoyage::geo doctest::doctest evoyage_warnings)

# doctest 2.4.12 includes <ciso646>, and libstdc++ >= 15 emits a #warning for that header
# (removed from the standard in C++20). Under -Werror, Clang (and thus clang-tidy over any
# compile database) turns that third-party warning into a fatal error. -Wno-cpp silences
# #warning directives on both GCC and Clang; our own sources use no #warning directives.
# Drop this once doctest ships a fixed release.
target_compile_options(evoyage_tests PRIVATE "-Wno-cpp")

# doctest_discover_tests registers one CTest test per TEST_CASE, so failures are reported
# individually. Each test runs from the test binary directory, where any future data
# fixtures will be staged.
doctest_discover_tests(evoyage_tests
    WORKING_DIRECTORY "${CMAKE_CURRENT_BINARY_DIR}"
    TEST_PREFIX "unit.")
