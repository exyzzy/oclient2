package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/exyzzy/oauth2"
	"github.com/gorilla/mux"
)

func main() {
	err := oauth2.InitOauth2("services.json")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	r.HandleFunc("/", PageHomeHandler)
	r.HandleFunc("/page/api", PageApiHandler)
	r.HandleFunc("/login/{authtype}/{service}", LoginHandler)
	r.HandleFunc("/authlink/{authtype}/{service}", AuthlinkHandler)
	r.HandleFunc("/redirect", RedirectHandler)
	r.HandleFunc("/strava/get/athlete", StravaGetAthleteHandler)
	r.HandleFunc("/strava/get/activities", StravaGetActivitiesHandler)
	r.HandleFunc("/linkedin/get/me", LinkedinGetMeHandler)
	r.HandleFunc("/spotify/get/me", SpotifyGetMeHandler)
	r.HandleFunc("/spotify/get/newreleases", SpotifyGetNewReleasesHandler)
	r.HandleFunc("/spotify/put/rename", SpotifyPutRenameHandler)
	r.HandleFunc("/github/get/user", GithubGetUserHandler)
	r.HandleFunc("/fitbit/get/user", FitbitGetUserHandler)
	r.HandleFunc("/fitbit/get/heartrate", FitbitGetHeartrateHandler)
	r.HandleFunc("/fitbit/get/sleep", FitbitGetSleepHandler)
	r.HandleFunc("/oura/get/user", OuraGetUserHandler)
	r.HandleFunc("/oura/get/sleep", OuraGetSleepHandler)
	r.HandleFunc("/oura/get/activity", OuraGetActivityHandler)
	r.HandleFunc("/oura/get/readiness", OuraGetReadinessHandler)
	r.HandleFunc("/google/get/user", GoogleGetUserHandler)
	r.HandleFunc("/amazon/get/user", AmazonGetUserHandler)
	r.HandleFunc("/withings/post/sleep", WithingsPostSleepHandler)
	http.Handle("/", r)
	fmt.Println(">>>>>>> OClient started at:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	return
}

func PageHomeHandler(w http.ResponseWriter, r *http.Request) {
	pageHandler(w, r, nil, "templates", "home.html")
}

func PageApiHandler(w http.ResponseWriter, r *http.Request) {
	pageHandler(w, r, nil, "templates", "api.html")
}

func pageHandler(w http.ResponseWriter, r *http.Request, data interface{}, dir string, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, path.Join(dir, file))
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authtype := vars["authtype"]
	service := vars["service"]
	authlink := oauth2.AuthLink(r, authtype, service)
	http.Redirect(w, r, authlink, http.StatusTemporaryRedirect)
}

func AuthlinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authtype := vars["authtype"]
	service := vars["service"]
	authlink := oauth2.AuthLink(r, authtype, service)
	fmt.Fprintln(w, authlink)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("Redirect Error: ", r.URL.RawQuery, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	code := m.Get("code")
	state := m.Get("state")
	err = oauth2.ExchangeCode(w, r, code, state) //do not write to w before this call
	if err != nil {
		http.Error(w, "Exchange Failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintln(w, "Code: ", code, " Scope: ", scope)
	http.Redirect(w, r, "/page/api", http.StatusTemporaryRedirect)
}

func processAPI(w http.ResponseWriter, r *http.Request, service string, action string, url string, data map[string]interface{}) (result string, err error) {
	resp, err := oauth2.ApiRequest(w, r, service, action, url, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result = string(body)
	return
}

//== API Examples

func StravaGetAthleteHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://www.strava.com/api/v3/athlete"
	result, err := processAPI(w, r, oauth2.STRAVA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func StravaGetActivitiesHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://www.strava.com/api/v3/athlete/activities?page=1&per_page=30"
	result, err := processAPI(w, r, oauth2.STRAVA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func LinkedinGetMeHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.linkedin.com/v2/me"
	result, err := processAPI(w, r, oauth2.LINKEDIN, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func SpotifyGetMeHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.spotify.com/v1/me"
	result, err := processAPI(w, r, oauth2.SPOTIFY, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func SpotifyGetNewReleasesHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.spotify.com/v1/browse/new-releases"
	result, err := processAPI(w, r, oauth2.SPOTIFY, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func SpotifyPutRenameHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"name":        "Updated Playlist Name",
		"description": "Updated playlist description",
		"public":      false,
	}

	url := "https://api.spotify.com/v1/playlists/2RmnrZSPoYtVyjou7DU8We"
	result, err := processAPI(w, r, oauth2.SPOTIFY, "PUT", url, data)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func GithubGetUserHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.github.com/user"
	result, err := processAPI(w, r, oauth2.GITHUB, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func FitbitGetUserHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.fitbit.com/1/user/-/profile.json"
	result, err := processAPI(w, r, oauth2.FITBIT, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}

}

func FitbitGetHeartrateHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.fitbit.com/1/user/-/activities/heart/date/today/1d/1sec.json"
	result, err := processAPI(w, r, oauth2.FITBIT, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func FitbitGetSleepHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.fitbit.com/1.2/user/-/sleep/date/2021-08-08.json?timezone=UTC"
	result, err := processAPI(w, r, oauth2.FITBIT, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func OuraGetUserHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.ouraring.com/v1/userinfo"
	result, err := processAPI(w, r, oauth2.OURA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func OuraGetSleepHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.ouraring.com/v1/sleep"
	result, err := processAPI(w, r, oauth2.OURA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func OuraGetActivityHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.ouraring.com/v1/activity"
	result, err := processAPI(w, r, oauth2.OURA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func OuraGetReadinessHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.ouraring.com/v1/readiness"
	result, err := processAPI(w, r, oauth2.OURA, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func GoogleGetUserHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://www.googleapis.com/oauth2/v3/userinfo"
	result, err := processAPI(w, r, oauth2.GOOGLE, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func AmazonGetUserHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://api.amazon.com/user/profile"
	result, err := processAPI(w, r, oauth2.AMAZON, "GET", url, nil)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}

func WithingsPostSleepHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"action":     "getsummary",
		"lastupdate": 1,
	}
	url := "https://wbsapi.withings.net/v2/sleep"
	result, err := processAPI(w, r, "withings", "POST", url, data)
	if err == nil {
		fmt.Fprintln(w, result)
	}
}
