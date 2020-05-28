package config

import (
	"github.com/stevejo12/go-cloudinary"
)

// CloudinaryService => Cloudinary connection
var CloudinaryService *cloudinary.Service

// ConnectToCloudinary => Establish connection to Cloudinary API
func ConnectToCloudinary() {
	var err error

	apiKey := "477787745735813"
	apiSecret := "Nac4woUWVnBOttHOSsSiaI5PdFc"
	cloudName := "drrd7ai50"

	uri := "cloudinary://" + apiKey + ":" + apiSecret + "@" + cloudName

	CloudinaryService, err = cloudinary.Dial(uri)

	if err != nil {
		panic(err.Error())
	}
}
