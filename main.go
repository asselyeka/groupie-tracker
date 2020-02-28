package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	grab "./data"
)

var data []grab.MyArtistFull

func mainPage(w http.ResponseWriter, r *http.Request) {
	err := grab.GetData()
	if err != nil {
		errors.New("Error by get data")
	}

	body := r.FormValue("search")
	data = Search(body)

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

func concertPage(w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("concert")
	id, _ := strconv.Atoi(idStr)
	artist, _ := grab.GetFullDataById(id)

	for key, value := range artist.DatesLocations {
		fmt.Print(key + "  - ")
		for _, e := range value {
			println(e)
		}
	}

	tmpl, err := template.ParseFiles("concert.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := tmpl.Execute(w, artist); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func ConverterStructToString() ([]string, error) {
	var data []string
	for i := 1; i <= len(grab.Artists); i++ {
		artist, err1 := grab.GetArtistByID(i)
		locations, err2 := grab.GetLocationByID(i)
		date, err3 := grab.GetDateByID(i)
		if err1 != nil || err2 != nil || err3 != nil {
			return data, errors.New("Error by converter")
		}

		str := artist.Name + " "
		for _, member := range artist.Members {
			str += member + " "
		}
		str += strconv.Itoa(artist.CreationDate) + " "
		str += artist.FirstAlbum + " "
		for _, location := range locations.Locations {
			str += location + " "
		}
		for _, d := range date.Dates {
			str += d + " "
		}
		data = append(data, str)
	}
	println("Convert to str Done!")
	return data, nil
}

func Search(search string) []grab.MyArtistFull {
	if search == "" {
		return grab.ArtistsFull
	}
	art, err := ConverterStructToString()
	if err != nil {
		errors.New("Error by converter")
	}
	var search_artist []grab.MyArtistFull

	for i, artist := range art {
		lower_band := strings.ToLower(artist)
		for i_name, l_name := range []byte(lower_band) {
			lower_search := strings.ToLower(search)
			if lower_search[0] == l_name {
				lenght_name := 0
				indx := i_name
				for _, l := range []byte(lower_search) {
					if l == lower_band[indx] {
						if indx+1 == len(lower_band) {
							break
						}
						indx++
						lenght_name++
					} else {
						break
					}
				}
				if len(search) == lenght_name {
					band, _ := grab.GetFullDataById(i + 1)
					search_artist = append(search_artist, band)
					break
				}
			}
		}

	}
	println("Search str Done!")
	return search_artist
}

func main() {

	// static folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/concert", concertPage)

	port := ":8080"
	println("Server listen on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Listen and Server", err)
	}

}
