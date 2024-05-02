// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
)

type (
	AirlinesClient interface {
		ListAirlines(ctx context.Context, requestOptions ...RequestOption) *Iter[Airline]
		GetAirline(ctx context.Context, id string, requestOptions ...RequestOption) (*Airline, error)
	}
)

func (a *API) ListAirlines(ctx context.Context, requestOptions ...RequestOption) *Iter[Airline] {
	return newRequestWithAPI[EmptyPayload, Airline](a).
		Get("/air/airlines").
		WithOptions(requestOptions...).
		Iter(ctx)
}

func (a *API) GetAirline(ctx context.Context, id string, requestOptions ...RequestOption) (*Airline, error) {
	return newRequestWithAPI[EmptyPayload, Airline](a).
		Getf("/air/airlines/%s", id).
		WithOptions(requestOptions...).
		Single(ctx)
}

var _ AirlinesClient = (*API)(nil)
