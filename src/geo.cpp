// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau
// Implementation of the geographic primitives used by the EVoyage planner.

#include "evoyage/geo.hpp"

#include <algorithm>
#include <cmath>
#include <numbers>
#include <stdexcept>
#include <string>

namespace evoyage::geo {
namespace {

constexpr double kFullTurnDegrees = 360.0;
constexpr double kHalfTurnDegrees = 180.0;
constexpr double kRadiansPerDegree = std::numbers::pi / kHalfTurnDegrees;
constexpr double kDegreesPerRadian = kHalfTurnDegrees / std::numbers::pi;

[[nodiscard]] double to_radians(double degrees) {
    return degrees * kRadiansPerDegree;
}

// Wrap an angle in degrees into [0, 360) so that computed headings stay normalized.
[[nodiscard]] double normalize_degrees(double degrees) {
    const double wrapped = std::fmod(degrees, kFullTurnDegrees);
    return wrapped < 0.0 ? wrapped + kFullTurnDegrees : wrapped;
}

}  // namespace

LatLng::LatLng(double latitude, double longitude) : lat{latitude}, lon{longitude} {
    if (!std::isfinite(lat) || std::abs(lat) > kMaxLatitudeDegrees) {
        throw std::invalid_argument{"invalid latitude (must be finite and in [-90, 90]): " +
                                    std::to_string(lat)};
    }
    if (!std::isfinite(lon) || std::abs(lon) > kMaxLongitudeDegrees) {
        throw std::invalid_argument{"invalid longitude (must be finite and in [-180, 180]): " +
                                    std::to_string(lon)};
    }
}

double haversine_distance(const LatLng& origin, const LatLng& destination, double radius_meters) {
    const double lat1 = to_radians(origin.lat);
    const double lat2 = to_radians(destination.lat);
    const double delta_lat = lat2 - lat1;
    const double delta_lon = to_radians(destination.lon - origin.lon);

    const double sin_half_lat = std::sin(delta_lat / 2.0);
    const double sin_half_lon = std::sin(delta_lon / 2.0);
    const double haversine = (sin_half_lat * sin_half_lat) +
                             (std::cos(lat1) * std::cos(lat2) * sin_half_lon * sin_half_lon);

    // Round-off can push the value just past 1 for nearly antipodal points; clamp before asin.
    return 2.0 * radius_meters * std::asin(std::sqrt(std::min(1.0, haversine)));
}

double initial_bearing(const LatLng& origin, const LatLng& destination) {
    const double lat1 = to_radians(origin.lat);
    const double lat2 = to_radians(destination.lat);
    const double delta_lon = to_radians(destination.lon - origin.lon);

    const double east_component = std::sin(delta_lon) * std::cos(lat2);
    const double north_component =
        (std::cos(lat1) * std::sin(lat2)) - (std::sin(lat1) * std::cos(lat2) * std::cos(delta_lon));
    return normalize_degrees(std::atan2(east_component, north_component) * kDegreesPerRadian);
}

double angle_between(double entry_bearing_degrees, double exit_bearing_degrees) {
    // OSRM guidance convention: fmod(entry - exit + 540, 360).
    return std::fmod(entry_bearing_degrees - exit_bearing_degrees + kHalfTurnDegrees +
                         kFullTurnDegrees,
                     kFullTurnDegrees);
}

}  // namespace evoyage::geo
