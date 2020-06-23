package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/gojsonq"
)

// Define structure for json attributes
type adNetwork struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Platform    string `json:"platform"`
	OsVersion   string `json:"osversion"`
	AppName     string `json:"appname"`
	AppVersion  string `json:"appversion"`
	CountryCode string `json:"countrycode"`
	AdType      string `json:"adtype"`
}

// Global variable for a list of Ad Networks
var adNetworks []adNetwork

// Path to the file where json data is written
var filePath = "output.txt"

// Return all Ad Networks (Everything contained in the file)
func returnAllAdNetworks(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(adNetworks)
}

// Return a specific Ad Networks by ad type
func returnAdType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adtype := vars["adtype"]

	var list []adNetwork

	for _, adNetwork := range adNetworks {
		if strings.ToLower(adNetwork.AdType) == strings.ToLower(adtype) {
			list = append(list, adNetwork)
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(list)
}

// Custom Query for Ad Networks
func queryAdNetworks(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// Check if platform and osversion parameters were input
	platform, platformErr := params["platform"]
	osversion, osversionErr := params["osversion"]

	jq := gojsonq.New().File(filePath)
	res := jq.WhereNotNil("id")

	// Query json data with gotten parameters
	i := 0
	for param, valueParam := range params {
		if strings.ToLower(param) == "id" {
			value, err := strconv.ParseInt(valueParam[0], 10, 64)
			if err != nil {
				fmt.Println("Query ERROR")
				return
			}
			res = jq.WhereEqual(strings.ToLower(param), value)
		} else {
			value := valueParam[0]
			res = jq.WhereEqual(strings.ToLower(param), strings.ToLower(value))
		}

		// Don't show admob Networks when the given platform was "android" and the osversion was "9"
		if platformErr == true && osversionErr == true && i == 0 {
			i++
			if strings.ToLower(platform[0]) == "android" && strings.ToLower(osversion[0]) == "9" {
				res = jq.WhereNotEqual("description", "admob")
			}
		}
	}

	resGet := res.SortBy("value", "desc").Get()
	file, _ := json.MarshalIndent(resGet, "", " ")

	// Show admod-optout only when the result does not contain any admob's
	if strings.Contains(string(file), "\"description\": \"admob\"") {
		resGet = gojsonq.New().JSONString(string(file)).WhereNotEqual("description", "admod-optout").Get()
	}

	// If the response is empty, fill it with some data
	if string(file) == "[]" {
		jq2 := gojsonq.New().File(filePath)
		res2 := jq2.WhereNotNil("id")
		for param, value := range params {
			// Query all the data to a given platform, if platform parameter was given
			if strings.ToLower(param) == "platform" {
				res2 = jq2.WhereEqual(param, strings.ToLower(value[0]))
			}
		}
		resGet2 := res2.SortBy("value", "desc").Get()
		enc := json.NewEncoder(w)
		enc.SetIndent("", "	")
		enc.Encode(resGet2)
	} else {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "	")
		enc.Encode(resGet)
	}
}

// Create a new Ad Network
func createAdNetwork(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var adNetwork adNetwork
	reqBody = []byte(strings.ToLower(string(reqBody)))
	json.Unmarshal(reqBody, &adNetwork)
	jq := gojsonq.New().File(filePath)
	// Get current max id, for autoincrementing ID
	res := jq.Max("id")
	adNetwork.ID = int64(res) + 1
	adNetworks = append(adNetworks, adNetwork)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(adNetwork)
	file, _ := json.MarshalIndent(adNetworks, "", " ")
	_ = ioutil.WriteFile(filePath, file, 0644)
}

// Delete an existing Ad Network via ID
func deleteAdNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		fmt.Println("Query ERROR")
		return
	}

	for index, adNetwork := range adNetworks {
		if adNetwork.ID == id {
			adNetworks = append(adNetworks[:index], adNetworks[index+1:]...)

			enc := json.NewEncoder(w)
			enc.SetIndent("", "	")
			enc.Encode(adNetwork)
		}
	}

	file, _ := json.MarshalIndent(adNetworks, "", " ")
	_ = ioutil.WriteFile(filePath, file, 0644)
}

// Update Value on an existing Ad Network via ID
func updateAdNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		fmt.Println("Query ERROR")
		return
	}
	for index, adNetwork := range adNetworks {
		if adNetwork.ID == id {
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &adNetwork)
			adNetworks[index].Value = adNetwork.Value
			enc := json.NewEncoder(w)
			enc.SetIndent("", "	")
			enc.Encode(adNetwork)
		}
	}
	file, _ := json.MarshalIndent(adNetworks, "", " ")
	_ = ioutil.WriteFile(filePath, file, 0644)
}

// Define URL requests and run the app
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/adnetworks", returnAllAdNetworks)
	myRouter.HandleFunc("/adnetwork", createAdNetwork).Methods("POST")
	myRouter.HandleFunc("/adnetwork/{id}", deleteAdNetwork).Methods("DELETE")
	myRouter.HandleFunc("/adnetwork/{id}", updateAdNetwork).Methods("POST")
	myRouter.HandleFunc("/adnetwork/{adtype}", returnAdType)
	myRouter.HandleFunc("/adnetwork", queryAdNetworks)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

// MAIN, read json data from the output.txt file and add it to the global variable adNetworks
func main() {
	s, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	json.Unmarshal([]byte(s), &adNetworks)

	handleRequests()
}
