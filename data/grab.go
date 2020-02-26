package grab

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

type MyArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type MyArtistFull struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
}

type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type MyRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type MyDates struct {
	Index []MyDate `json:"index"`
}

type MyLocations struct {
	Index []MyLocation `json:"index"`
}
type MyRelations struct {
	Index []MyRelation `json:"index"`
}

var ArtistsFull []MyArtistFull
var Artists []MyArtist
var Dates MyDates
var Locations MyLocations
var Relations MyRelations

func GetArtistsData() error {
	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}

func GetDatesData() error {
	resp, err := http.Get(baseURL + "/dates")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Dates)
	return nil
}

func GetLocationsData() error {
	resp, err := http.Get(baseURL + "/locations")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Locations)
	return nil
}

func GetRelationsData() error {
	resp, err := http.Get(baseURL + "/relation")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Relations)
	return nil
}

func GetData() error {
	err1 := GetArtistsData()
	err2 := GetLocationsData()
	err3 := GetDatesData()
	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("Error by get data artists, locations, dates")
	}
	for i := range Artists {
		var tmpl MyArtistFull
		tmpl.ID = i + 1
		tmpl.Image = Artists[i].Image
		tmpl.Name = Artists[i].Name
		tmpl.Members = Artists[i].Members
		tmpl.CreationDate = Artists[i].CreationDate
		tmpl.FirstAlbum = Artists[i].FirstAlbum
		tmpl.Locations = Locations.Index[i].Locations
		tmpl.ConcertDates = Dates.Index[i].Dates
		ArtistsFull = append(ArtistsFull, tmpl)
	}
	return nil
}

func GetArtistByID(id int) (MyArtist, error) {
	for _, artist := range Artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtist{}, errors.New("Not found")
}

func GetDateByID(id int) (MyDate, error) {
	for _, date := range Dates.Index {
		if date.ID == id {
			return date, nil
		}
	}
	return MyDate{}, errors.New("Not found")
}

func GetLocationByID(id int) (MyLocation, error) {
	for _, location := range Locations.Index {
		if location.ID == id {
			return location, nil
		}
	}
	return MyLocation{}, errors.New("Not found")
}

func GetRelationByID(id int) (MyRelation, error) {
	for _, relation := range Relations.Index {
		if relation.ID == id {
			return relation, nil
		}
	}
	return MyRelation{}, errors.New("Not found")
}

func GetFullDataById(id int) (MyArtistFull, error) {
	for _, artist := range ArtistsFull {
		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtistFull{}, errors.New("Not found")
}

/*
func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	if Artists == nil {
		err := getArtistsData()
		if err != nil {
			internalError(w, r)
			return
		}
	}
	ids := r.URL.Query()["id"]
	if len(ids) > 1 {
		badRequest(w, r)
		return
	}

	if ids != nil {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			badRequest(w, r)
			return
		}
		artist, err := getArtistByID(id)
		if err != nil {
			internalError(w, r)
			return
		}
		response, err := json.Marshal(artist)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		response, err := json.Marshal(Artists)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func getDatesHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	if Dates.Index == nil {
		err := getDatesData()
		if err != nil {
			internalError(w, r)
			return
		}
	}
	ids := r.URL.Query()["id"]
	if len(ids) > 1 {
		badRequest(w, r)
		return
	}

	if ids != nil {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			badRequest(w, r)
			return
		}
		date, err := getDateByID(id)
		if err != nil {
			internalError(w, r)
			return
		}
		response, err := json.Marshal(date)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		response, err := json.Marshal(Dates)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	if Locations.Index == nil {
		err := getLocationsData()
		if err != nil {
			internalError(w, r)
			return
		}
	}
	ids := r.URL.Query()["id"]
	if len(ids) > 1 {
		badRequest(w, r)
		return
	}

	if ids != nil {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			badRequest(w, r)
			return
		}
		location, err := getLocationByID(id)
		if err != nil {
			internalError(w, r)
			return
		}
		response, err := json.Marshal(location)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		response, err := json.Marshal(Locations)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func getRelationsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	if Relations.Index == nil {
		err := getRelationsData()
		if err != nil {
			internalError(w, r)
			return
		}
	}
	ids := r.URL.Query()["id"]
	if len(ids) > 1 {
		badRequest(w, r)
		return
	}

	if ids != nil {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			badRequest(w, r)
			return
		}
		relation, err := getRelationByID(id)
		if err != nil {
			internalError(w, r)
			return
		}
		response, err := json.Marshal(relation)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		response, err := json.Marshal(Relations)
		if err != nil {
			internalError(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func fileServerWithCustom404(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			notFound(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/api/artists", getArtistsHandler)
	http.HandleFunc("/api/dates", getDatesHandler)
	http.HandleFunc("/api/locations", getLocationsHandler)
	http.HandleFunc("/api/relation", getRelationsHandler)
	http.Handle("/", fileServerWithCustom404(http.Dir("./frontend/dist")))
	http.ListenAndServe(":3000", nil)
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Vary", "Accept-Encoding")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./frontend/dist/404.html"))
}

func internalError(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./frontend/dist/500.html"))
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./frontend/dist/400.html"))
}
*/
