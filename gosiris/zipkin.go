package gosiris

import (
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/opentracing/opentracing-go"
	"context"
)

var (
	ctx       context.Context
	span      opentracing.Span
	collector zipkin.Collector
)

type ZipkinOptions struct {
	url string
}

func zipkinInit(url string, debug bool, hostPort string, actorSystemName string, sameSpan bool) error {
	c, err := zipkin.NewHTTPCollector(url)

	collector = c

	if err != nil {
		ErrorLogger.Printf("Failed to create a Zipkin collector: %v", err)
		return err
	}

	recorder := zipkin.NewRecorder(collector, debug, hostPort, actorSystemName)

	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
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

	return nil
}

func logEvent(msg string) {
	span.LogEvent(msg)

}
