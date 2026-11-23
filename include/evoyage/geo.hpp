// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau
// Geographic primitives for EVoyage: coordinates, great-circle distance, and headings.

#pragma once

namespace evoyage::geo {

// Mean Earth radius (IUGG) in meters; the default sphere used for distance computations.
inline constexpr double kEarthRadiusMeters = 6'371'008.8;

// Valid range for a latitude in decimal degrees.
inline constexpr double kMaxLatitudeDegrees = 90.0;

// Valid range for a longitude in decimal degrees.
inline constexpr double kMaxLongitudeDegrees = 180.0;

// A WGS84 coordinate in decimal degrees. Construction validates the ranges, so an object of
// this type always holds a usable coordinate.
class LatLng {
public:
    // The origin, (0, 0), which is always valid.
    LatLng() = default;

    // Throws std::invalid_argument unless the latitude is finite and in [-90, 90] and the
    // longitude is finite and in [-180, 180].
    LatLng(double latitude, double longitude);

    double lat{};  // degrees, [-90, 90]
    double lon{};  // degrees, [-180, 180]
};

// Great-circle distance between two coordinates in meters, on a sphere of the given radius.
[[nodiscard]] double haversine_distance(const LatLng& origin,
                                        const LatLng& destination,
                                        double radius_meters = kEarthRadiusMeters);

// Initial bearing (forward azimuth) from `origin` to `destination`: degrees clockwise from
// true north, normalized to [0, 360).
[[nodiscard]] double initial_bearing(const LatLng& origin, const LatLng& destination);

// Turn angle between an entry and an exit bearing, following the OSRM guidance convention:
// 180 is straight ahead, 0 is a U-turn, values in (0, 180) are right turns, and values in
// (180, 360) are left turns.
[[nodiscard]] double angle_between(double entry_bearing_degrees, double exit_bearing_degrees);

}  // namespace evoyage::geo
