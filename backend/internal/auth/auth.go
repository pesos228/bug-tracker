package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pesos228/bug-tracker/internal/config"
	"golang.org/x/oauth2"
)

type Client struct {
	Provider *oidc.Provider
	OIDC     *oidc.IDTokenVerifier
	Oauth    oauth2.Config
}

func New(ctx context.Context, cfg *config.AuthConfig) (*Client, error) {
	discoveryUrl := fmt.Sprintf("%s/realms/%s", cfg.InternalBaselUrl, cfg.Realm)
	expectedIssuerUrl := fmt.Sprintf("%s/realms/%s", cfg.PublicBaseUrl, cfg.Realm)
	newCtx := oidc.InsecureIssuerURLContext(ctx, expectedIssuerUrl)

	provider, err := oidc.NewProvider(newCtx, discoveryUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	internalJWKSURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.InternalBaselUrl, cfg.Realm)
	keySet := oidc.NewRemoteKeySet(ctx, internalJWKSURL)

	verifierConfig := &oidc.Config{
		ClientID: cfg.ClientId,
	}

	verifier := oidc.NewVerifier(expectedIssuerUrl, keySet, verifierConfig)

	internalEndpoint := oauth2.Endpoint{
		AuthURL:  provider.Endpoint().AuthURL,
		TokenURL: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", cfg.InternalBaselUrl, cfg.Realm),
	}

	ouath2Config := oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectUrl,
		Endpoint:     internalEndpoint,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	return &Client{
		Provider: provider,
		OIDC:     verifier,
		Oauth:    ouath2Config,
	}, nil
}
