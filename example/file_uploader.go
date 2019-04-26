package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/utahta/go-twitter"
)

// Example file uploader. Set all credentials, set image to be downloaded and attached to tweet.
func main() {

	twitter.SetConsumerCredentials(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))

	cl, err := twitter.New(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		fmt.Printf("err accessing twitter application :%v \n", err)
		return
	}

	imgUrl := "https://www.igneous.io/hs-fs/hubfs/gopher3.png?width=400&height=214&name=gopher3.png" //gopher
	_, err = cl.TweetImageURLs("testing img upload", []string{imgUrl}, url.Values{})
	if err != nil {
		fmt.Printf("error creating tweet with attachment: %v \n", err)
		return
	}
	fmt.Println("tweet with attachment created")
}
