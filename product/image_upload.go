package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ph *ProductHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	fileSize := file.Size
	if fileSize < 10240 { // Less than 10KB
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size is too small (less than 10KB)"})
		return
	}
	if fileSize > 2097152 { // More than 2MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size is too large (more than 2MB)"})
		return
	}

	fileExt := filepath.Ext(file.Filename)
	if fileExt != ".jpg" && fileExt != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File format not supported (must be *.jpg or *.jpeg)"})
		return
	}

	// Generate UUID v4 for filename
	uuid := uuid.New()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), ""),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	// Get the file from the form data
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Specify the bucket and the key (filename) under which you want to store the image in S3
	bucketName := os.Getenv("S3_BASE_URL")
	key := uuid.String() + fileExt

	// Upload the image to S3
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(key),
		Body:               src,
		ACL:                aws.String("public-read"),
		ContentLength:      aws.Int64(fileSize),
		ContentType:        aws.String(file.Header.Get("Content-Type")),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file to S3: %s", err.Error())})
		return
	}

	// Send success response with image URL
	imageURL := fmt.Sprintf("https://%s/%s", bucketName, key)
	c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
}
