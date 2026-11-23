// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package planner_test

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/jaja360/evoyage/internal/planner"
)

// sampleRequest returns a fully populated, valid request.
func sampleRequest() planner.RouteRequest {
	return planner.RouteRequest{
		Origin:         planner.LatLng{Lat: 45.5019, Lng: -73.5674},
		Destination:    planner.LatLng{Lat: 46.8139, Lng: -71.208},
		AutonomyMeters: 350_000,
		DepartureTime:  time.Date(2026, time.May, 1, 9, 30, 0, 0, time.UTC),
		Vehicle: planner.VehicleProfile{
			BatteryKWh:             77,
			ConsumptionKWhPer100Km: 17.5,
			ChargingCurve: []planner.CurvePoint{
				{KWh: 0, Minutes: 0},
				{KWh: 20, Minutes: 8},
				{KWh: 77, Minutes: 42},
			},
		},
		Options: planner.Options{Language: "fr", IncludeSteps: true},
	}
}

// sampleResponse returns a fully populated response.
func sampleResponse() planner.RouteResponse {
	return planner.RouteResponse{
		DistanceMeters: 253_000,
		DriveSeconds:   9_000,
		ChargeSeconds:  1_800,
		WaitSeconds:    300,
		Feasible:       true,
		Legs: []planner.RouteLeg{
			{
				FromStation:    "",
				ToStation:      "Drummondville",
				DistanceMeters: 101_000,
				DriveSeconds:   3_600,
				ChargeSeconds:  1_200,
				WaitSeconds:    0,
			},
			{
				FromStation:    "Drummondville",
				ToStation:      "",
				DistanceMeters: 152_000,
				DriveSeconds:   5_400,
				ChargeSeconds:  600,
				WaitSeconds:    300,
			},
		},
		Steps: []planner.GuidanceStep{
			{
				Type:           "depart",
				Modifier:       "",
				Name:           "Rue Saint-Denis",
				Ref:            "",
				ExitNumber:     "",
				DistanceMeters: 450,
				BearingBefore:  0,
				BearingAfter:   15,
			},
			{
				Type:           "turn",
				Modifier:       "slight right",
				Name:           "Autoroute 20",
				Ref:            "A-20",
				ExitNumber:     "42",
				DistanceMeters: 12_300,
				BearingBefore:  15,
				BearingAfter:   75,
			},
		},
	}
}

func TestRouteRequestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	request := sampleRequest()

	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded planner.RouteRequest
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	reEncoded, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("re-marshal: %v", err)
	}

	if string(encoded) != string(reEncoded) {
		t.Fatalf("round trip changed the payload:\n first: %s\nsecond: %s", encoded, reEncoded)
	}

	if !decoded.DepartureTime.Equal(request.DepartureTime.UTC()) {
		t.Errorf("departure time = %s, want %s", decoded.DepartureTime, request.DepartureTime)
	}
}

func TestRouteRequestJSONFieldNames(t *testing.T) {
	t.Parallel()

	var fields map[string]json.RawMessage
	decodeInto(t, sampleRequest(), &fields)

	assertKeys(t, "route request", fields, []string{
		"autonomy_meters", "departure_time", "destination", "options", "origin", "vehicle",
	})

	var vehicle map[string]json.RawMessage
	if err := json.Unmarshal(fields["vehicle"], &vehicle); err != nil {
		t.Fatalf("decode vehicle: %v", err)
	}
	assertKeys(t, "vehicle", vehicle, []string{"battery_kwh", "charging_curve", "consumption_kwh_per_100km"})

	var curve []map[string]json.RawMessage
	if err := json.Unmarshal(vehicle["charging_curve"], &curve); err != nil {
		t.Fatalf("decode charging curve: %v", err)
	}
	if len(curve) != 3 {
		t.Fatalf("charging curve has %d points, want 3", len(curve))
	}
	assertKeys(t, "curve point", curve[0], []string{"kwh", "minutes"})

	var options map[string]json.RawMessage
	if err := json.Unmarshal(fields["options"], &options); err != nil {
		t.Fatalf("decode options: %v", err)
	}
	assertKeys(t, "options", options, []string{"include_steps", "language"})

	if got, want := string(fields["departure_time"]), `"2026-05-01T09:30:00Z"`; got != want {
		t.Errorf("departure_time = %s, want %s", got, want)
	}
}

func TestRouteResponseJSONRoundTrip(t *testing.T) {
	t.Parallel()

	response := sampleResponse()

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded planner.RouteResponse
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(decoded, response) {
		t.Fatalf("round trip mismatch:\n got: %+v\nwant: %+v", decoded, response)
	}
}

func TestRouteResponseJSONFieldNames(t *testing.T) {
	t.Parallel()

	var fields map[string]json.RawMessage
	decodeInto(t, sampleResponse(), &fields)

	assertKeys(t, "route response", fields, []string{
		"charge_seconds", "distance_meters", "drive_seconds", "feasible", "legs", "steps", "wait_seconds",
	})

	var legs []map[string]json.RawMessage
	if err := json.Unmarshal(fields["legs"], &legs); err != nil {
		t.Fatalf("decode legs: %v", err)
	}
	assertKeys(t, "route leg", legs[0], []string{
		"charge_seconds", "distance_meters", "drive_seconds", "from_station", "to_station", "wait_seconds",
	})

	var steps []map[string]json.RawMessage
	if err := json.Unmarshal(fields["steps"], &steps); err != nil {
		t.Fatalf("decode steps: %v", err)
	}
	assertKeys(t, "guidance step", steps[0], []string{
		"bearing_after", "bearing_before", "distance_meters", "exit_number", "modifier", "name", "ref", "type",
	})
}

// decodeInto marshals value and unmarshals it into target.
func decodeInto(t *testing.T, value, target any) {
	t.Helper()

	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := json.Unmarshal(encoded, target); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

// assertKeys fails unless object has exactly the wanted keys.
func assertKeys(t *testing.T, label string, object map[string]json.RawMessage, want []string) {
	t.Helper()

	got := make([]string, 0, len(object))
	for key := range object {
		got = append(got, key)
	}
	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s keys = %v, want %v", label, got, want)
	}
}
