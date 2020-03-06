package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	grab "./data"
)

var data []grab.MyArtistFull

func mainPage(w http.ResponseWriter, r *http.Request) {
	err := grab.GetData()
	if err != nil {
		errors.New("Error by get data")
	}

	main := r.FormValue("main")
	search := r.FormValue("search")
	filterByCreationFrom := r.FormValue("startCD")
	filterByCreationTill := r.FormValue("endCD")
	location := r.FormValue("location-filter")
	var mem1, mem2, mem3, mem4, mem5, mem6, mem7, mem8 int
	mem1, err1 := strconv.Atoi(r.FormValue("mem1"))
	if err1 != nil {
		mem1 = 0
	}
	mem2, err2 := strconv.Atoi(r.FormValue("mem2"))
	if err2 != nil {
		mem2 = 0
	}
	mem3, err3 := strconv.Atoi(r.FormValue("mem3"))
	if err3 != nil {
		mem3 = 0
	}
	mem4, err4 := strconv.Atoi(r.FormValue("mem4"))
	if err4 != nil {
		mem4 = 0
	}
	mem5, err5 := strconv.Atoi(r.FormValue("mem5"))
	if err5 != nil {
		mem5 = 0
	}
	mem6, err6 := strconv.Atoi(r.FormValue("mem6"))
	if err6 != nil {
		mem6 = 0
	}
	mem7, err7 := strconv.Atoi(r.FormValue("mem7"))
	if err7 != nil {
		mem7 = 0
	}
	mem8, err8 := strconv.Atoi(r.FormValue("mem8"))
	if err8 != nil {
		mem8 = 0
	}
	mem := []int{mem1, mem2, mem3, mem4, mem5, mem6, mem7, mem8}
	sum := 0
	for _, n := range mem {
		sum += n
	}
	println("startCD:", filterByCreationFrom)
	fmt.Println("mem:", mem)

	filterByFA := r.FormValue("startFA")
	filterByFAend := r.FormValue("endFA")
	fmt.Println("filterFA:", filterByFA)
	fmt.Println("filterFAend:", filterByFAend)

	if !(search == "" && len(data) != 0) {
		data = Search(search)
	}

	if filterByCreationFrom != "" || filterByCreationTill != "" {
		if filterByCreationFrom == "" {
			filterByCreationFrom = "1900"
		}
		if filterByCreationTill == "" {
			filterByCreationTill = "2020"
		}
		data = FilterByCreation(data, filterByCreationFrom, filterByCreationTill)
	}

	if filterByFA != "" || filterByFAend != "" {
		if filterByFA == "" {
			filterByFA = "1900-01-01"
		}
		if filterByFAend == "" {
			filterByFAend = "2020-03-03"
		}
		data = FilterByAlbumDate(data, filterByFA, filterByFAend)
	}

	if sum != 0 {
		data = FilterByMember(data, mem)
	}

	if location != "" {
		data = FilterByLocation(data, location)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if main == "Main Page" {
		data = Search("a")
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

func FilterByCreation(data []grab.MyArtistFull, from string, till string) []grab.MyArtistFull {
	if from == "" || till == "" {
		return data
	}
	fromInt, err1 := strconv.Atoi(from)
	tillInt, err2 := strconv.Atoi(till)
	if err1 != nil || err2 != nil {
		errors.New("Error by filter by creation date data")
	}

	var search_artist []grab.MyArtistFull

	for _, artist := range data {
		if artist.CreationDate >= fromInt && artist.CreationDate <= tillInt {
			search_artist = append(search_artist, artist)
		}
	}
	return search_artist
}

func FilterByMember(data []grab.MyArtistFull, mem []int) []grab.MyArtistFull {
	var tmp []grab.MyArtistFull
	for _, bandmem := range data {
		for _, num := range mem {
			if len(bandmem.Members) == num {
				tmp = append(tmp, bandmem)
			}

		}
	}
	return tmp
}

func FilterByAlbumDate(data []grab.MyArtistFull, albumStart string, albumEnd string) []grab.MyArtistFull {
	layOut := "2006-01-02"

	dateStart, err1 := time.Parse(layOut, albumStart)
	dateEnd, err2 := time.Parse(layOut, albumEnd)
	if err1 != nil || err2 != nil {
		errors.New("Error by First Album date convert for filter")
	}
	var tmp []grab.MyArtistFull

	layOutData := "02-01-2006"
	for _, band := range data {
		date, err := time.Parse(layOutData, band.FirstAlbum)
		if err != nil {
			errors.New("Error by First Album date convert for filter")
		}
		if dateStart.Before(date) && dateEnd.After(date) {
			tmp = append(tmp, band)
		}
	}
	return tmp
}

func FilterByLocation(data []grab.MyArtistFull, location string) []grab.MyArtistFull {
	var tmp []grab.MyArtistFull

	for _, band := range data {
		for _, loc := range band.Locations {
			if location == loc {
				tmp = append(tmp, band)
			}
		}
	}
	return tmp
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
