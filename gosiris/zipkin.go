package gosiris

import (
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	zipkinSystemInitialized bool
	collector               zipkin.Collector
	tracer                  opentracing.Tracer
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

	t, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(options.SameSpan),
		zipkin.TraceID128Bit(true),
	)

	tracer = t

	if err != nil {
		ErrorLogger.Printf("Failed to create a Zipkin tracer: %v", err)
		return err
	}

	opentracing.InitGlobalTracer(tracer)

	zipkinSystemInitialized = true

	InfoLogger.Printf("Zipkin tracer started")

	return nil
}

//func logZipkinFields(span opentracing.Span, fields ...zlog.Field) {
//	span.LogFields(fields...)
//}

func logZipkinMessage(span opentracing.Span, event string) {
	span.LogEvent(event)
}

func inject(span opentracing.Span) (opentracing.TextMapCarrier, error) {
	if span == nil {
		return nil, nil
	}

	carrier := opentracing.TextMapCarrier{}

	err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)

	if err != nil {
		ErrorLogger.Printf("Failed to inject: %v", err)
		return nil, err
	}

	return carrier, nil
}

func extract(carrier opentracing.TextMapCarrier) (opentracing.SpanContext, error) {
	return tracer.Extract(opentracing.TextMap, carrier)
}

func startZipkinSpan(spanName, operationName string) opentracing.Span {
	if !zipkinSystemInitialized {
		ErrorLogger.Printf("Zipkin system not started")
		return nil
	}

	span := opentracing.StartSpan(spanName)
	span.SetOperationName(operationName)

	InfoLogger.Printf("Span %v started", spanName)

	return span
}

func stopZipkinSpan(span opentracing.Span) {
	span.Finish()
}

func closeZipkinSystem() {
	if zipkinSystemInitialized {
		InfoLogger.Printf("Closing Zipkin system")
		collector.Close()
		zipkinSystemInitialized = false
	}
}
