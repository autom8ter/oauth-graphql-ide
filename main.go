package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/autom8ter/oauth-graphql-playground/internal/logger"
	"github.com/autom8ter/oauth-graphql-playground/internal/playground"
	"github.com/autom8ter/oauth-graphql-playground/internal/session"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func init() {
	godotenv.Load()
}

func main() {
	var (
		clientID         = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CLIENT_ID")
		clientSecret     = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CLIENT_SECRET")
		authurl          = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_AUTHORIZATION_URL")
		tokenUrl         = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_TOKEN_URL")
		redirectUrl      = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_REDIRECT_URL")
		scopes           = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SCOPES")
		sessionManager   = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SESSION_MANAGER")
		port             = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_PORT")
		useIDToken       = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_OPEN_ID") == "true"
		graphqlServer    = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SERVER_ENDPOINT")
		debug            = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_DEBUG") == "true"
		corsAllowOrigins = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CORS_ALLOW_ORIGINS")
		corsAllowMethods = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CORS_ALLOW_METHODS")
		corsAllowHeaders = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CORS_ALLOW_HEADERS")
		tlsCertFile      = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_TLS_CERT_FILE")
		tlsKeyFile       = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_TLS_KEY_FILE")
	)
	lgger := logger.New(debug)
	if port == "" {
		port = "5000"
	}
	if graphqlServer == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_SERVER_ENDPOINT")
		return
	}
	if clientID == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_CLIENT_ID")
		return
	}
	if clientSecret == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_CLIENT_SECRET")
		return
	}
	if authurl == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_AUTHORIZATION_URL")
		return
	}
	if tokenUrl == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_TOKEN_URL")
		return
	}
	if redirectUrl == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_REDIRECT_URL")
		return
	}
	if scopes == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_SCOPES")
		return
	}
	if sessionManager == "" {
		lgger.Error("validation error - empty OAUTH_GRAPHQL_PLAYGROUND_SESSION_MANAGER")
		return
	}
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authurl,
			TokenURL: tokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      strings.Split(scopes, ","),
	}
	sessionCfg := map[string]string{}
	if err := json.Unmarshal([]byte(sessionManager), &sessionCfg); err != nil {
		lgger.Error("validation error - failed to decode JSON string - OAUTH_GRAPHQL_PLAYGROUND_SESSION_MANAGER")
		return
	}
	manager, err := session.GetSessionManager(config, sessionCfg)
	if err != nil {
		lgger.Error("failed to setup session manager", zap.Error(err))
		return
	}
	graphqlServerURL, err := url.Parse(graphqlServer)
	if err != nil {
		lgger.Error("failed to parse graphql server URL", zap.Error(err))
		return
	}
	p := playground.NewPlayground(manager, lgger, "/", useIDToken, graphqlServerURL)
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.Playground())

	mux.Handle("/proxy", http.StripPrefix("/proxy", p.Proxy()))
	mux.HandleFunc("/oauth2/callback", p.OAuthCallback())
	lgger.Debug("starting server",
		zap.String("port", port),
		zap.String("playground_path", "/"),
		zap.String("proxy_path", "/proxy"),
		zap.String("oauth_callback_path", "/oauth2/callback"),
	)
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(corsAllowOrigins, ","),
		AllowedMethods:   strings.Split(corsAllowMethods, ","),
		AllowedHeaders:   strings.Split(corsAllowHeaders, ","),
		AllowCredentials: true,
	})
	server := http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		Handler:   c.Handler(mux),
		TLSConfig: nil,
	}
	if tlsKeyFile != "" && tlsCertFile != "" {
		cer, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			lgger.Error("failed to load tls config", zap.Error(err))
			return
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cer},
		}
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		lgger.Error("failed to setup tcp listener", zap.Error(err))
		return
	}
	defer lis.Close()
	if err := server.Serve(lis); err != nil {
		lgger.Error("server failure", zap.Error(err))
	}
}
