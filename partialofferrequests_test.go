package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCreatePartialOfferRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/partial_offer_requests").
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-partial-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.CreatePartialOfferRequest(ctx, PartialOfferRequestInput{
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
				Origin:        "NYC",
				Destination:   "ATL",
			},
			{
				DepartureDate: Date(time.Now().AddDate(0, 0, 9)),
				Origin:        "ATL",
				Destination:   "NYC",
			},
		},
	})
	a.NoError(err)
	a.NotNil(data)

	a.Equal("973.73 GBP", data.Offers[0].TotalAmount().String())
	a.Equal("148.54 GBP", data.Offers[0].TaxAmount().String())
	a.Equal("city", data.Slices[0].OriginType)
	a.Len(data.Slices, 2) // inbound and outbound
}

func TestGetPartialOfferRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/partial_offer_requests/prq_0000AZPy1jdXi7327O8H9k").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-partial-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.GetPartialOfferRequest(ctx, "prq_0000AZPy1jdXi7327O8H9k")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("973.73 GBP", data.Offers[0].TotalAmount().String())
	a.Equal("148.54 GBP", data.Offers[0].TaxAmount().String())
	a.Equal("airport", data.Offers[0].Slices[0].DestinationType)
	a.Len(data.Slices, 2)
}

func TestListPartialOfferRequestFares(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/partial_offer_requests/prq_0000AZPy1jdXi7327O8H9k/fares").
		MatchParam("selected_partial_offer", "off_0000AZPy1jzsN4udEea6yM_0").
		MatchParam("selected_partial_offer", "off_0000AZPy1jzsN4udEea6yM_1").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-partial-offer-request-fares.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.ListPartialOfferRequestFares(ctx, "prq_0000AZPy1jdXi7327O8H9k",
		[]ListPartialOfferRequestFaresParams{
			{SelectedPartialOffer: "off_0000AZPy1jzsN4udEea6yM_0"}, {SelectedPartialOffer: "off_0000AZPy1jzsN4udEea6yM_1"},
		})
	a.NoError(err)
	a.NotNil(data)
	a.Equal("273.22 GBP", data.Offers[0].TotalAmount().String())
	a.Equal("41.68 GBP", data.Offers[0].TaxAmount().String())
	// Second offer
	a.Equal("337.78 GBP", data.Offers[1].TotalAmount().String())
	a.Equal("51.53 GBP", data.Offers[1].TaxAmount().String())
	a.Len(data.Slices, 2)
}
