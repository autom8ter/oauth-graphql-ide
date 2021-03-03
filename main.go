package main

import (
	"encoding/json"
	"fmt"
	"github.com/autom8ter/oauth-graphql-playground/internal/logger"
	"github.com/autom8ter/oauth-graphql-playground/internal/playground"
	"github.com/autom8ter/oauth-graphql-playground/internal/session"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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
		clientID       = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CLIENT_ID")
		clientSecret   = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_CLIENT_SECRET")
		authurl        = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_AUTHORIZATION_URL")
		tokenUrl       = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_TOKEN_URL")
		redirectUrl    = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_REDIRECT_URL")
		scopes         = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SCOPES")
		sessionManager = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SESSION_MANAGER")
		port           = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_PORT")
		useIDToken     = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_OPEN_ID") == "true"
		graphqlServer  = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_SERVER_ENDPOINT")
		debug          = os.Getenv("OAUTH_GRAPHQL_PLAYGROUND_DEBUG") == "true"
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
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		lgger.Error("server failure", zap.Error(err))
	}
}
