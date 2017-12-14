package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/oauth2"

	"github.com/anxiousmodernman/goth"
	"github.com/anxiousmodernman/goth/gothic"
	"github.com/anxiousmodernman/goth/providers/auth0"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/johanbrandhorst/protobuf/wsproxy"
	"github.com/sirupsen/logrus"

	auth0z "github.com/auth0-community/go-auth0"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"gitlab.com/DSASanFrancisco/co-chair/backend"
	"gitlab.com/DSASanFrancisco/co-chair/frontend/bundle"
	"gitlab.com/DSASanFrancisco/co-chair/proto/server"
)

var (
	// Store is our sessions store.
	Store *sessions.CookieStore
)

// TODO pass this down to my object
var logger *logrus.Logger

func init() {
	Store = sessions.NewCookieStore([]byte("something-very-secret"))
	logger = logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})
	// Should only be done from init functions
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))

	store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))
	store.MaxLength(math.MaxInt64)
	gothic.Store = store
}

func main() {

	goth.UseProviders(
		auth0.New(
			os.Getenv("COCHAIR_AUTH0_CLIENTID"),
			os.Getenv("COCHAIR_AUTH0_SECRET"), "https://localhost:2016/auth/auth0/callback",
			os.Getenv("COCHAIR_AUTH0_DOMAIN")),
	)

	// prevent error:
	// gob: type not registered for interface: map[string]interface {}
	var t map[string]interface{}
	gob.Register(t)

	// NewProxy gives us a Proxy, our concrete implementation of the
	// interface generated by the grpc protobuf compiler.
	px, err := backend.NewProxy("co-chair.db")
	if err != nil {
		log.Fatalf("proxy init: %v", err)
	}

	gs := grpc.NewServer()
	server.RegisterProxyServer(gs, px)
	wrappedServer := grpcweb.WrapServer(gs)

	clientCreds, err := credentials.NewClientTLSFromFile("./cert.pem", "")
	if err != nil {
		logger.WithError(err).Fatal("Failed to get local server client credentials, did you run `make generate_cert`?")
	}

	wsproxy := wsproxy.WrapServer(
		http.HandlerFunc(wrappedServer.ServeHTTP),
		wsproxy.WithLogger(logger),
		wsproxy.WithTransportCredentials(clientCreds))

	// Note: routes are evaluated in the order they're defined.
	p := mux.NewRouter()

	p.Handle("/login", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.Wrap(http.HandlerFunc(loginLink)),
	)).Methods("GET")

	p.Handle("/auth/{provider}/callback", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.Wrap(http.HandlerFunc(oauthCallbackHandler)),
	)).Methods("GET")

	p.Handle("/auth/{provider}", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.Wrap(http.HandlerFunc(loginHandler)),
	)).Methods("GET")

	p.Handle("/logout/{provider}", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.Wrap(http.HandlerFunc(logoutHandler)),
	)).Methods("GET")

	p.Handle("/", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.HandlerFunc(IsAuthenticated),
		negroni.Wrap(websocketsProxy(wsproxy)),
	)).Methods("POST")

	p.Handle("/", negroni.New(
		negroni.HandlerFunc(withLog),
		negroni.HandlerFunc(IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(homeHandler)),
	)).Methods("GET")

	addr := "localhost:2016"
	httpsSrv := &http.Server{
		Addr:    addr,
		Handler: p,
		// Some security settings
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		},
	}

	logger.Info("Serving on https://" + addr)
	err = httpsSrv.ListenAndServeTLS("./cert.pem", "./key.pem")
	logger.Fatalf("handler exit: %s", err.Error())
}

var indexTemplate = `
<p><a href="/auth/auth0">Log in with auth0</a></p>
`

