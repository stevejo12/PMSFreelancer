package config

import (
	"github.com/stevejo12/go-cloudinary"
)

// CloudinaryService => Cloudinary connection
var CloudinaryService *cloudinary.Service

// ConnectToCloudinary => Establish connection to Cloudinary API
func ConnectToCloudinary() {
	var err error

	apiKey := "645723149374711"
	apiSecret := "nr7-pjsZYMnZxawJIFcUNTonh8g"
	cloudName := "dvah7jvpa"

	uri := "cloudinary://" + apiKey + ":" + apiSecret + "@" + cloudName

	CloudinaryService, err = cloudinary.Dial(uri)

	if err != nil {
		panic(err.Error())
	}
}
