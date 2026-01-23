package config

import "os"

func GetMPAccessToken() string {
	token := os.Getenv("MP_ACCESS_TOKEN")
	if token == "" {
		return "YOUR_MERCADO_PAGO_ACCESS_TOKEN_HERE"
	}
	return token
}

func GetStripeSecretKey() string {
	token := os.Getenv("STRIPE_SECRET_KEY")
	if token == "" {
		return "your_stripe_secret_key"
	}
	return token
}

func GetStripeWebhookSecret() string {
	token := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if token == "" {
		return "your_stripe_webhook_secret"
	}
	return token
}

func GetJWTSecret() string {
	token := os.Getenv("JWT_SECRET")
	if token == "" {
		return "your_jwt_secret_key"
	}
	return token
}