func extractor(r *http.Request) (string, error) {

	if s, err := jwtmiddleware.FromParameter("code")(r); err == nil {
		if s != "" {
			logger.Debug("extractor used param code")
			return s, nil
		}
	}

	fmt.Println("dumping cookies")
	for _, c := range r.Cookies() {
		fmt.Println("cookie:", c.Value)
	}

	if c, err := r.Cookie("auth0_gothic_session"); err == nil {
		logger.Debug("extractor used cookie")
		return c.Value, nil
	}

	return "", errors.New("no authorization found")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	domain := "dsasf.auth0.com"
	aud := ""

	conf := &oauth2.Config{
		ClientID:     os.Getenv("COCHAIR_AUTH0_CLIENTID"),
		ClientSecret: os.Getenv("COCHAIR_AUTH0_SECRET"),
		RedirectURL:  "https://localhost:2016/auth/auth0/callback",
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	if aud == "" {
		aud = "https://" + domain + "/userinfo"
	}

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := Store.Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	// add "code" here?
	url := conf.AuthCodeURL(state, audience)
	logger.Debug("auth code url: ", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func loginLink(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := os.Getenv("COCHAIR_AUTH0_DOMAIN")

	var u *url.URL
	u, err := url.Parse("https://" + domain)

	if err != nil {
		panic("boom")
	}

	u.Path += "/auth/auth0/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", "https://localhost:2016")
	parameters.Add("client_id", os.Getenv("COCHAIR_AUTH0_CLIENTID"))
	u.RawQuery = parameters.Encode()

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	domain := "dsasf.auth0.com"

	conf := &oauth2.Config{
		ClientID:     os.Getenv("COCHAIR_AUTH0_CLIENTID"),
		ClientSecret: os.Getenv("COCHAIR_AUTH0_SECRET"),
		RedirectURL:  "https://localhost:2016/auth/auth0/callback",
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}
	// Validate state before calling Exchange
	state := r.URL.Query().Get("state")
	session, err := Store.Get(r, "state")
	if err != nil {
		logger.Errorf("could not get session state: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if state != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	var code string
	code = r.URL.Query().Get("code")
	if code == "" {
		logger.Error("code is not in query params")
	}
	code = r.FormValue("code")
	if code == "" {
		logger.Error("code is not in form")
	}
	idToken := r.URL.Query().Get("id_token")
	if idToken != "" {
		logger.Debugf("got a token: %s", idToken)
	}

	// package oauth2 docs:
	// The code will be in the *http.Request.FormValue("code").
	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		logger.Error("oauth exchange failure: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		logger.Errorf("error calling userinfo endpoint: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		logger.Errorf("could not decode userinfo response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err = Store.Get(r, "auth-session")
	if err != nil {
		logger.Errorf("could not get session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		logger.Errorf("could not save session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(bundle.Assets).ServeHTTP(w, r)
}

func websocketsProxy(wsproxy http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") ||
			websocket.IsWebSocketUpgrade(r) {
			wsproxy.ServeHTTP(w, r)
		}
	}
}

func checkJWT(h http.HandlerFunc) http.HandlerFunc {

	type Response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		JWKS_URI := "https://" + os.Getenv("COCHAIR_AUTH0_DOMAIN") + "/.well-known/jwks.json"
		client := auth0z.NewJWKClient(auth0z.JWKClientOptions{URI: JWKS_URI})
		aud := "dsasf." //os.Getenv("AUTH0_AUDIENCE")
		audience := []string{aud}

		var AUTH0_API_ISSUER = "https://" + os.Getenv("COCHAIR_AUTH0_DOMAIN") + "/"
		configuration := auth0z.NewConfiguration(client, audience, AUTH0_API_ISSUER, jose.RS256)
		validator := auth0z.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println("Token is not valid or missing token")

			response := Response{
				Message: "Missing or invalid token.",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)

		} else {
			// Ensure the token has the correct scope
			result := checkScope(r, validator, token)
			if result == true {
				// If the token is valid and we have the right scope, we'll pass through the middleware
				h.ServeHTTP(w, r)
			} else {
				response := Response{
					Message: "You do not have the read:messages scope.",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)

			}
		}
	}
}

func checkScope(r *http.Request, validator *auth0z.JWTValidator, token *jwt.JSONWebToken) bool {
	claims := map[string]interface{}{}
	err := validator.Claims(r, token, &claims)

	if err != nil {
		fmt.Println(err)
		return false
	}

	logger.Debug("scopes:", claims["scope"])

	if claims["scope"] != nil && strings.Contains(claims["scope"].(string), "read:messages") {
		return true
	} else {
		// TODO: insecure on purpose here...
		return true
	}
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		logger.Errorf("session profile not found; redirecting.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		next(w, r)
	}
}

func withLog(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	logger.Debug("path:", r.URL.Path)
	next(w, r)
}
