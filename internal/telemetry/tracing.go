// Package telemetry initialises some instrumentation for the webapp.
// The aim is to make as monitoring as seamless and transparent as possible,
// that developers do not need to manually call any functions.
package telemetry

import (
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"os"
	"sync"

	"github.com/kaphos/webapp/internal/log"
)

var jaegerLogger = log.Get("JAEG")
var exporter *texporter.Exporter
var once sync.Once
var jaegerInitialised = false

func initJaeger() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		jaegerLogger.Warn("No project ID provided. Not initialising.")
		return
	}

	var err error
	exporter, err = texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		jaegerLogger.Error("Error connecting. Not initialising.")
		return
	}

	jaegerLogger.Info("Connected to Jaeger.")
	jaegerInitialised = true
}

// newTracerProvider is a wrapper around trace.NewTracerProvider,
// while abstracting away repeated/automatically configured settings.
func newTracerProvider(serviceName string) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("env", os.Getenv("ENV")),
		)),
	)
}

// NewTracer creates a new tracer provider, and returns a new tracer
// as well. Abstracts away repeated/automatically configured settings.
// Also automatically initialises a Jaeger connection, if `JAEGER_URL`
// is set, the first time this function is called.
func NewTracer(tracerName, serviceName string) oteltrace.Tracer {
	once.Do(initJaeger)

	if !jaegerInitialised {
		return otel.Tracer(tracerName)
	}

	tp := newTracerProvider(serviceName)
	return tp.Tracer(tracerName)
}
