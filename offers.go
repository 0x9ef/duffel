// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/bojanz/currency"
)

const offerIDPrefix = "off_"

type (
	OfferClient interface {
		UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input PassengerUpdateInput, requestOptions ...RequestOption) (*OfferRequestPassenger, error)
		ListOffers(ctx context.Context, offerRequestID string, options []ListOffersParams, requestOptions ...RequestOption) *Iter[Offer]
		GetOffer(ctx context.Context, offerID string, params []GetOfferParams, requestOptions ...RequestOption) (*Offer, error)
	}

	Offer struct {
		ID                                      string                  `json:"id"`
		LiveMode                                bool                    `json:"live_mode"`
		CreatedAt                               time.Time               `json:"created_at"`
		UpdatedAt                               time.Time               `json:"updated_at"`
		ExpiresAt                               time.Time               `json:"expires_at"`
		TotalEmissionsKg                        string                  `json:"total_emissions_kg"`
		RawTotalCurrency                        string                  `json:"total_currency"`
		RawTotalAmount                          string                  `json:"total_amount"`
		RawTaxAmount                            string                  `json:"tax_amount"`
		RawTaxCurrency                          string                  `json:"tax_currency"`
		RawBaseAmount                           string                  `json:"base_amount"`
		RawBaseCurrency                         string                  `json:"base_currency"`
		Owner                                   Airline                 `json:"owner"`
		Slices                                  []Slice                 `json:"slices"`
		Passengers                              []OfferRequestPassenger `json:"passengers"`
		Partial                                 bool                    `json:"partial"`
		PassengerIdentityDocumentsRequired      bool                    `json:"passenger_identity_documents_required"`
		AllowedPassengerIdentityDocumentTypes   []string                `json:"allowed_passenger_identity_document_types"`
		SupportedPassengerIdentityDocumentTypes []string                `json:"supported_passenger_identity_document_types"`
		PaymentRequirements                     OfferPaymentRequirement `json:"payment_requirements"`
		AvailableServices                       []AvailableService      `json:"available_services"`
		Conditions                              Conditions              `json:"conditions"`
	}

	AvailableService struct {
		// Duffel's unique identifier for service
		ID               string         `json:"id"`
		MaximumQuantity  int            `json:"maximum_quantity"`
		Metadata         map[string]any `json:"metadata"`
		PassengerIDs     []string       `json:"passenger_ids"`
		SegmentIDs       []string       `json:"segment_ids"`
		RawTotalAmount   string         `json:"total_amount"`
		RawTotalCurrency string         `json:"total_currency"`

		// Possible values: "baggage"
		Type string `json:"type"`
	}

	OfferPaymentRequirement struct {
		RequiresInstantPayment  bool      `json:"requires_instant_payment"`
		PriceGuaranteeExpiresAt *DateTime `json:"price_guarantee_expires_at"`
		PaymentRequiredBy       *DateTime `json:"payment_required_by"`
	}

	ListOffersSortParam string

	ListOffersParams struct {
		Sort           ListOffersSortParam `url:"sort,omitempty"`
		MaxConnections int                 `url:"max_connections,omitempty"`
	}

	GetOfferParams struct {
		ReturnAvailableServices bool
	}
)

const (
	ListOffersSortTotalAmount   ListOffersSortParam = "total_amount"
	ListOffersSortTotalDuration ListOffersSortParam = "total_duration"
)

// UpdateOfferPassenger updates a single offer passenger.
func (a *API) UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input PassengerUpdateInput, requestOptions ...RequestOption) (*OfferRequestPassenger, error) {
	if err := validateID(offerRequestID, offerRequestIDPrefix); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/air/offers/%s/passengers/%s", offerRequestID, passengerID)
	return newRequestWithAPI[PassengerUpdateInput, OfferRequestPassenger](a).
		Patch(url, &input).
		WithOptions(requestOptions...).
		Single(ctx)
}

// ListOffers lists all the offers for an offer request. Returns an iterator.
func (a *API) ListOffers(ctx context.Context, offerRequestID string, options []ListOffersParams, requestOptions ...RequestOption) *Iter[Offer] {
	if err := validateID(offerRequestID, offerRequestIDPrefix); err != nil {
		return ErrIter[Offer](err)
	}

	return newRequestWithAPI[ListOffersParams, Offer](a).Get("/air/offers").
		WithParam("offer_request_id", offerRequestID).
		WithParams(normalizeParams(options)...).
		WithOptions(requestOptions...).
		Iter(ctx)
}

// GetOffer gets a single offer by ID.
func (a *API) GetOffer(ctx context.Context, offerID string, params []GetOfferParams, requestOptions ...RequestOption) (*Offer, error) {
	if err := validateID(offerID, offerIDPrefix); err != nil {
		return nil, err
	}

	return newRequestWithAPI[GetOfferParams, Offer](a).
		Getf("/air/offers/%s", offerID).
		WithParams(normalizeParams(params)...).
		WithOptions(requestOptions...).
		Single(ctx)
}

func (o ListOffersParams) Encode(q url.Values) error {
	if o.Sort != "" {
		q.Set("sort", string(o.Sort))
	}

	if o.MaxConnections != 0 {
		q.Set("max_connections", fmt.Sprintf("%d", o.MaxConnections))
	}

	return nil
}

func (o GetOfferParams) Encode(q url.Values) error {
	if o.ReturnAvailableServices {
		q.Set("return_available_services", "true")
	}
	return nil
}

func (o *Offer) BaseAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawBaseAmount, o.RawBaseCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

func (o *Offer) TotalAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawTotalAmount, o.RawTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

func (o *Offer) TaxAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawTaxAmount, o.RawTaxCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

// Less will sort ascending by total amount
func (o Offers) Less(i, j int) bool {
	cmp, err := o[i].TotalAmount().Cmp(o[j].TotalAmount())
	if err != nil {
		return false
	}
	return cmp < 0
}

func (o Offers) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o Offers) Len() int {
	return len(o)
}

var _ OfferClient = (*API)(nil)
