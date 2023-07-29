package client

import (
	"context"
	"time"
)

// Client is the interface that wraps the basic method of a client
type Client interface {
	// Call invokes the named function, waits for it to complete,
	// and returns its error status.
	// The args are the arguments for the function call.
	Inovke(ctx context.Context, req, rsp interface{}, path string, opts ...Option) error
}

// defaultClient is the default implementation of the Client interface
type defaultClient struct {
	opts *Options
}

// Options defines the parameters for the Client
type Options struct {
	serName    string        // service name
	method     string        // method name
	target     string        // target address
	timeout    time.Duration // timeout
	network    string        // network type
	protocol   string        // protocol type
	serialType string        // serialization type

	// todo
	// transportOpts transport.ClientOption // transport options
	// interceptors  []interceptor.ClientInterceptor // interceptors
	// selectorName  string // selector name
}

// Option is used to initialize a Option instance
type Option func(*Options)

// WithServiceName is used to set the service name
func WithServiceName(serName string) Option {
	return func(o *Options) {
		o.serName = serName
	}
}

// WithMethodName is used to set the method name
func WithMethodName(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

// WithTimeout is used to set the timeout
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

// WithNetwork is used to set the network type
func WithNetwork(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

// WithProtocol is used to set the protocol type
func WithProtocol(protocol string) Option {
	return func(o *Options) {
		o.protocol = protocol
	}
}

// WithSerialType is used to set the serialization type
func WithSerialType(serialType string) Option {
	return func(o *Options) {
		o.serialType = serialType
	}
}

// singleton pattern to create a defaultClient instance
var defaultClientInstance = NewDefaultClient()

var NewDefaultClient = func() *defaultClient {
	return &defaultClient{
		opts: &Options{
			serName:    "default",
			method:     "default",
			timeout:    time.Second * 5,
			network:    "tcp",
			protocol:   "json",
			serialType: "json",
		},
	}
}

// the first method, through reflection, server dynamically
// generates the corresponding method, and then calls the method
// so append a Call method to be called by the upstream users
// Use the idea of proxy pattern
// Call is used to invoke the named function
func (c *defaultClient) Call(ctx context.Context, req interface{}, rsp interface{}, path string, opts ...Option) error {
	// reflection calls need to be serialized
	callOpts := make([]Option, 0, len(opts)+1)
	callOpts = append(callOpts, opts...)
	//callOpts = append(callOpts, WithSerialType(codec.MsgPack))

	err := c.Inovke(ctx, req, rsp, path, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

// the second method, proto help the server to generate the
// corresponding method, and generate some stub code to be called
// by the client, so it was not necessary to append a Call method

// Inovke is used to complete the full process of calling the remote method
func (c *defaultClient) Inovke(ctx context.Context, req interface{}, rsp interface{}, path string, opts ...Option) error {
	// serialize the request
	// serialization := codec.GetSerialization(c.opts.serialType)
	// payLoad, err := serialization.Marshal(req)
	// if err != nil {
	// 	return errors.Wrap(err, "client: failed to marshal request")
	// }

	// package the request
	// request := addReqHeader(ctx, payLoad)
	// reqbuf, err := proto.Marshal(request)
	// if err != nil {
	// 	return errors.Wrap(err, "client: failed to marshal request")
	// }

	// through transport to send the request to the server, and get the response
	// transport := c.NewTransport()
	// transportOpts := []transport.ClientTransportOption{
	// 	transport.WithServiceName(c.opts.serName),
	// 	transport.WithTarget(c.opts.target),
	// 	transport.WithNetwork(c.opts.network),
	//  transport.WithClientPool(connPoll.GetPool("default")),
	//  transport.WithSelector(selector.GetSelector(c.opts.selectorName)),
	// 	transport.WithTimeout(c.opts.timeout),
	//  }
	// frame, err := transport.Send(ctx, reqbuf, transportOpts...)

	// parse the response
	// rspbuf, err := clientCodec.Decode(frame)
	// if err != nil {
	// 	return errors.Wrap(err, "client: failed to decode response")
	// }

	// parse the response header
	// response := &protocol.Response{}
	// if err := proto.Unmarshal(rspbuf, response); err != nil {
	// 	return errors.Wrap(err, "client: failed to unmarshal response")
	// }

	// unmashal the response body
	// serialization.Unmarshal(response.Payload, rsp)
	return nil
}
