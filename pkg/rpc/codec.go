package rpc

import (
	"bufio"
	"io"
	"net/rpc"

	"github.com/vmihailenco/msgpack"
)

// NewMsgpackServerCodec :
func NewMsgpackServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	buf := bufio.NewWriter(conn)
	return &msgpackServerCodec{
		rwc:    conn,
		dec:    msgpack.NewDecoder(conn),
		enc:    msgpack.NewEncoder(buf),
		encBuf: buf,
	}
}

// NewMsgpackClientCodec :
func NewMsgpackClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	encBuf := bufio.NewWriter(conn)
	return &msgpackClientCodec{
		rwc:    conn,
		dec:    msgpack.NewDecoder(conn),
		enc:    msgpack.NewEncoder(encBuf),
		encBuf: encBuf,
	}
}

type msgpackServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *msgpack.Decoder
	enc    *msgpack.Encoder
	encBuf *bufio.Writer
	closed bool
}

func (c *msgpackServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.dec.Decode(r)
}

func (c *msgpackServerCodec) ReadRequestBody(body interface{}) error {
	var err error

	if body == nil {
		body, err = c.dec.DecodeInterface()
		return err
	}

	return c.dec.Decode(body)
}

func (c *msgpackServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {

	if err = c.enc.Encode(r); err != nil {
		if c.encBuf.Flush() == nil {
			// Msgpack couldn't encode the header. Should not happen, so if it does,
			// shut down the connection to signal that the connection is broken.
			c.Close()
		}
		return
	}
	if err = c.enc.Encode(body); err != nil {
		if c.encBuf.Flush() == nil {
			// Was a msgpack problem encoding the body but the header has been written.
			// Shut down the connection to signal that the connection is broken.
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *msgpackServerCodec) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

type msgpackClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *msgpack.Decoder
	enc    *msgpack.Encoder
	encBuf *bufio.Writer
}

func (c *msgpackClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *msgpackClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *msgpackClientCodec) ReadResponseBody(body interface{}) error {

	var err error

	if body == nil {
		body, err = c.dec.DecodeInterface()
		return err
	}

	return c.dec.Decode(body)
}

func (c *msgpackClientCodec) Close() error {
	return c.rwc.Close()
}
