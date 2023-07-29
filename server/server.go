package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server apply the ability to register services and start the server
type Server struct {
	opts     *ServerOptions     // server options
	services map[string]Service // registered services
}

// Service represents a service that can be registered and started
type Service interface {
	// Register(string, HandlerFunc) error
	Serve(*ServerOptions) error
	Close() error
}

// service is the default implementation of the Service interface
type service struct {
	svr         interface{}        // server instance
	ctx         context.Context    // every service has a context to control its lifecycle
	cancel      context.CancelFunc // cancel function
	serviceName string             // service name
	// handlers    map[string]HandlerFunc // service handlers
	opts *ServerOptions // server options
}

// ServerOptions is used to initialize the server
// by passing in the configuration
type ServerOptions struct {
	addr       string        // listening address, e.g. :( ip://127.0.0.1:8080)
	network    string        // network types, e.g. : (tcp, udp)
	protocol   string        // protocol types, e.g. : (json, proto)
	timeout    time.Duration // timeout, e.g. : (1s, 500ms)
	serialType string        // serialization types, e.g. : (json, proto)

	// todo
	// selectorSvrAddr string // selector address, requeired when using the selector
	// tracSvrAddr string // trace address, required when using the trace
	// tracSpanName string // trace span name, required when using the trace
	// pluginNames []string // plugin names, required when using the plugin
	// intercrptors []interceptor.ServerInterceptor // interceptors, required when using the interceptor
}

// ServerOption is used to initialize a ServerOption instance
// this is the key to functional option pattern
type ServerOption func(*ServerOptions)

// WithAddr is used to set the listening address
func WithAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.addr = addr
	}
}

// WithNetwork is used to set the network type
func WithNetwork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.network = network
	}
}

// WithProtocol is used to set the protocol type
func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.protocol = protocol
	}
}

// WithTimeout is used to set the timeout
func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

// WithSerialType is used to set the serialization type
func WithSerialType(serialType string) ServerOption {
	return func(o *ServerOptions) {
		o.serialType = serialType
	}
}

// // WithSelectorSvrAddr is used to set the selector address
// func WithSelectorSvrAddr(selectorSvrAddr string) ServerOption {
// 	return func(o *ServerOptions) {
// 		// o.selectorSvrAddr = selectorSvrAddr
// 	}
// }

// // WithTracSvrAddr is used to set the trace address
// func WithTracSvrAddr(tracSvrAddr string) ServerOption {
// 	return func(o *ServerOptions) {
// 		// o.tracSvrAddr = tracSvrAddr
// 	}
// }

// // WithTracSpanName is used to set the trace span name
// func WithTracSpanName(tracSpanName string) ServerOption {
// 	return func(o *ServerOptions) {
// 		// o.tracSpanName = tracSpanName
// 	}
// }

// // WithPluginNames is used to set the plugin names
// func WithPluginNames(pluginNames []string) ServerOption {
// 	return func(o *ServerOptions) {
// 		// o.pluginNames = pluginNames
// 	}
// }

// // WithIntercrptors is used to set the interceptors
// func WithIntercrptors(intercrptors []interceptor.ServerInterceptor) ServerOption {
// 	return func(o *ServerOptions) {
// 		// o.intercrptors = intercrptors
// 	}
// }

// NewServer creates a new server instance
func NewServer(opt ...ServerOption) *Server {
	s := &Server{
		opts:     &ServerOptions{},
		services: make(map[string]Service),
	}
	for _, o := range opt {
		o(s.opts)
	}
	return s
}

// Serve is the key method to start the server, which will start all registered services
func (s *Server) Serve() {
	for _, service := range s.services {
		go service.Serve(s.opts)
	}

	// this ch is used to receive the signal from the system want to stop the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch

	s.Close()
}

// Close is used to close all registered services
func (s *Server) Close() {
	for _, service := range s.services {
		service.Close()
	}
}
