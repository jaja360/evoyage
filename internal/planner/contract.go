// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

// Package planner defines the API contract between the EVoyage web service and
// the route planner sidecar.
//
// The contract is transport-agnostic on purpose. The sidecar will later speak
// JSON over a unix socket; the structs below are the shared v0 protocol, and
// Client is the seam the web layer codes against.
package planner

import "time"

// LatLng is a WGS84 geographic coordinate in decimal degrees.
type LatLng struct {
	// Lat is the latitude in degrees, within [-90, 90].
	Lat float64 `json:"lat"`
	// Lng is the longitude in degrees, within [-180, 180].
	Lng float64 `json:"lng"`
}

// CurvePoint is one sample of a vehicle's DC charging curve.
//
// The curve maps the energy already stored in the battery (KWh) to the time
// (Minutes) needed to reach that level, so it must be strictly increasing in
// time.
type CurvePoint struct {
	// KWh is the battery level reached at this point, in kilowatt-hours.
	KWh float64 `json:"kwh"`
	// Minutes is the elapsed charging time, in minutes.
	Minutes float64 `json:"minutes"`
}

// VehicleProfile describes the vehicle used for a route.
type VehicleProfile struct {
	// BatteryKWh is the usable battery capacity in kilowatt-hours.
	BatteryKWh float64 `json:"battery_kwh"`
	// ConsumptionKWhPer100Km is the average consumption in kWh per 100 km.
	ConsumptionKWhPer100Km float64 `json:"consumption_kwh_per_100km"`
	// ChargingCurve is the DC charging curve, ordered by increasing time.
	ChargingCurve []CurvePoint `json:"charging_curve"`
}

// Options tweaks how a route is computed and returned.
type Options struct {
	// Language is the BCP 47-ish language tag for guidance text. Supported
	// values are "en" and "fr".
	Language string `json:"language"`
	// IncludeSteps asks the planner to fill RouteResponse.Steps.
	IncludeSteps bool `json:"include_steps"`
}

// RouteRequest asks the planner for a single route.
type RouteRequest struct {
	// Origin is where the trip starts.
	Origin LatLng `json:"origin"`
	// Destination is where the trip ends.
	Destination LatLng `json:"destination"`
	// AutonomyMeters is the maximum driving range on a full battery.
	AutonomyMeters float64 `json:"autonomy_meters"`
	// DepartureTime is the planned departure instant, encoded as RFC 3339.
	DepartureTime time.Time `json:"departure_time"`
	// Vehicle describes the vehicle used for the trip.
	Vehicle VehicleProfile `json:"vehicle"`
	// Options tweaks the computation and the response.
	Options Options `json:"options"`
}

// RouteLeg is one drive-and-charge segment between two charging stations.
type RouteLeg struct {
	// FromStation is the origin station of the leg; empty for the first leg.
	FromStation string `json:"from_station"`
	// ToStation is the destination station of the leg; empty for the last leg.
	ToStation string `json:"to_station"`
	// DistanceMeters is the leg length in meters.
	DistanceMeters float64 `json:"distance_meters"`
	// DriveSeconds is the time spent driving during the leg.
	DriveSeconds float64 `json:"drive_seconds"`
	// ChargeSeconds is the time spent charging at the leg's station.
	ChargeSeconds float64 `json:"charge_seconds"`
	// WaitSeconds is the time spent queuing for a charger.
	WaitSeconds float64 `json:"wait_seconds"`
}

// GuidanceStep is one turn-by-turn instruction.
type GuidanceStep struct {
	// Type is the maneuver type, for example "turn" or "arrive".
	Type string `json:"type"`
	// Modifier refines Type, for example "left" or "slight right".
	Modifier string `json:"modifier"`
	// Name is the road name, when known.
	Name string `json:"name"`
	// Ref is the road reference, for example "A-20".
	Ref string `json:"ref"`
	// ExitNumber is the exit identifier on motorways, when applicable.
	ExitNumber string `json:"exit_number"`
	// DistanceMeters is the length of the step in meters.
	DistanceMeters float64 `json:"distance_meters"`
	// BearingBefore is the heading entering the step, in degrees.
	BearingBefore int `json:"bearing_before"`
	// BearingAfter is the heading leaving the step, in degrees.
	BearingAfter int `json:"bearing_after"`
}

// RouteResponse is the planner's answer for a RouteRequest.
type RouteResponse struct {
	// DistanceMeters is the total trip distance in meters.
	DistanceMeters float64 `json:"distance_meters"`
	// DriveSeconds is the total driving time.
	DriveSeconds float64 `json:"drive_seconds"`
	// ChargeSeconds is the total charging time.
	ChargeSeconds float64 `json:"charge_seconds"`
	// WaitSeconds is the total time spent waiting for chargers.
	WaitSeconds float64 `json:"wait_seconds"`
	// Feasible reports whether a route satisfying the request was found.
	Feasible bool `json:"feasible"`
	// Legs is the ordered list of trip legs, between charging stops.
	Legs []RouteLeg `json:"legs"`
	// Steps is the turn-by-turn guidance, when requested.
	Steps []GuidanceStep `json:"steps"`
}
