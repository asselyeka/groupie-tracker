package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var art Artists

type Artists []struct {
	Id                int      `json:id`
	Image             string   `json:image`
	Name              string   `json:name`
	Members           []string `json:members`
	CreationDate_year int      `json:creationDate`
	FirstAlbum_date   int      `json:firstAlbum`
	Locations         string   `json:locations`
	ConcertDates      string   `json:concertDates`
	Relations         string   `json:relations`
}

func ParseArtists(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	ourStr, orr := ioutil.ReadAll(response.Body)
	if orr != nil {
		panic(orr)
	}

	json.Unmarshal(ourStr, &art)

	fmt.Println("Artists")
}

func main() {
	artists := "https://groupietrackers.herokuapp.com/api/artists"
	// locations := "https://groupietrackers.herokuapp.com/api/locations"
	// dates := "https://groupietrackers.herokuapp.com/api/dates"
	// relation :=	"https://groupietrackers.herokuapp.com/api/relation"
	ParseArtists(artists)
	println(art[0].Name)

}
