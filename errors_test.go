package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGatewayTimeoutError(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/partial_offer_requests").
		Reply(504).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/504-gateway-timeout.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	_, err := client.CreatePartialOfferRequest(ctx, PartialOfferRequestInput{
		Passengers: []OfferRequestPassenger{
			{
				Type: PassengerTypeAdult,
			},
		},
		CabinClass: CabinClassEconomy,
		Slices: []OfferRequestSlice{
			{
				DepartureDate: Date(time.Now().AddDate(0, 0, 7)),
				Origin:        "STN",
				Destination:   "LHR",
			},
		},
	})
	a.Equal(true, IsErrorCode(err, GatewayTimeout), "is 504 gateway timeout")
}
