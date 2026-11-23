// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package planner_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/jaja360/evoyage/server/internal/planner"
)

func TestRouteRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*planner.RouteRequest)
		wantErr []string
	}{
		{
			name:   "valid request",
			mutate: func(*planner.RouteRequest) {},
		},
		{
			name: "valid in english",
			mutate: func(r *planner.RouteRequest) {
				r.Options.Language = "en"
			},
		},
		{
			name: "origin latitude above range",
			mutate: func(r *planner.RouteRequest) {
				r.Origin.Lat = 90.1
			},
			wantErr: []string{"origin.lat"},
		},
		{
			name: "origin longitude below range",
			mutate: func(r *planner.RouteRequest) {
				r.Origin.Lng = -180.5
			},
			wantErr: []string{"origin.lng"},
		},
		{
			name: "destination latitude is NaN",
			mutate: func(r *planner.RouteRequest) {
				r.Destination.Lat = math.NaN()
			},
			wantErr: []string{"destination.lat"},
		},
		{
			name: "destination longitude is infinite",
			mutate: func(r *planner.RouteRequest) {
				r.Destination.Lng = math.Inf(1)
			},
			wantErr: []string{"destination.lng"},
		},
		{
			name: "autonomy is zero",
			mutate: func(r *planner.RouteRequest) {
				r.AutonomyMeters = 0
			},
			wantErr: []string{"autonomy_meters"},
		},
		{
			name: "autonomy is infinite",
			mutate: func(r *planner.RouteRequest) {
				r.AutonomyMeters = math.Inf(1)
			},
			wantErr: []string{"autonomy_meters"},
		},
		{
			name: "departure time is missing",
			mutate: func(r *planner.RouteRequest) {
				r.DepartureTime = time.Time{}
			},
			wantErr: []string{"departure_time"},
		},
		{
			name: "battery is zero",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.BatteryKWh = 0
			},
			wantErr: []string{"vehicle.battery_kwh"},
		},
		{
			name: "battery is NaN",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.BatteryKWh = math.NaN()
			},
			wantErr: []string{"vehicle.battery_kwh"},
		},
		{
			name: "consumption is negative",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ConsumptionKWhPer100Km = -1
			},
			wantErr: []string{"vehicle.consumption_kwh_per_100km"},
		},
		{
			name: "charging curve is empty",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = nil
			},
			wantErr: []string{"vehicle.charging_curve"},
		},
		{
			name: "curve point kWh is negative",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: -1, Minutes: 0}}
			},
			wantErr: []string{"vehicle.charging_curve[0].kwh"},
		},
		{
			name: "curve point exceeds battery",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: 10, Minutes: 0}, {KWh: 80, Minutes: 30}}
			},
			wantErr: []string{"vehicle.charging_curve[1].kwh"},
		},
		{
			name: "curve point kWh is NaN",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: math.NaN(), Minutes: 0}}
			},
			wantErr: []string{"vehicle.charging_curve[0].kwh"},
		},
		{
			name: "curve minutes are not increasing",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: 10, Minutes: 10}, {KWh: 20, Minutes: 10}}
			},
			wantErr: []string{"vehicle.charging_curve[1].minutes"},
		},
		{
			name: "curve minutes are negative",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: 10, Minutes: -1}}
			},
			wantErr: []string{"vehicle.charging_curve[0].minutes"},
		},
		{
			name: "curve minutes are NaN",
			mutate: func(r *planner.RouteRequest) {
				r.Vehicle.ChargingCurve = []planner.CurvePoint{{KWh: 10, Minutes: math.NaN()}}
			},
			wantErr: []string{"vehicle.charging_curve[0].minutes"},
		},
		{
			name: "language is unsupported",
			mutate: func(r *planner.RouteRequest) {
				r.Options.Language = "de"
			},
			wantErr: []string{"options.language"},
		},
		{
			name: "language is empty",
			mutate: func(r *planner.RouteRequest) {
				r.Options.Language = ""
			},
			wantErr: []string{"options.language"},
		},
		{
			name: "several violations are reported together",
			mutate: func(r *planner.RouteRequest) {
				r.Origin.Lat = 100
				r.Vehicle.BatteryKWh = -5
				r.Options.Language = "es"
			},
			wantErr: []string{"origin.lat", "vehicle.battery_kwh", "options.language"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := sampleRequest()
			test.mutate(&request)

			err := request.Validate()
			if len(test.wantErr) == 0 {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() = nil, want an error about %v", test.wantErr)
			}
			for _, want := range test.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Validate() = %q, want it to mention %q", err, want)
				}
			}
		})
	}
}
