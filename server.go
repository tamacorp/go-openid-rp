package main

import(
  "encoding/json"
  "log"
  "net/http"
  oidc "github.com/coreos/go-oidc"
  "golang.org/x/net/context"
  "golang.org/x/oauth2"
)

var config oauth2.Config
var clientID string
var clientSecret string
var state string
var verifier *oidc.IDTokenVerifier

func main() {

  clientID = "test"
  clientSecret = "test"
  state = "12345"

  ctx := context.Background()
  provider, err := oidc.NewProvider(ctx, "https://sso.tamacorp.co/oauth")
  if err != nil {
    log.Fatal(err)
  }
  oidcConfig := &oidc.Config{
    ClientID: clientID,
  }
  verifier = provider.Verifier(oidcConfig)

  config = oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    Endpoint:     provider.Endpoint(),
    RedirectURL:  "http://localhost:3000/assert",
    Scopes:       []string{"openid", "profile"},
  }

  http.HandleFunc("/login", login)
  http.HandleFunc("/assert", assert)
  log.Println("RP is running on port 3000")
  http.ListenAndServe(":3000", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
}

func assert(w http.ResponseWriter, r *http.Request) {
  if r.URL.Query().Get("state") != state {
    http.Error(w, "state did not match", http.StatusBadRequest)
    return
  }

  oauth2Token, err := config.Exchange(context.Background(), r.URL.Query().Get("code"))
  if err != nil {
    http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
    return
  }
  rawIDToken, ok := oauth2Token.Extra("id_token").(string)
  if !ok {
    http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
    return
  }
  idToken, err := verifier.Verify(context.Background(), rawIDToken)
  if err != nil {
    http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
    return
  }

  resp := struct {
    OAuth2Token   *oauth2.Token
    IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
  }{oauth2Token, new(json.RawMessage)}

  if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  data, err := json.MarshalIndent(resp, "", "    ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(data)
}

