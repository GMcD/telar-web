package config

import "time"

type (
	Configuration struct {
		OAuthProvider          string
		OAuthProviderBaseURL   string
		OAuthTelarBaseURL      string
		ClientID               string
		ClientSecret           string
		OAuthClientSecretPath  string // OAuthClientSecretPath when given overrides the ClientSecret env-var
		ExternalRedirectDomain string
		Scope                  string
		CookieRootDomain       string
		CookieExpiresIn        time.Duration
		BaseRoute              string
		VerifyType             string
		QueryPrettyURL         bool
		Debug                  bool // Debug enables verbose logging of claims / cookies
	}
)

// AuthConfig holds the configuration values from auth-config.yml file
var AuthConfig Configuration
