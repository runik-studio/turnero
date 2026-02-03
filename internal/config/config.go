package config

import (
	"encoding/json"
	"log"
	"os"
)

func GetFirebaseProjectID() string {
	if id := os.Getenv("FIREBASE_PROJECT_ID"); id != "" {
		return id
	}

	data, err := os.ReadFile("firebaseCredentials.json")
	if err != nil {
		log.Printf("Warning: Could not read firebaseCredentials.json to find project_id: %v", err)
		return ""
	}

	var creds map[string]interface{}
	if err := json.Unmarshal(data, &creds); err != nil {
		log.Printf("Warning: Could not parse firebaseCredentials.json: %v", err)
		return ""
	}

	if id, ok := creds["project_id"].(string); ok && id != "" {
		return id
	}

	return ""
}

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
