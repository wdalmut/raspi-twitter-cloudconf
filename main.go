package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/darkhelmet/twitterstream"
	"github.com/wdalmut/twitterstream/async"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
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
		fmt.Printf("File error: %v\n", e)
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

		rand.Seed(42)
		num := rand.Intn(1000000)

		filename := strconv.Itoa(num)

		cmd := exec.Command("raspistill", "--quality", "10", "-o", "/tmp/pic.jpg")
		err := cmd.Run()

		if err != nil {
			fmt.Println("Unable to get the picture!")
			return
		}

		path := name + "/" + filename + ".jpg"
		fileToBeUploaded := "/tmp/pic.jpg"

		file, err := os.Open(fileToBeUploaded)

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
			fmt.Println(err)
		}
		file.Close()

		fmt.Printf("Tweet: %s is %s\n\n", text, name)

	})
}
