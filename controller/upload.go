package controller

import (
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"

	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func uploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	arr := helpers.SplitDot(header.Filename)
	var fileExt string

	if len(arr) == 2 {
		// take the second part since it is the extension file
		fileExt = arr[1]
	} else {
		return "", errors.New("File format is not correct ex: *.pdf or *.txt")
	}

	// make it in the same folder as this file
	absolutePath, _ := filepath.Abs("./")

	// make a temporary file in the disk
	// will be deleted after uploading finishes
	tempFile, err := ioutil.TempFile(absolutePath, "upload-*."+fileExt)
	if err != nil {
		return "", errors.New("Server unable to create temporary file for uploading")
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		return "", errors.New("Server unable to read the temporary file")
	}

	fileName := tempFile.Name()

	tempFile.Write(fileBytes)

	// TO DO: Mungkin bisa dibuat restriction berdasarkan tipe file yg diupload
	// filename itu extension nya

	var url string
	if fileExt == "pdf" {
		url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 1)
	} else {
		url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 3)
	}
	// url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 3)

	if err != nil {
		return "", errors.New("Server is unable to upload the file")
	}

	// 0 represent the iota or code for image in go-cloudinary
	var urlFile string
	if fileExt == "pdf" {
		urlFile = config.CloudinaryService.Url(url, 1)
	} else {
		urlFile = config.CloudinaryService.Url(url, 3)
	}
	// urlFile = config.CloudinaryService.Url(url, 3)

	// remove the file after using
	err = tempFile.Close()
	if err != nil {
		return "", errors.New("Server is unable to close the temporary file")
	}

	err = os.RemoveAll(tempFile.Name())
	if err != nil {
		return "", errors.New("Server is unable to remove the temporary file")
	}

	return urlFile, nil
}

// UploadPicture => Upload Image to Cloudinary and get the URL response
// UploadPicture godoc
// @Summary Uploading Picture here
// @Tags User
// @Accept multipart/form-data
// @Param token header string true "Token Header"
// @Param file formData file true "Upload File"
// @Success 200 {object} models.ResponseWithStringData
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /uploadPicture [post]
func UploadPicture(c *gin.Context) {
	id := idToken

	// accept the images and store it in the tempfile
	c.Request.ParseMultipartForm(32 << 20)

	file, _, err := c.Request.FormFile("file")
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

	// TO DO: Store the image into an id (add parameter that accept user id)
	_, err = config.DB.Query("UPDATE login SET picture=? WHERE id=?", urlImage, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to store url in the database"})
		return
	}

	// response OK
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Uploading successful",
		"data":    urlImage})

	return
}

// UploadAttachment => Upload file other than image to this
// UploadAttachment godoc
// @Summary Uploading Attachment here
// @Tags Project
// @Accept multipart/form-data
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Param file formData file true "Upload File"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /uploadAttachment/{id} [post]
func UploadAttachment(c *gin.Context) {
	// accept the images and store it in the tempfile
	c.Request.ParseMultipartForm(5 * 1024 * 1024)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to read the uploaded file"})
		return
	}
	defer file.Close()

	arr := helpers.SplitDot(header.Filename)
	var fileExt string

	if len(arr) == 2 {
		// take the second part since it is the extension file
		fileExt = arr[1]
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "File format is not correct ex: *.pdf or *.txt"})
		return
	}

	// make it in the same folder as this file
	absolutePath, _ := filepath.Abs("./")

	// make a temporary file in the disk
	// will be deleted after uploading finishes
	tempFile, err := ioutil.TempFile(absolutePath, "upload-*."+fileExt)
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

	// TO DO: Mungkin bisa dibuat restriction berdasarkan tipe file yg diupload
	// filename itu extension nya

	var url string
	if fileExt == "pdf" {
		url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 1)
	} else {
		url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 3)
	}
	// url, err = config.CloudinaryService.Upload(fileName, nil, "", true, 3)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to upload the file"})

		return
	}

	// 0 represent the iota or code for image in go-cloudinary
	var urlFile string
	if fileExt == "pdf" {
		urlFile = config.CloudinaryService.Url(url, 1)
	} else {
		urlFile = config.CloudinaryService.Url(url, 3)
	}
	// urlFile = config.CloudinaryService.Url(url, 3)

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
		"data":    gin.H{"url": urlFile}})
	return
}
