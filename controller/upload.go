package controller

import (
	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadImage => Upload Image to Cloudinary and get the URL response
func UploadImage(c *gin.Context) {
	// accept the images and store it in the tempfile
	c.Request.ParseMultipartForm(32 << 20)

	file, _, err := c.Request.FormFile("myFile")
	defer file.Close()

	// make it in the same folder as this file
	absolutePath, _ := filepath.Abs("./")

	// make a temporary file in the disk
	// will be deleted after uploading finishes
	tempFile, err := ioutil.TempFile(absolutePath, "upload-*.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to create temporary file for uploading"})

		return
	}

	// TO DO: remove temporary file.
	// defer tempFile.Close()
	// defer os.RemoveAll(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to read the temporary file"})

		return
	}

	fileName := tempFile.Name()

	tempFile.Write(fileBytes)

	url, err := config.CloudinaryService.UploadImage(fileName, nil, "")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to upload the file"})

		return
	}

	// 0 represent the iota or code for image in go-cloudinary
	urlImage := config.CloudinaryService.Url(url, 0)

	// remove the file after using
	err = tempFile.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to close the temporary file"})

		return
	}

	err = os.RemoveAll(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to remove the temporary file"})

		return
	}

	// response OK
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Uploading successful",
		"data":    urlImage})

	return
}
