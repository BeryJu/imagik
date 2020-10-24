package auth

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

type OIDCAuth struct {
	context  context.Context
	config   oauth2.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	logger   *log.Entry
}

func (oa *OIDCAuth) handleRedirect(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	state := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	session.Values["oidc_state"] = state
	http.Redirect(w, r, oa.config.AuthCodeURL(state), http.StatusFound)
}

func (oa *OIDCAuth) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	oauth2Token, err := oa.config.Exchange(oa.context, r.URL.Query().Get("code"))
	if err != nil {
		oa.logger.Warn(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		oa.logger.Warn(err)
	}

	// Parse and verify ID Token payload.
	idToken, err := oa.verifier.Verify(oa.context, rawIDToken)
	if err != nil {
		oa.logger.Warn(err)
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		oa.logger.Warn(err)
	}
	session, _ := store.Get(r, "session-name")
	session.Values["test"] = "test"
}

func (oa *OIDCAuth) Init() {
	oa.logger = log.WithField("component", "OIDC-auth")
	oa.context = context.Background()
	provider, err := oidc.NewProvider(oa.context, config.C.AuthOIDCConfig.URL)
	if err != nil {
		oa.logger.Warn(err)
	}
	oa.provider = provider
	oa.config = oauth2.Config{
		ClientID:     config.C.AuthOIDCConfig.ClientID,
		ClientSecret: config.C.AuthOIDCConfig.ClientSecret,
		RedirectURL:  config.C.AuthOIDCConfig.Redirect,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	oa.verifier = provider.Verifier(&oidc.Config{ClientID: oa.config.ClientID})
}

func (oa *OIDCAuth) InitRoutes(r *mux.Router) {
	r.Path("/oidc/redirect").HandlerFunc(oa.handleRedirect)
	r.Path("/oidc/callback").HandlerFunc(oa.handleOAuth2Callback)
}

func (oa *OIDCAuth) AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	session, _ := store.Get(r, "session-name")
	if _, ok := session.Values["oidc_state"]; !ok {
		http.Redirect(w, r, "/api/pub/oidc/redirect", http.StatusFound)
	}
}
