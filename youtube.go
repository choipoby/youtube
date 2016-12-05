package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	url = flag.String("url", "", "")
)

func findVideoId(url string) string {
	//var videoID string
	videoURL := url
	if strings.Contains(videoURL, "youtu") || strings.ContainsAny(videoURL, "\"?&/<%=") {
		re_list := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range re_list {
			if is_match := re.MatchString(videoURL); is_match {
				subs := re.FindStringSubmatch(videoURL)
				fmt.Println(subs)
				videoURL = subs[1]
			}
		}
	}

	log.Printf("Found video id: '%s'", videoURL)
	return videoURL
}

func main() {
	flag.Parse()

	videoid := findVideoId(*url)

	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyDNHK2hPPpLw9lB8_cBWBkVYVWsdYTrc-0"},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	videosService := youtube.NewVideosService(service)

	// Make the API call to YouTube.
	call := videosService.List("snippet").
		Id(videoid).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	/*
		//func NewThumbnailsService(s *Service) *ThumbnailsService
		thumbnailService := youtube.NewThumbnailsService(service)

		//func (r *ThumbnailsService) Set(videoId string) *ThumbnailsSetCall
		thumnailServiceSet := thumbnailService.Set(videoid)

		//func (c *ThumbnailsSetCall) Do(opts ...googleapi.CallOption) (*ThumbnailSetResponse, error)
		response, err := thumnailServiceSet.Do()
		if err != nil {
			log.Fatalf("Error making search API call: %v", err)
		}
	*/
	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		// fmt.Println(item.ContentDetails)
		// fmt.Println(item.Id)
		// fmt.Println(item.Kind)
		fmt.Println(item.Snippet.Thumbnails.High.Url)
		fmt.Println(item.Snippet.Thumbnails.Medium.Url)
		fmt.Println(item.Snippet.Thumbnails.Standard.Url)
		fmt.Println(item.Snippet.Thumbnails.Default.Url)
		fmt.Println(item.Snippet.Description)
		fmt.Println(item.Snippet.Title)
		//fmt.Println(item)
	}
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
