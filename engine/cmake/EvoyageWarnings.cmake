# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2026 Jaël Champagne Gareau
# Shared warning flags, applied as errors to every EVoyage target.
# Sanitizer flags are injected by the CMake presets so that they stay per-configuration.

add_library(evoyage_warnings INTERFACE)

target_compile_options(evoyage_warnings INTERFACE
    $<$<COMPILE_LANG_AND_ID:CXX,GNU,Clang,AppleClang>:-Wall>
    $<$<COMPILE_LANG_AND_ID:CXX,GNU,Clang,AppleClang>:-Wextra>
    $<$<COMPILE_LANG_AND_ID:CXX,GNU,Clang,AppleClang>:-Wpedantic>
    $<$<COMPILE_LANG_AND_ID:CXX,GNU,Clang,AppleClang>:-Werror>)
