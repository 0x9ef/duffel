package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateBatchOfferRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/batch_offer_requests").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-create-batch-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.CreateBatchOfferRequest(ctx, CreateBatchOfferRequestInput{
		Passengers: []OfferRequestPassenger{
			{
				FamilyName: "Earhardt",
				GivenName:  "Amelia",
				Type:       PassengerTypeAdult,
			},
			{
				Age: 14,
			},
		},
		CabinClass: CabinClassEconomy,
		Slices: []OfferRequestSlice{
			{
				DepartureDate: Date(time.Now().AddDate(0, 0, 7)),
				Origin:        "JFK",
				Destination:   "AUS",
			},
		},
	})
	a.NoError(err)
	a.NotNil(data)

	a.Equal(7, data.RemainingBatches)
	a.Equal(7, data.TotalBatches)
}

func TestGetBatchOfferRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/batch_offer_requests/orq_0000AhTmH2Thpl6RrM97qK").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-batch-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")

	data, err := client.GetBatchOfferRequest(ctx, "orq_0000AhTmH2Thpl6RrM97qK")
	a.NoError(err)
	a.NotNil(data)

	a.Equal(2, data.TotalBatches)
	a.Equal(2, data.RemainingBatches)
}
