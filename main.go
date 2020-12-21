package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"os/exec"
	"os"
	"sync"
	"strings"
)

type song struct {
	title string
	artist string
	album string
	length string
}


func main() {
	// Create the sync group
	var wg sync.WaitGroup

	// Get the playlist url and target dir from os args
	playlistURLLists := getURL(os.Args[1:])
	targetPath := getTargetPath(os.Args[1:])

	// Loop through each url and download them concurrently
	for _, playlistURL := range playlistURLLists {

		// Add the go func to the work group
		wg.Add(1)
		go getAndDownloadSong(&wg, playlistURL, targetPath)
	}

	// Join point
	wg.Wait()

	fmt.Println("Main: All jobs done")
}

func getAndDownloadSong(wg *sync.WaitGroup, playlistURL string, targetPath string) {

	// Will notify wg this function is done once its finished returning
	defer wg.Done()

	// Get the songs form the playlist
	songList :=	getSongs(playlistURL)
	
	// download each song
	for _, song := range songList {
		fmt.Println("Now downloading "+ song.title)
		searchTerm := ""+song.title+" by "+song.artist
		downloadSong(searchTerm, targetPath)
	}
}

func downloadSong(searchTerm string, outputPath string) {
	
	// Create the outputfile template
	outputFile := outputPath + "%(title)s.%(ext)s"

	// use the command line to execute youtube-dl to search and download the song
	cmd := exec.Command("cmd", "/c", "youtube-dl", "ytsearch1:"+searchTerm, "--add-metadata","-f bestaudio[ext=m4a]","-o"+outputFile)
	
	// Capture the output and any errors
	output, err := cmd.Output()

	// Output the error if any or else display the results
	if err != nil {
		fmt.Println("Errors: ", err)
	} else {
		fmt.Println(string(output))
	}
}


func getSongs(url string) []song {

	// Wil hold all out song info
	var songList []song

	// Get the url and html document
	res, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	// Finds the title of the playlist
	playListTitle := doc.Find("div > h1").First().Text()
	fmt.Println("Downloading: ", playListTitle)

	// Finds the element that holds the songs
	divList := doc.Find("[role]")

	// Loop through each song and get its title, artist, length and album
	divList.Each(func(i int, s *goquery.Selection) {

		span := s.Find("span")
		if len(span.Text())!=0 && span.Text() != "Now downloading To play this content, you'll need the Spotify app." {

			title := s.Find("span").First().Text()
			artist := s.Find("span").Eq(2).Text()
			length := s.Find("span").Eq(-2).Text()
			album := s.Find("span").Eq(-3).Text()

			newSong := song {
				title: title,
				artist: artist,
				album: album,
				length: length,
			}

			songList = append(songList, newSong)
			
		}
	})
	return songList
}

// Checks which arguments are urls
func getURL(args []string) []string {
	var playlistURL []string
	for _, arg := range args {
		if strings.Contains(arg, "https://") {
			playlistURL = append(playlistURL, arg)
		}
	}
	return playlistURL
}

// Checks which argument is a file path
func getTargetPath(arg []string) string {
	for _, arg := range arg {
		if !strings.Contains(arg, "https://") {
			return arg
		}
	}
	return `C:\Users\Shannon\Desktop\TEST_SONGS\`
}