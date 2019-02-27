// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package endpoint

import (
	service "github.com/GDGVIT/Project-Hades/coupons/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateSchemaEndpoint endpoint.Endpoint
	MarkPresentEndpoint  endpoint.Endpoint
	RedeemCouponEndpoint endpoint.Endpoint
	DeleteCouponEndpoint endpoint.Endpoint
	DeleteSchemaEndpoint endpoint.Endpoint
	ViewSchemaEndpoint   endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.CouponsService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateSchemaEndpoint: MakeCreateSchemaEndpoint(s),
		DeleteCouponEndpoint: MakeDeleteCouponEndpoint(s),
		DeleteSchemaEndpoint: MakeDeleteSchemaEndpoint(s),
		MarkPresentEndpoint:  MakeMarkPresentEndpoint(s),
		RedeemCouponEndpoint: MakeRedeemCouponEndpoint(s),
		ViewSchemaEndpoint:   MakeViewSchemaEndpoint(s),
	}
	for _, m := range mdw["CreateSchema"] {
		eps.CreateSchemaEndpoint = m(eps.CreateSchemaEndpoint)
	}
	for _, m := range mdw["MarkPresent"] {
		eps.MarkPresentEndpoint = m(eps.MarkPresentEndpoint)
	}
	for _, m := range mdw["RedeemCoupon"] {
		eps.RedeemCouponEndpoint = m(eps.RedeemCouponEndpoint)
	}
	for _, m := range mdw["DeleteCoupon"] {
		eps.DeleteCouponEndpoint = m(eps.DeleteCouponEndpoint)
	}
	for _, m := range mdw["DeleteSchema"] {
		eps.DeleteSchemaEndpoint = m(eps.DeleteSchemaEndpoint)
	}
	for _, m := range mdw["ViewSchema"] {
		eps.ViewSchemaEndpoint = m(eps.ViewSchemaEndpoint)
	}
	return eps
}
