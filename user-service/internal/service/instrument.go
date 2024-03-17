package service

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var tracer = otel.Tracer("github.com/zaenalarifin12/user-service/internal/service")

var meter = otel.Meter("github.com/zaenalarifin12/user-service/internal/service")

var orderCounter, _ = meter.Int64Counter("login",
	metric.WithDescription("number of order with its status"),
	metric.WithUnit("1"))

var orderStatusKey = attribute.Key("status")
