// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListPlaces(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/places/suggestions").
		MatchParam("query", "Lond").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-places-suggestion.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	places, err := client.ListPlaces(ctx, []ListPlacesParams{{Name: "Lond"}})
	a.NoError(err)
	a.NotNil(places)

	a.Equal("Heathrow", places[0].Name)
	a.Equal("EGLL", places[0].ICAOCode)
	a.Equal("London", places[0].City.Name)
	a.Equal("London", places[0].CityName)
	a.Equal("Heathrow", places[0].Airports[0].Name)
}

func TestListPlacesByLongitude(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/places/suggestions").
		MatchParam("lat", "51.885508").
		MatchParam("lng", "0.236933").
		MatchParam("rad", "100").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-places-by-cord.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	places, err := client.ListPlaces(ctx, []ListPlacesParams{
		{
			Lat: 51.885508,
		},
		{
			Long: 0.236933,
		},
		{
			Rad: 100,
		},
	})
	a.NoError(err)
	a.NotNil(places)

	a.Equal("London Stansted Airport", places[0].Name)
	a.Equal("London", places[0].CityName)
	a.Equal("EGSS", places[0].ICAOCode)
	a.Equal("STN", places[0].IATACode)
}
