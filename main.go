package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/naoya/go-pit"
	"log"
	"os"
	"os/exec"
)

func main() {
	config := pit.Get("cdn.bloghackers.net")

	var accessKeyID, secretAccessKey string
	if config["aws_access_key_id"] == "" {
		accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	} else {
		accessKeyID = config["aws_access_key_id"]
		secretAccessKey = config["aws_secret_access_key"]
	}

	// TODO: 設定ファイルとかで設定できるように
	gozo := NewGozo(
		accessKeyID,
		secretAccessKey,
		"files.bloghackers.net",
		aws.APNortheast,
		"http://cdn.bloghackers.net/",
	)

	var url string
	var err error

	if len(os.Args) == 1 {
		url, err = gozo.SendCapture()
	} else {
		url, err = gozo.SendFile(os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}

	if url == "" {
		os.Exit(1)
	}

	err = exec.Command("open", url).Run()
	if err != nil {
		log.Fatal(err)
	}
}
