package gosiris

import (
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/opentracing/opentracing-go"
	zlog "github.com/opentracing/opentracing-go/log"
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

	return nil
}

func logZipkinFields(span opentracing.Span, fields ...zlog.Field) {
	span.LogFields(fields...)
}

func logZipkinMessage(span opentracing.Span) {
}

//
//func logZipkinKV(spanName string, alternatingKeyValues ...interface{}) {
//	if !zipkinSystemInitialized {
//		ErrorLogger.Printf("Zipkin system not started")
//		return
//	}
//
//	span, exists := spans[spanName]
//	if !exists {
//		ErrorLogger.Printf("Span %v not started", spanName)
//		return
//	}
//
//	span.LogKV(alternatingKeyValues...)
//}

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

//func startZipkinChildSpan(parentSpanName, spanName, operationName string) {
//	if !zipkinSystemInitialized {
//		ErrorLogger.Printf("Zipkin system not started")
//		return
//	}
//
//	parentSpan, exists := spans[parentSpanName]
//	if !exists {
//		ErrorLogger.Printf("Parent span %v not started", parentSpanName)
//		return
//	}
//
//	span := opentracing.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
//	spans[spanName] = span
//
//	InfoLogger.Printf("Child span %v of parent %v started", spanName, parentSpanName)
//}
//
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
