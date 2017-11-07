package gosiris

import (
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/opentracing/opentracing-go"
	"context"
)

var (
	initialized bool
	ctx         context.Context
	span        opentracing.Span
	collector   zipkin.Collector
)

type ZipkinOptions struct {
	Url      string
	Debug    bool
	HostPort string
	SameSpan bool
}

func initZipkinSystem(actorSystemName string, options ZipkinOptions) error {
	c, err := zipkin.NewHTTPCollector(options.Url)

	collector = c

	if err != nil {
		ErrorLogger.Printf("Failed to create a Zipkin collector: %v", err)
		return err
	}

	recorder := zipkin.NewRecorder(collector, options.Debug, options.HostPort, actorSystemName)

	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(options.SameSpan),
		zipkin.TraceID128Bit(true),
	)

	if err != nil {
		ErrorLogger.Printf("Failed to create a Zipkin tracer: %v", err)
		return err
	}

	opentracing.InitGlobalTracer(tracer)

	//client := svc1.NewHTTPClient(tracer, svc1Endpoint)

	span = opentracing.StartSpan("Run")
	ctx = opentracing.ContextWithSpan(context.Background(), span)

	span.SetOperationName("System")

	initialized = true

	return nil
}

func logZipkinEvent(msg string) {
	if initialized {
		span.LogEvent(msg)
	}
}

func FinishSpan() {
	if initialized {
		span.Finish()
		collector.Close()
	}
}

func closeZipkinSystem() {
	//if initialized {
	//	span.Finish()
	//	collector.Close()
	//}
}
