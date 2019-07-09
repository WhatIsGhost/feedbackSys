package minio

import (
	"feedbackSys/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/minio/minio-go"
	"log"
	"strings"
)

var minioClient *minio.Client
var location = "us-east-1"

func InitMinio() {
	var err error
	endpoint := beego.AppConfig.String("minio::endpoint")
	accessKeyID := beego.AppConfig.String("minio::accessKeyID")
	secretAccessKey := beego.AppConfig.String("minio::secretAccessKey")
	useSSL, _ := beego.AppConfig.Bool("minio::useSSL")
	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		utils.LogError(fmt.Sprintf("init minio error %s", err))
	}

}

func MakeBucket(bucketName string) {
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil || !exists {
		err := minioClient.MakeBucket(bucketName, location)
		if err != nil {
			utils.LogError(fmt.Sprintf("make bucket error %s", err))
		}
	}
}

func UploadFile(bucketName string, filePath string) string {
	MakeBucket(bucketName)
	pathSlice := strings.Split(filePath, "/")
	objectName := pathSlice[len(pathSlice)-1]
	nameSlice := strings.Split(objectName, ".")
	suffix := nameSlice[len(nameSlice)-1]
	contentType := "application/" + suffix
	num, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		utils.LogError(fmt.Sprintf("upload file error %s", err))
		return ""
	}
	utils.LogDebug(fmt.Sprintf("Successfully uploaded %s of size %d\n", objectName, num))
	return "http://" + beego.AppConfig.String("minio::endpoint") + "/" + bucketName + "/" + objectName
}
