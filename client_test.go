// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestClientError(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)
	gock.New("https://api.duffel.com/air/offer_requests").
		MatchParam("return_offers", "false").
		Reply(400).
		File("fixtures/400-bad-request.json")
	defer gock.Off()

	client := New("duffel_test_123")
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{
		ReturnOffers: false,
	})
	a.Error(err)
	a.Nil(data)

	a.Equal("duffel: The airline responded with an unexpected error, please contact support", err.Error())

	derr := err.(*DuffelError)
	a.True(derr.IsType(AirlineError))
	a.True(derr.IsCode(AirlineUnknown))
	a.True(IsErrorType(err, AirlineError))
	a.True(IsErrorCode(err, AirlineUnknown))
	a.True(ErrIsRetryable(err))

	reqId, ok := RequestIDFromError(err)
	a.True(ok)
	a.Equal("FZW0H3HdJwKk5HMAAKxB", reqId)
}

func TestClientErrorBadGateway(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)
	gock.New("https://api.duffel.com/air/offer_requests").
		Reply(502).
		AddHeader("Content-Type", "text/html").
		File("fixtures/502-bad-gateway.html")
	defer gock.Off()

	client := New("duffel_test_123")
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{
		ReturnOffers: true,
	})
	a.Error(err)
	a.Nil(data)
	a.Equal("duffel: An internal server error occurred. Please try again later.", err.Error())
}

func TestClientRetry(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)
	gock.New("https://api.duffel.com/air/offer_requests").
		Persist().
		Reply(502).
		AddHeader("Content-Type", "text/html").
		File("fixtures/502-bad-gateway.html")
	defer gock.Off()

	client := New("duffel_test_123",
		WithRetry(3, time.Second, time.Second*5, ExponentalBackoff),
		WithRetryCondition(func(resp *http.Response, err error) bool {
			return err != nil
		}),
	)
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{ReturnOffers: true})
	a.Error(err)
	a.Nil(data)
	a.Equal("duffel: An internal server error occurred. Please try again later.", err.Error())
}
