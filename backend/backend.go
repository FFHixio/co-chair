package backend

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/Rudd-O/curvetls"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/helpers/derhelpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"

	"github.com/anxiousmodernman/co-chair/proto/server"
	"github.com/asdine/storm"
)

// Proxy is our server.ProxyServer implementation.
type Proxy struct {
	DB  *storm.DB
	mtx *sync.Mutex
}

// NewProxy is our constructor for the server.ProxyServer implementation.
func NewProxy(path string) (*Proxy, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	// find our KeyPair for protecting external grpc conns, and
	// create a keypair if it doesn't exist.
	var kp KeyPair
	err = db.One("Name", "server", &kp)
	if err != nil {
		if err == storm.ErrNotFound {
			priv, pub, err := curvetls.GenKeyPair()
			if err != nil {
				return nil, err
			}
			newKP := KeyPair{Name: "server", Pub: pub.String(), Priv: priv.String()}
			if err := db.Save(&newKP); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &Proxy{db, &sync.Mutex{}}, nil
}

// assert that Proxy is a server.ProxyServer at compile time.
var _ server.ProxyServer = (*Proxy)(nil)

// State returns the state of the proxy. The number of backends returned is
// controlled by the domain field of the request. A blank domain returns all.
func (p *Proxy) State(_ context.Context, req *server.StateRequest) (*server.ProxyState, error) {
	var resp server.ProxyState
	var backends []*BackendData
	var err error
	if req.Domain == "" {
		err = p.DB.All(&backends)
	} else {
		err = p.DB.Find("Domain", req.Domain, &backends)
	}
	if err != nil {
		return nil, fmt.Errorf("domain: %s; db error: %v", req.Domain, err)
	}
	for _, b := range backends {
		// do not leak private keys here
		resp.Backends = append(resp.Backends, b.AsBackend())
	}

	return &resp, nil
}

// Put adds a backend to our pool of proxied Backends.
func (p *Proxy) Put(ctx context.Context, b *server.Backend) (*server.OpResult, error) {

	var bd BackendData
	err := p.DB.One("Domain", b.Domain, &bd)
	// ignore ErrNotFound: always overwrite the BackendData
	if err != nil && err != storm.ErrNotFound {
		return &server.OpResult{}, fmt.Errorf("domain lookup: %v", err)
	}

	bd.Domain = b.Domain
	bd.IPs = combine(bd.IPs, b.Ips)
	bd.Protocol = b.Protocol
	bd.MatchHeaders = b.MatchHeaders

	if b.BackendCert != nil {
		bd.BackendCert = b.BackendCert.Cert
		bd.BackendKey = b.BackendCert.Key
	}

	// Possible feature: generate BackendCerts if blank?
	if false {
		// generate cert
		c := csr.New()
		rootCACert, csrPEM, rootCAKey, err := initca.New(c)
		if err != nil {
			return nil, err
		}
		_ = csrPEM
		crt, err := tls.X509KeyPair(rootCAKey, rootCACert)
		if err != nil {
			return nil, err
		}

		var password string // blank for now
		derredUp, err := helpers.GetKeyDERFromPEM(rootCAKey, []byte(password))
		if err != nil {
			return nil, err
		}
		priv, err := derhelpers.ParsePrivateKeyDER(derredUp)
		if err != nil {
			return nil, err
		}
		signr, err := local.NewSigner(priv, crt.Leaf, x509.ECDSAWithSHA512, nil)
		if err != nil {
			return nil, err
		}
		// our actual server cert and key?
		req := csr.CertificateRequest{KeyRequest: csr.NewBasicKeyRequest()}
		var key, csrBytes []byte
		g := &csr.Generator{Validator: genkey.Validator}
		csrBytes, key, err = g.ProcessRequest(&req)

		var signReq signer.SignRequest
		signReq.Hosts = []string{b.Domain}
		signReq.Request = string(csrBytes)
		newCert, err := signr.Sign(signReq)
		if err != nil {
			return nil, err
		}
		_, _ = newCert, key
	}

	err = p.DB.Save(&bd)
	if err != nil {
		return nil, fmt.Errorf("save: %v", err)
	}

	resp := &server.OpResult{Code: 200, Status: "Ok"}

	return resp, nil
}

// PutKVStream lets us stream key-value pairs into our db.
func (p *Proxy) PutKVStream(stream server.Proxy_PutKVStreamServer) error {

	for {
		kv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := p.DB.SetBytes("streams", kv.Key, kv.Value); err != nil {
			return fmt.Errorf("SetBytes: %v", err)
		}
	}
	return nil
}

// GetKVStream scans a keyspace.
func (p *Proxy) GetKVStream(key *server.Key, stream server.Proxy_GetKVStreamServer) error {

	tx, err := p.DB.Bolt.Begin(false)
	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}
	c := tx.Bucket([]byte("streams")).Cursor()

	// RFC3339 example: the 90's decade.
	//min := []byte("1990-01-01T00:00:00Z")
	//max := []byte("2000-01-01T00:00:00Z")

	// Iterate over the 90's.
	for k, v := c.Seek(key.Prefix); k != nil || false; /* could do something besides "false" */ k, v = c.Next() {
		fmt.Printf("%s: %s\n", k, v)
		kv := server.KV{Key: k, Value: v}
		err := stream.Send(&kv)
		if err != nil {
			return fmt.Errorf("send: %v", err)
		}
	}
	return nil
}

func combine(a, b []string) []string {
	// let's pre-allocate enough space
	both := make([]string, 0, len(a)+len(b))
	both = append(both, a...)
	both = append(both, b...)
	sort.Strings(both)
	var val string
	var res []string
	for _, x := range both {
		if strings.TrimSpace(x) == strings.TrimSpace(val) {
			continue
		}
		val = x
		res = append(res, strings.TrimSpace(x))
	}
	return res
}

// Remove ...
func (p *Proxy) Remove(_ context.Context, b *server.Backend) (*server.OpResult, error) {
	// match on domain name exactly
	var bd BackendData
	if err := p.DB.One("Domain", b.Domain, &bd); err != nil {
		return nil, err
	}
	if err := p.DB.DeleteStruct(&bd); err != nil {
		return nil, err
	}

	res := &server.OpResult{Code: 200, Status: fmt.Sprintf("removed: %s", bd.Domain)}
	return res, nil
}

// BackendData is our type for the storm ORM. We can define field-level
// constraints and indexes on struct tags. It is unfortunate that we
// need an intermediary type, but it seems better than going in and
// adding storm struct tags to protobuf-generated code.
//
// See issue: https://github.com/golang/protobuf/issues/52
type BackendData struct {
	ID     int    `storm:"id,increment"`
	Domain string `storm:"unique"`
	IPs    []string
	// An optional endpoint we can call, expecting HTTP 200
	HealthCheck string
	// one of HTTP1, HTTP2, GRPC
	Protocol server.Backend_Protocol
	// Our TLS certs and keys.
	BackendCert, BackendKey []byte
	// Headers to match on during backend selection when we first
	// get a connection.
	MatchHeaders map[string]string
}

// AsBackend is a conversion method to a grpc-sendable type.
func (bd BackendData) AsBackend() *server.Backend {
	var b server.Backend
	b.Domain = bd.Domain
	b.Ips = bd.IPs
	b.Protocol = bd.Protocol
	return &b
}

// KeyPair is a database type that represents curvetls key pairs. A KeyPair
// must be in the database for each pure grpc client that wants to connect.
// Not used for grpc websocket clients. We rely on HTTPS for those.
type KeyPair struct {
	Name string `storm:"unique,id"`
	// Pub and Priv are base64 strings that represent curvetls keys
	// for servers or clients.
	Pub  string
	Priv string
}

// RetrieveServerKeys gets the curvetls public key for our co-chair instance's api.
// Clients will need to know the server's public key.
func RetrieveServerKeys(db *storm.DB) (curvetls.Pubkey, curvetls.Privkey, error) {
	var kp KeyPair
	err := db.One("Name", "server", &kp)
	if err != nil {
		return curvetls.Pubkey{}, curvetls.Privkey{}, err
	}
	priv, err := curvetls.PrivkeyFromString(kp.Priv)
	if err != nil {
		return curvetls.Pubkey{}, curvetls.Privkey{}, err
	}
	pub, err := curvetls.PubkeyFromString(kp.Pub)
	if err != nil {
		return curvetls.Pubkey{}, curvetls.Privkey{}, err
	}
	return pub, priv, nil
}
