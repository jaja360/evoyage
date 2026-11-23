// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau
// Unit tests for the geographic primitives (distance, bearing, turn angle, validation).

#include "evoyage/geo.hpp"

#include <doctest/doctest.h>

#include <array>
#include <cmath>
#include <limits>
#include <stdexcept>

namespace {

using evoyage::geo::angle_between;
using evoyage::geo::haversine_distance;
using evoyage::geo::initial_bearing;
using evoyage::geo::kEarthRadiusMeters;
using evoyage::geo::LatLng;

constexpr double kTolerance = 1e-9;

}  // namespace

// Note: keep TEST_CASE names free of '['. doctest 2.4.12's CTest discovery script misparses
// such names and emits one merged, invalid test entry.

TEST_CASE("haversine_distance: Montreal to Quebec City is about 233 km") {
    const LatLng montreal{45.5019, -73.5674};
    const LatLng quebec_city{46.8139, -71.2080};

    const double distance = haversine_distance(montreal, quebec_city);
    CHECK(distance > 228'000.0);
    CHECK(distance < 238'000.0);
    CHECK(distance == doctest::Approx(233'000.0).epsilon(2e-2));
}

TEST_CASE("haversine_distance: zero for coincident points and linear in the radius") {
    const LatLng paris{48.8566, 2.3522};
    CHECK(haversine_distance(paris, paris) == 0.0);

    const LatLng london{51.5074, -0.1278};
    const double on_earth = haversine_distance(paris, london);
    CHECK(on_earth > 300'000.0);
    CHECK(on_earth < 400'000.0);

    const double on_doubled_sphere = haversine_distance(paris, london, 2.0 * kEarthRadiusMeters);
    CHECK(on_doubled_sphere == doctest::Approx(2.0 * on_earth).epsilon(1e-12));
}

TEST_CASE("initial_bearing: cardinal directions") {
    const LatLng origin{0.0, 0.0};

    CHECK(initial_bearing(origin, LatLng(10.0, 0.0)) == doctest::Approx(0.0).epsilon(kTolerance));
    CHECK(initial_bearing(origin, LatLng(0.0, 10.0)) == doctest::Approx(90.0).epsilon(kTolerance));
    CHECK(initial_bearing(origin, LatLng(-10.0, 0.0)) ==
          doctest::Approx(180.0).epsilon(kTolerance));
    CHECK(initial_bearing(origin, LatLng(0.0, -10.0)) ==
          doctest::Approx(270.0).epsilon(kTolerance));

    // The north-east diagonal sits between north and east.
    const double northeast = initial_bearing(origin, LatLng(10.0, 10.0));
    CHECK(northeast > 44.0);
    CHECK(northeast < 46.0);

    // Montreal to Quebec City heads roughly north-north-east.
    const double toward_quebec =
        initial_bearing(LatLng(45.5019, -73.5674), LatLng(46.8139, -71.2080));
    CHECK(toward_quebec > 45.0);
    CHECK(toward_quebec < 65.0);
}

TEST_CASE("initial_bearing: always normalized to a single turn") {
    const std::array<LatLng, 4> points{
        LatLng{45.5, -73.5},
        LatLng{-33.8, 151.2},
        LatLng{0.0, 179.9},
        LatLng{0.0, -179.9},
    };

    for (const LatLng& lhs : points) {
        for (const LatLng& rhs : points) {
            const double bearing = initial_bearing(lhs, rhs);
            CHECK(bearing >= 0.0);
            CHECK(bearing < 360.0);
        }
    }
}

TEST_CASE("angle_between: OSRM guidance convention") {
    // Straight ahead.
    CHECK(angle_between(0.0, 0.0) == doctest::Approx(180.0).epsilon(kTolerance));
    CHECK(angle_between(45.0, 45.0) == doctest::Approx(180.0).epsilon(kTolerance));

    // Right turns are in (0, 180).
    CHECK(angle_between(0.0, 90.0) == doctest::Approx(90.0).epsilon(kTolerance));
    CHECK(angle_between(350.0, 80.0) == doctest::Approx(90.0).epsilon(kTolerance));
    CHECK(angle_between(350.0, 10.0) == doctest::Approx(160.0).epsilon(kTolerance));

    // U-turns are 0.
    CHECK(angle_between(0.0, 180.0) == doctest::Approx(0.0).epsilon(kTolerance));
    CHECK(angle_between(90.0, 270.0) == doctest::Approx(0.0).epsilon(kTolerance));

    // Left turns are in (180, 360).
    CHECK(angle_between(0.0, 270.0) == doctest::Approx(270.0).epsilon(kTolerance));
    CHECK(angle_between(90.0, 0.0) == doctest::Approx(270.0).epsilon(kTolerance));
    CHECK(angle_between(10.0, 350.0) == doctest::Approx(200.0).epsilon(kTolerance));
}

TEST_CASE("angle_between: normalized result and complementary angles") {
    constexpr std::array<double, 8> bearings{0.0, 17.0, 45.0, 90.0, 180.0, 270.0, 359.0, 359.9};

    for (const double entry : bearings) {
        for (const double exit : bearings) {
            const double angle = angle_between(entry, exit);
            CHECK(angle >= 0.0);
            CHECK(angle < 360.0);

            const double total = angle + angle_between(exit, entry);
            CHECK(std::fmod(total, 360.0) == doctest::Approx(0.0).epsilon(kTolerance));
        }
    }
}

TEST_CASE("haversine_distance: symmetric") {
    const std::array<LatLng, 4> points{
        LatLng{45.5019, -73.5674},
        LatLng{46.8139, -71.2080},
        LatLng{-33.8688, 151.2093},
        LatLng{0.0, 0.0},
    };

    for (const LatLng& lhs : points) {
        for (const LatLng& rhs : points) {
            CHECK(haversine_distance(lhs, rhs) ==
                  doctest::Approx(haversine_distance(rhs, lhs)).epsilon(1e-12));
        }
    }
}

TEST_CASE("LatLng: rejects invalid coordinates") {
    CHECK_THROWS_AS(LatLng(91.0, 0.0), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(-90.5, 0.0), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(0.0, 181.0), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(0.0, -180.5), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(std::numeric_limits<double>::quiet_NaN(), 0.0), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(0.0, std::numeric_limits<double>::quiet_NaN()), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(std::numeric_limits<double>::infinity(), 0.0), std::invalid_argument);
    CHECK_THROWS_AS(LatLng(-std::numeric_limits<double>::infinity(), 0.0), std::invalid_argument);
}

TEST_CASE("LatLng: accepts boundary coordinates and preserves them") {
    CHECK_NOTHROW(LatLng(90.0, 180.0));
    CHECK_NOTHROW(LatLng(-90.0, -180.0));

    const LatLng montreal{45.5019, -73.5674};
    CHECK(montreal.lat == 45.5019);
    CHECK(montreal.lon == -73.5674);

    const LatLng origin;
    CHECK(origin.lat == 0.0);
    CHECK(origin.lon == 0.0);
}
