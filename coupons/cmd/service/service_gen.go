// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package service

import (
	endpoint "github.com/GDGVIT/Project-Hades/coupons/pkg/endpoint"
	http1 "github.com/GDGVIT/Project-Hades/coupons/pkg/http"
	endpoint1 "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	http "github.com/go-kit/kit/transport/http"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
)

func createService(endpoints endpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	initHttpHandler(endpoints, g)
	return g
}
func defaultHttpOptions(logger log.Logger, tracer opentracinggo.Tracer) map[string][]http.ServerOption {
	options := map[string][]http.ServerOption{
		"CreateSchema": {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "CreateSchema", logger))},
		"DeleteCoupon": {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "DeleteCoupon", logger))},
		"DeleteSchema": {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "DeleteSchema", logger))},
		"MarkPresent":  {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "MarkPresent", logger))},
		"RedeemCoupon": {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "RedeemCoupon", logger))},
		"ViewSchema":   {http.ServerErrorEncoder(http1.ErrorEncoder), http.ServerErrorLogger(logger), http.ServerBefore(opentracing.HTTPToContext(tracer, "ViewSchema", logger))},
	}
	return options
}
func addEndpointMiddlewareToAllMethods(mw map[string][]endpoint1.Middleware, m endpoint1.Middleware) {
	methods := []string{"CreateSchema", "MarkPresent", "RedeemCoupon", "DeleteCoupon", "DeleteSchema", "ViewSchema"}
	for _, v := range methods {
		mw[v] = append(mw[v], m)
	}
}
