package config

import "github.com/stevejo12/PMSFreelancer/models"

// import "PMSFreelancer/models"

// Config => Configuration for admin key to access swagger
var Config models.AuthenticationKeyForSwagger

// LoadConfig => Initialize configuration for header
func LoadConfig() {
	Config.APIKey = "AuthorizationSPIRITS"
}
