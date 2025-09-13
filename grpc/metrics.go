package grpc

import (
	"context"
	"strings"
	"sync"
	"time"

	coreMetrics "github.com/Rasikrr/core/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	once    sync.Once
	metrics *Metrics
)

type Metrics struct {
	reqTotal    coreMetrics.CounterVec
	inflight    coreMetrics.GaugeVec
	latencySec  coreMetrics.HistogramVec
	msgInBytes  coreMetrics.HistogramVec
	msgOutBytes coreMetrics.HistogramVec
	streamRecv  coreMetrics.CounterVec
	streamSend  coreMetrics.CounterVec
}

func initGRPCMetrics() {
	once.Do(func() {
		dur := []float64{0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
		sz := []float64{128, 256, 512, 1 << 10, 2 << 10, 4 << 10, 8 << 10, 16 << 10, 32 << 10, 64 << 10, 128 << 10, 1 << 20, 5 << 20, 10 << 20}

		metrics = &Metrics{
			reqTotal:    coreMetrics.NewCounterVec("grpc", "requests_total", "gRPC requests", []string{"type", "service", "method", "code"}, nil),
			inflight:    coreMetrics.NewGaugeVec("grpc", "inflight", "In-flight RPCs", []string{"type", "service", "method"}, nil),
			latencySec:  coreMetrics.NewHistogramVec("grpc", "request_seconds", "RPC latency", dur, []string{"type", "service", "method", "code"}, nil),
			msgInBytes:  coreMetrics.NewHistogramVec("grpc", "msg_in_bytes", "Inbound message size", sz, []string{"type", "service", "method"}, nil),
			msgOutBytes: coreMetrics.NewHistogramVec("grpc", "msg_out_bytes", "Outbound message size", sz, []string{"type", "service", "method"}, nil),
			streamRecv:  coreMetrics.NewCounterVec("grpc", "stream_recv_total", "Stream messages received", []string{"service", "method"}, nil),
			streamSend:  coreMetrics.NewCounterVec("grpc", "stream_send_total", "Stream messages sent", []string{"service", "method"}, nil),
		}
	})
}

func (m *Metrics) UnaryServer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		svc, method := split(info.FullMethod)
		t := "unary"
		m.inflight.WithLabelValues(t, svc, method).Inc()
		start := time.Now()
		defer func() { m.inflight.WithLabelValues(t, svc, method).Dec() }()

		if pm, ok := req.(proto.Message); ok {
			m.msgInBytes.WithLabelValues(t, svc, method).Observe(float64(proto.Size(pm)))
		}

		resp, err = handler(ctx, req)
		code := status.Code(err).String()
		m.reqTotal.WithLabelValues(t, svc, method, code).Inc()
		m.latencySec.WithLabelValues(t, svc, method, code).Observe(time.Since(start).Seconds())

		if pm, ok := resp.(proto.Message); ok {
			m.msgOutBytes.WithLabelValues(t, svc, method).Observe(float64(proto.Size(pm)))
		}
		return resp, err
	}
}

func (m *Metrics) StreamServer() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		svc, method := split(info.FullMethod)
		t := "stream"
		m.inflight.WithLabelValues(t, svc, method).Inc()
		start := time.Now()
		defer func() { m.inflight.WithLabelValues(t, svc, method).Dec() }()

		w := &serverStream{ServerStream: ss, m: m, svc: svc, method: method}
		err := handler(srv, w)
		code := status.Code(err).String()
		m.reqTotal.WithLabelValues(t, svc, method, code).Inc()
		m.latencySec.WithLabelValues(t, svc, method, code).Observe(time.Since(start).Seconds())
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	m      *Metrics
	svc    string
	method string
}

func (w *serverStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.m.streamRecv.WithLabelValues(w.svc, w.method).Inc()
		if pm, ok := m.(proto.Message); ok {
			w.m.msgInBytes.WithLabelValues("stream", w.svc, w.method).Observe(float64(proto.Size(pm)))
		}
	}
	return err
}

func (w *serverStream) SendMsg(m interface{}) error {
	if pm, ok := m.(proto.Message); ok {
		w.m.msgOutBytes.WithLabelValues("stream", w.svc, w.method).Observe(float64(proto.Size(pm)))
	}
	w.m.streamSend.WithLabelValues(w.svc, w.method).Inc()
	return w.ServerStream.SendMsg(m)
}

func (m *Metrics) UnaryClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		svc, meth := split(method)
		t := "unary"
		m.inflight.WithLabelValues(t, svc, meth).Inc()
		start := time.Now()
		defer func() { m.inflight.WithLabelValues(t, svc, meth).Dec() }()

		if pm, ok := req.(proto.Message); ok {
			m.msgOutBytes.WithLabelValues(t, svc, meth).Observe(float64(proto.Size(pm)))
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		code := status.Code(err).String()
		m.reqTotal.WithLabelValues(t, svc, meth, code).Inc()
		m.latencySec.WithLabelValues(t, svc, meth, code).Observe(time.Since(start).Seconds())

		if pm, ok := reply.(proto.Message); ok {
			m.msgInBytes.WithLabelValues(t, svc, meth).Observe(float64(proto.Size(pm)))
		}
		return err
	}
}

func (m *Metrics) StreamClient() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		svc, meth := split(method)
		t := "stream"
		m.inflight.WithLabelValues(t, svc, meth).Inc()
		start := time.Now()

		cs, err := streamer(ctx, desc, cc, method, opts...)
		code := status.Code(err).String()
		if err != nil {
			m.reqTotal.WithLabelValues(t, svc, meth, code).Inc()
			m.latencySec.WithLabelValues(t, svc, meth, code).Observe(time.Since(start).Seconds())
			m.inflight.WithLabelValues(t, svc, meth).Dec()
			return cs, err
		}

		w := &clientStream{ClientStream: cs, m: m, svc: svc, method: meth, started: start}
		return w, nil
	}
}

type clientStream struct {
	grpc.ClientStream
	m       *Metrics
	svc     string
	method  string
	started time.Time
}

func (w *clientStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if err == nil {
		w.m.streamRecv.WithLabelValues(w.svc, w.method).Inc()
		if pm, ok := m.(proto.Message); ok {
			w.m.msgInBytes.WithLabelValues("stream", w.svc, w.method).Observe(float64(proto.Size(pm)))
		}
	}
	return err
}

func (w *clientStream) SendMsg(m interface{}) error {
	if pm, ok := m.(proto.Message); ok {
		w.m.msgOutBytes.WithLabelValues("stream", w.svc, w.method).Observe(float64(proto.Size(pm)))
	}
	w.m.streamSend.WithLabelValues(w.svc, w.method).Inc()
	return w.ClientStream.SendMsg(m)
}

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()
	code := status.Code(err).String()
	w.m.reqTotal.WithLabelValues("stream", w.svc, w.method, code).Inc()
	w.m.latencySec.WithLabelValues("stream", w.svc, w.method, code).Observe(time.Since(w.started).Seconds())
	w.m.inflight.WithLabelValues("stream", w.svc, w.method).Dec()
	return err
}

func split(full string) (service, method string) {
	// "/pkg.Service/Method"
	full = strings.TrimPrefix(full, "/")
	i := strings.LastIndex(full, "/")
	if i < 0 {
		return "unknown", full
	}
	return full[:i], full[i+1:]
}
