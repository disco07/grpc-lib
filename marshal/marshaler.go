package marshal

import (
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
)

// WithMultipartFormMarshaler returns a ServeMuxOption which associates inbound and outbound Marshalers to a MIME type in mux.
func WithMultipartFormMarshaler() runtime.ServeMuxOption {
	return runtime.WithMarshalerOption("multipart/form-data", &multipartFormMarshaler{
		HTTPBodyMarshaler: &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions:   protojson.MarshalOptions{EmitUnpopulated: true},
				UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
			},
		},
	})
}

type multipartFormMarshaler struct {
	*runtime.HTTPBodyMarshaler
}

func (h *multipartFormMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &multipartFormDecoder{
		Decoder: h.HTTPBodyMarshaler.NewDecoder(r),
		body:    r,
		buf:     make([]byte, 0, 512),
		eof:     false,
	}
}

type multipartFormDecoder struct {
	runtime.Decoder

	body io.Reader
	buf  []byte
	eof  bool
}

func (d *multipartFormDecoder) Decode(v interface{}) error {
	body, ok := v.(*httpbody.HttpBody)
	if !ok {
		return d.Decoder.Decode(v)
	}

	if d.eof {
		return io.EOF
	}

	n, err := io.ReadFull(d.body, d.buf)
	if n > 0 {
		body.Data = d.buf[:n]
	}

	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		d.eof = true
		return nil
	}

	return err
}
