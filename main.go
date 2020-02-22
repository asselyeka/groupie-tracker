package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var art []Artists

type Artists struct {
	Id           int      `json:id`
	Image        string   `json:image`
	Name         string   `json:name`
	Members      []string `json:members`
	CreationDate int      `json:creationDate`
	FirstAlbum   string   `json:firstAlbum`
	Locations    string   `json:locations`
	ConcertDates string   `json:concertDates`
	Relations    string   `json:relations`
}

func ParseArtists(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	ourStr, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		panic(err1)
	}

	json.Unmarshal(ourStr, &art)

	fmt.Println("Artists")
}

// locations := "https://groupietrackers.herokuapp.com/api/locations"
// dates := "https://groupietrackers.herokuapp.com/api/dates"
// relation :=	"https://groupietrackers.herokuapp.com/api/relation" (

var artists = "https://groupietrackers.herokuapp.com/api/artists"

func mainPage(w http.ResponseWriter, r *http.Request) {

	body := r.FormValue("body")

	//font := r.FormValue("fonts")

	data := SearchArtist(body)

	for i := range data {
		fmt.Println(data[i].CreationDate)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func SearchArtist(name string) []Artists {
	if name == "" {
		return nil
	}
	var search_artist []Artists

	for i, artist := range art {
		if name[0] == artist.Name[0] {
			lenght_name := 0
			for _, l := range []byte(name) {
				if l == art[i].Name[lenght_name] {
					lenght_name++
				} else {
					break
				}
			}
			if len(name) == lenght_name {
				search_artist = append(search_artist, art[i])
			}
		}
	}
	return search_artist
}

func main() {

	ParseArtists(artists)
	println(art[0].Name)

	http.HandleFunc("/", mainPage)

	port := ":8080"
	println("Server listen on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Listen and Server", err)
	}

}
