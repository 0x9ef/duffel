// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"fmt"
	"net/url"
)

type (
	PlacesClient interface {
		ListPlaces(ctx context.Context, params []ListPlacesParams, requestOptions ...RequestOption) ([]*Place, error)
		Cities(ctx context.Context, requestOptions ...RequestOption) *Iter[City]
		City(ctx context.Context, id string, requestOptions ...RequestOption) (*City, error)
	}

	Place struct {
		ID              string     `json:"id"`
		Airports        []*Airport `json:"airports"`
		City            *City      `json:"city"`
		CityName        string     `json:"city_name"`
		CountryName     string     `json:"country_name"`
		IATACityCode    string     `json:"iata_city_code"`
		IATACode        string     `json:"iata_code"`
		IATACountryCode string     `json:"iata_country_code"`
		ICAOCode        string     `json:"icao_code"`
		Latitude        float64    `json:"latitude"`
		Longitude       float64    `json:"longitude"`
		Name            string     `json:"name"`
		TimeZone        string     `json:"time_zone"`
		Type            PlaceType  `json:"type"`
	}

	PlaceType string

	ListPlacesParams struct {
		// (Deprecated) A search string for finding matching Places by name.
		// This is deprecated. Please use the name query instead for equivalent behaviour.
		Query string
		// A search string for finding matching Places by name
		Name string
		// The radius, in metres, to search within
		Rad int
		// The latitude to search by
		Lat float64
		// The longitude to search by
		Long float64
	}
)

const PlaceTypeAirport = "airport"
const PlaceTypeCity = "city"

func (a *API) ListPlaces(ctx context.Context, params []ListPlacesParams, requestOptions ...RequestOption) ([]*Place, error) {
	return newRequestWithAPI[ListPlacesParams, Place](a).
		Get("/places/suggestions").
		WithParams(normalizeParams(params)...).
		WithOptions(requestOptions...).
		Slice(ctx)
}

func (a *API) Cities(ctx context.Context, requestOptions ...RequestOption) *Iter[City] {
	return newRequestWithAPI[EmptyPayload, City](a).
		Get("/air/cities").
		WithOptions(requestOptions...).
		Iter(ctx)
}

func (a *API) City(ctx context.Context, id string, requestOptions ...RequestOption) (*City, error) {
	return newRequestWithAPI[EmptyPayload, City](a).
		Getf("/air/cities/%s", id).
		WithOptions(requestOptions...).
		Single(ctx)
}

func (p ListPlacesParams) Encode(v url.Values) error {
	if p.Query != "" {
		v.Set("name", p.Query)
	}
	if p.Name != "" {
		v.Set("name", p.Name)
	}
	if p.Rad != 0 {
		v.Set("rad", fmt.Sprintf("%d", p.Rad))
	}
	if p.Lat != 0 {
		v.Set("lat", fmt.Sprintf("%f", p.Lat))
	}
	if p.Long != 0 {
		v.Set("lng", fmt.Sprintf("%f", p.Long))
	}

	return nil
}
