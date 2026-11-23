// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package planner

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Valid coordinate bounds in decimal degrees.
const (
	minLatitude  = -90
	maxLatitude  = 90
	minLongitude = -180
	maxLongitude = 180
)

// supportedLanguages lists the languages the planner can produce guidance in.
var supportedLanguages = []string{"en", "fr"}

// Validate reports every reason why r cannot be planned. It returns nil when r
// is well formed, and an error joining all violations otherwise.
func (r RouteRequest) Validate() error {
	var errs []error

	errs = append(errs, r.Origin.validate("origin")...)
	errs = append(errs, r.Destination.validate("destination")...)

	if !isPositive(r.AutonomyMeters) {
		errs = append(errs, fmt.Errorf("autonomy_meters: must be positive, got %v", r.AutonomyMeters))
	}
	if r.DepartureTime.IsZero() {
		errs = append(errs, errors.New("departure_time: must be set"))
	}

	errs = append(errs, r.Vehicle.validate()...)

	if !containsString(supportedLanguages, r.Options.Language) {
		errs = append(errs, fmt.Errorf("options.language: %q is not supported, want one of %s",
			r.Options.Language, strings.Join(supportedLanguages, ", ")))
	}

	return errors.Join(errs...)
}

// validate returns the coordinate violations of p, labeled with field.
func (p LatLng) validate(field string) []error {
	var errs []error

	if !inRange(p.Lat, minLatitude, maxLatitude) {
		errs = append(errs, fmt.Errorf("%s.lat: must be within [%d, %d], got %v",
			field, minLatitude, maxLatitude, p.Lat))
	}
	if !inRange(p.Lng, minLongitude, maxLongitude) {
		errs = append(errs, fmt.Errorf("%s.lng: must be within [%d, %d], got %v",
			field, minLongitude, maxLongitude, p.Lng))
	}

	return errs
}

// validate returns the violations of the vehicle profile.
func (v VehicleProfile) validate() []error {
	var errs []error

	if !isPositive(v.BatteryKWh) {
		errs = append(errs, fmt.Errorf("vehicle.battery_kwh: must be positive, got %v", v.BatteryKWh))
	}
	if !isPositive(v.ConsumptionKWhPer100Km) {
		errs = append(errs, fmt.Errorf("vehicle.consumption_kwh_per_100km: must be positive, got %v",
			v.ConsumptionKWhPer100Km))
	}
	if len(v.ChargingCurve) == 0 {
		errs = append(errs, errors.New("vehicle.charging_curve: must contain at least one point"))
	}

	prevMinutes := -1.0
	for i, point := range v.ChargingCurve {
		if math.IsNaN(point.KWh) || point.KWh < 0 || (isPositive(v.BatteryKWh) && point.KWh > v.BatteryKWh) {
			errs = append(errs, fmt.Errorf(
				"vehicle.charging_curve[%d].kwh: must be within [0, battery_kwh], got %v", i, point.KWh))
		}
		if math.IsNaN(point.Minutes) || point.Minutes <= prevMinutes {
			errs = append(errs, fmt.Errorf(
				"vehicle.charging_curve[%d].minutes: must increase strictly from the previous point, got %v",
				i, point.Minutes))
		}
		prevMinutes = point.Minutes
	}

	return errs
}

// isPositive reports whether f is a finite, strictly positive number.
func isPositive(f float64) bool {
	return f > 0 && !math.IsInf(f, 0)
}

// inRange reports whether v is finite and within [minV, maxV].
func inRange(v, minV, maxV float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= minV && v <= maxV
}

// containsString reports whether values contains want.
func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
