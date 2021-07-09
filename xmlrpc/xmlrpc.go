package xmlrpc

import (
	"encoding/xml"
	"io"
	"net/rpc"
	"sync"
)

type clientResponse struct {
	Raw []byte `xml:",innerxml"`
}

func (r *clientResponse) reset() {
	r.Raw = nil
}

type clientRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     [1]struct {
		Text  string `xml:",chardata"`
		Param p      `xml:"param"`
	} `xml:"params"`
}

type p struct {
	Text  string `xml:",chardata"`
	Value struct {
		Text string `xml:",chardata"`
	} `xml:"value"`
}

type clientCodec struct {
	dec *xml.Decoder
	enc *xml.Encoder
	c   io.Closer

	req  clientRequest
	resp clientResponse

	mu      sync.Mutex
	pending map[uint64]string
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	c.req.MethodName = r.ServiceMethod
	par := p{
		Text: "",
		Value: struct {
			Text string "xml:\",chardata\""
		}{"lmao"},
	}
	c.req.Params[0] = struct {
		Text  string "xml:\",chardata\""
		Param p      "xml:\"param\""
	}{"", par}
	return c.enc.Encode(&c.req)
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil
	}
	return xml.Unmarshal(c.resp.Raw, x)
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.resp.reset()
	if err := c.dec.Decode(&c.resp); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		dec: xml.NewDecoder(conn),
		enc: xml.NewEncoder(conn),
		c:   conn,
	}
}
