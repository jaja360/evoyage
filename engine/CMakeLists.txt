# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2026 Jaël Champagne Gareau
# Top-level build for EVoyage, an EV route planner with charging-stop detours.

cmake_minimum_required(VERSION 3.28)

project(evoyage
    VERSION 0.1.0
    DESCRIPTION "EV route planner that accounts for charging stops"
    HOMEPAGE_URL "https://github.com/jaja360/evoyage"
    LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# EVoyage does not use C++20 modules. Turning module scanning off keeps compile_commands.json
# free of GCC-only module flags, so clang-tidy can consume it (see .clang-tidy) and builds skip
# the unnecessary standard-library module scans.
set(CMAKE_CXX_SCAN_FOR_MODULES OFF)

if(NOT CMAKE_BUILD_TYPE AND NOT CMAKE_CONFIGURATION_TYPES)
    set(CMAKE_BUILD_TYPE Debug CACHE STRING "Build type" FORCE)
    set_property(CACHE CMAKE_BUILD_TYPE PROPERTY STRINGS Debug Release RelWithDebInfo MinSizeRel)
endif()

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(EvoyageWarnings)
include(EvoyageDependencies)

enable_testing()

# Expose the current commit in --version. Configure-time capture keeps the build hermetic;
# re-run the configure step to refresh it.
set(EVOYAGE_GIT_COMMIT "unknown")
find_package(Git QUIET)
if(Git_FOUND)
    execute_process(
        COMMAND "${GIT_EXECUTABLE}" rev-parse --short=8 HEAD
        WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
        RESULT_VARIABLE _evoyage_git_result
        OUTPUT_VARIABLE _evoyage_git_commit
        OUTPUT_STRIP_TRAILING_WHITESPACE
        ERROR_QUIET)
    if(_evoyage_git_result EQUAL 0 AND _evoyage_git_commit)
        set(EVOYAGE_GIT_COMMIT "${_evoyage_git_commit}")
        execute_process(
            COMMAND "${GIT_EXECUTABLE}" diff-index --quiet HEAD --
            WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
            RESULT_VARIABLE _evoyage_git_dirty)
        if(NOT _evoyage_git_dirty EQUAL 0)
            string(APPEND EVOYAGE_GIT_COMMIT "-dirty")
        endif()
    endif()
endif()

# Core library: geographic primitives and, later, the routing algorithm itself.
add_library(evoyage_geo STATIC src/geo.cpp)
add_library(evoyage::geo ALIAS evoyage_geo)
target_include_directories(evoyage_geo
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>)
target_link_libraries(evoyage_geo PRIVATE evoyage_warnings)

# Command-line entry point. Deliberately a stub until the routing core exists.
add_executable(evoyage-planner src/main.cpp)
target_link_libraries(evoyage-planner PRIVATE evoyage::geo evoyage_warnings)
target_compile_definitions(evoyage-planner PRIVATE
    EVOYAGE_VERSION="${PROJECT_VERSION}"
    EVOYAGE_GIT_COMMIT="${EVOYAGE_GIT_COMMIT}")

option(EVOYAGE_BUILD_TESTS "Build the EVoyage unit tests" ON)
if(EVOYAGE_BUILD_TESTS)
    add_subdirectory(tests)
endif()
