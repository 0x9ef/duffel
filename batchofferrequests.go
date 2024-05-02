package duffel

import (
	"context"
)

type (
	BatchOfferRequestClient interface {
		CreateBatchOfferRequest(ctx context.Context, requestInput CreateBatchOfferRequestInput) (*BatchOfferRequest, error)
		GetBatchOfferRequest(ctx context.Context, id string) (*BatchOfferRequest, error)
	}

	CreateBatchOfferRequestInput struct {
		// The passengers who want to travel. If you specify an age for a passenger, the type may differ for the same passenger in different offers due to airline's different rules. e.g. one airline may treat a 14 year old as an adult, and another as a young adult. You may only specify an age or a type – not both.
		Passengers []OfferRequestPassenger `json:"passengers" url:"-"`
		// The slices that make up this offer request. One-way journeys can be expressed using one slice, whereas return trips will need two.
		Slices []OfferRequestSlice `json:"slices" url:"-"`
		// The cabin that the passengers want to travel in
		CabinClass CabinClass `json:"cabin_class" url:"-"`
		// The maximum number of connections within any slice of the offer. For example 0 means a direct flight which will have a single segment within each slice and 1 means a maximum of two segments within each slice of the offer.
		MaxConnections *int `json:"max_connections,omitempty" url:"-"`
		// The maximum amount of time in milliseconds to wait for each airline to respond
		SupplierTimeout int `json:"-" url:"supplier_timeout,omitempty"`
	}

	BatchOfferRequest struct {
		TotalBatches     int      `json:"total_batches"`
		RemainingBatches int      `json:"remaining_batches"`
		ID               string   `json:"id"`
		Offers           []Offer  `json:"offers,omitempty"`
		CreatedAt        DateTime `json:"created_at"`
	}
)

func (a *API) CreateBatchOfferRequest(ctx context.Context, requestInput CreateBatchOfferRequestInput) (*BatchOfferRequest, error) {
	return newRequestWithAPI[CreateBatchOfferRequestInput, BatchOfferRequest](a).
		Post("/air/batch_offer_requests", &requestInput).
		Single(ctx)
}

func (a *API) GetBatchOfferRequest(ctx context.Context, id string) (*BatchOfferRequest, error) {
	return newRequestWithAPI[EmptyPayload, BatchOfferRequest](a).
		Getf("/air/batch_offer_requests/%s", id).
		Single(ctx)
}
