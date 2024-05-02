// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
)

type (
	AircraftClient interface {
		ListAircraft(ctx context.Context, requestOptions ...RequestOption) *Iter[Aircraft]
		GetAircraft(ctx context.Context, id string, requestOptions ...RequestOption) (*Aircraft, error)
	}
)

func (a *API) ListAircraft(ctx context.Context, requestOptions ...RequestOption) *Iter[Aircraft] {
	return newRequestWithAPI[ListAirportsParams, Aircraft](a).
		Get("/air/aircraft").
		WithOptions(requestOptions...).
		Iter(ctx)
}

func (a *API) GetAircraft(ctx context.Context, id string, requestOptions ...RequestOption) (*Aircraft, error) {
	return newRequestWithAPI[EmptyPayload, Aircraft](a).
		Getf("/air/aircraft/%s", id).
		WithOptions(requestOptions...).
		Single(ctx)
}

var _ AircraftClient = (*API)(nil)
