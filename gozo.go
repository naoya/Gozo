package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"log"
	"time"
)

type Gozo struct {
	accessKey       string
	secretAccessKey string
	bucketName      string
	region          aws.Region
	rootUrl         string
}

func NewGozo(accessKey string, secretAccessKey string, bucketName string, region aws.Region, rootUrl string) *Gozo {
	gozo := &Gozo{
		accessKey,
		secretAccessKey,
		bucketName,
		region,
		rootUrl,
	}
	return gozo
}

func (gozo Gozo) SendImage(filename string) (url string, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	auth, err := aws.GetAuth(
		gozo.accessKey,
		gozo.secretAccessKey,
	)
	if err != nil {
		return
	}

	s3client := s3.New(auth, gozo.region)
	bucket := s3client.Bucket(gozo.bucketName)

	path := "images/" + hexdigest(fmt.Sprintf("%s-%d", filename, time.Now().Unix())) + ".png"
	err = bucket.Put(path, data, "image/png", s3.PublicRead)
	if err != nil {
		return
	}

	url = gozo.rootUrl + path
	return
}

func (gozo Gozo) SendCapture() (url string, err error) {
	img, err := CaptureImage()
	if err != nil {
		return
	}

	if img == nil {
		url = ""
		return
	}

	defer func() {
		if img != nil {
			if err := img.Remove(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	url, err = gozo.SendImage(img.Name())
	return
}

func (gozo Gozo) SendFile(filename string) (url string, err error) {
	img := new(Image)
	img.name = filename
	if err = img.RemoveProfile(); err != nil {
		return
	}

	url, err = gozo.SendImage(img.Name())
	return
}

func hexdigest(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
