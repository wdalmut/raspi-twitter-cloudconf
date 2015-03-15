package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/darkhelmet/twitterstream"
	"github.com/wdalmut/twitterstream/async"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Key            string
	Secret         string
}

func main() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.WithFields(log.Fields{
			"type": "bootstrap",
		}).Fatal("Unable to find the configuration file 'config.json'")
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)

	client := async.NewClient(
		config.ConsumerKey,
		config.ConsumerSecret,
		config.AccessToken,
		config.AccessSecret,
	)

	AWSAuth := aws.Auth{
		AccessKey: config.Key,
		SecretKey: config.Secret,
	}

	region := aws.EUWest

	connection := s3.New(AWSAuth, region)
	bucket := connection.Bucket("example.walterdalmut.com")

	client.TrackAndServe("cloudconf2015", func(tweet *twitterstream.Tweet) {
		text := tweet.Text
		name := tweet.User.ScreenName

		log.WithFields(log.Fields{
			"user": name,
			"type": "tweet",
		}).Info(text)

		filename := strconv.FormatInt(tweet.Id, 10)

		cmd := exec.Command(
			"raspistill",
			"-a", "www.cloudconf.it - #cloudconf2015",
			"-t", "500",
			"-vf", "-hf",
			"-w", "1024", "-h", "768",
			"--quality", "60",
			"-o", "/tmp/pic.jpg")
		err := cmd.Run()

		if err != nil {
			log.WithFields(log.Fields{
				"user": name,
				"type": "pic",
			}).Error("Unable to get the picture, send the latest picture...")
		}

		path := name + "/" + filename + ".jpg"
		fileToBeUploaded := "/tmp/pic.jpg"

		file, err := os.Open(fileToBeUploaded)
		defer file.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		fileInfo, _ := file.Stat()
		var size int64 = fileInfo.Size()
		bytes := make([]byte, size)

		buffer := bufio.NewReader(file)
		_, err = buffer.Read(bytes)

		err = bucket.Put(path, bytes, "image/png", s3.ACL("public-read"), s3.Options{})

		if err != nil {
			log.WithFields(log.Fields{
				"user": name,
				"type": "upload",
			}).Error("Unable to upload the picture on S3!!! I have to exit...")

			return
		}

		log.WithFields(log.Fields{
			"user": name,
			"type": "tweet",
		}).Info(fmt.Printf("Tweet correctly uploaded! %s", text))
	})
}
