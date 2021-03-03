package playground

import (
	"fmt"
	"github.com/autom8ter/oauth-graphql-playground/internal/logger"
	session2 "github.com/autom8ter/oauth-graphql-playground/internal/session"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Playground struct {
	session        session2.SessionManager
	logger         *logger.Logger
	playgroundPath string
	useIDToken     bool
	proxy          *httputil.ReverseProxy
}

func NewPlayground(session session2.SessionManager, logger *logger.Logger, playgroundPath string, useIDToken bool, endpoint *url.URL) *Playground {
	director := func(req *http.Request) {
		req.URL.Scheme = endpoint.Scheme
		req.URL.Host = endpoint.Host
		req.URL.Path = endpoint.Path
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	p := &Playground{session: session, logger: logger, playgroundPath: playgroundPath, useIDToken: useIDToken, proxy: &httputil.ReverseProxy{}}
	p.proxy.Director = func(req *http.Request) {
		fmt.Println(req.URL.String())
		authToken, err := p.session.GetToken(req)
		if err != nil {
			director(req)
			return
		}
		if authToken == nil {
			director(req)
			return
		}
		if !authToken.Token.Valid() {
			director(req)
			return
		}
		token, err := p.session.GetToken(req)
		if err != nil {
			director(req)
			return
		}
		if p.useIDToken {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.IDToken))
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token.AccessToken))
		}
		director(req)
		logger.Debug("proxying graphQL request", zap.String("url", req.URL.String()))
	}
	return p
}

func (p *Playground) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p.proxy.ServeHTTP(w, req)
	}
}

func (p *Playground) Playground() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authToken, err := p.session.GetToken(req)
		if err != nil {
			p.logger.Error("playground: failed to get session - redirecting", zap.Error(err))
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		if authToken == nil {
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		if !authToken.Token.Valid() {
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		token, err := p.session.GetToken(req)
		if err != nil {
			p.logger.Error("playground: failed to get token - redirecting", zap.Error(err))
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		w.Header().Add("Content-Type", "text/html")
		if p.useIDToken {
			playground.Execute(w, map[string]string{
				"token": token.IDToken,
			})
		} else {
			playground.Execute(w, map[string]string{
				"token": token.Token.AccessToken,
			})
		}

	}
}

func (p *Playground) OAuthCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		state := req.URL.Query().Get("state")
		if code == "" {
			p.logger.Error("playground: empty authorization code - redirecting")
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		if state == "" {
			p.logger.Error("playground: empty authorization state - redirecting")
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}

		stateVal, err := p.session.GetState(req)
		if err != nil {
			p.logger.Error("playground: failed to get session state - redirecting", zap.Error(err))
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		if stateVal != state {
			p.logger.Error("playground: session state mismatch - redirecting")
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}

		if err := p.session.Exchange(w, req, code); err != nil {
			p.logger.Error("playground: failed to exchange authorization code - redirecting", zap.Error(err))
			if err := p.session.RedirectLogin(w, req); err != nil {
				p.logger.Error("playground: failed to redirect", zap.Error(err))
			}
			return
		}
		http.Redirect(w, req, p.playgroundPath, http.StatusTemporaryRedirect)
	}
}