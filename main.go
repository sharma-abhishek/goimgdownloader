package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

var (
	NumOfWorkers = flag.Int("n", 5, "number of workers to start")
	HTTPAddress  = flag.String("http", "127.0.0.1:8090", "host and port to listen on")
)

//WorkQueue : buffered channel of work request
var WorkQueue = make(chan WorkRequest, 100)

// Each request status is tracked in a map. 
var requestTrackMap = make(map[string]WorkRequestStatus)

// Function to download an image with random name (uuid)
func downloadImage(url string) error {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("Trouble making GET photo request!")
		return err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Trouble reading response body!")
		return err
	}

	filename := uuid.Must(uuid.NewV4()).String() + ".jpg"

	err = ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		log.Fatal("Trouble creating file! -- ", err)
		return err
	}
	return nil
}

// function that takes work request and worker and download number of images.
// (picsum.photos service is used as it returns random image each time)
func downloadImages(w WorkRequest, worker Worker) error {
	fmt.Println("processing request")
	n := w.NumberOfImagesToDownload
	status := w.status
	status.Queued = n
	requestTrackMap[w.ID] = status
	for i := 0; i < n; i++ {
		status.Queued = status.Queued - 1
		status.Downloading = status.Downloading + 1
		downloadError := downloadImage("https://picsum.photos/200/300/?random")
		status.Downloading = status.Downloading - 1
		if downloadError == nil {
			status.Downloaded = status.Downloaded + 1
		} else {
			status.Failed = status.Failed + 1
		}
		requestTrackMap[w.ID] = status
	}
	worker.Stop()
	return nil
}

// Handler function to get status by providing requestid
func ProcessRequestStatus(w http.ResponseWriter, r *http.Request) {
	requestid := r.URL.Query().Get("requestid")
	if val, ok := requestTrackMap[requestid]; ok {
		
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		
		responsejson, _ := json.Marshal(val)
		w.Write([]byte(responsejson))
	} else {
		http.Error(w, "Request id not found", http.StatusNotFound)
		return
	}

}

// Handler function to accept Work Request to download n images
func ProcessRequest(w http.ResponseWriter, r *http.Request) {

	numberOfImageToDownload := r.URL.Query().Get("n")

	if numberOfImageToDownload == "" {
		http.Error(w, "You must specify number of images to download.", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(numberOfImageToDownload)

	if err != nil {
		http.Error(w, "Query param {n} must be integer", http.StatusBadRequest)
		return
	}

	workStatus := WorkRequestStatus{
		Downloaded:  0,
		Downloading: 0,
		Queued:      0,
		Failed:      0,
	}

	requestID := uuid.Must(uuid.NewV4()).String()

	work := WorkRequest{
		ID:   requestID,
		Call: downloadImages,
		NumberOfImagesToDownload: i,
		status: workStatus,
	}
	requestTrackMap[requestID] = workStatus

	WorkQueue <- work
	fmt.Println("Work request queued")

	// let your user know that their work request is created and send requestid for tracking
	response := RequestIdResponse{
		requestID,
	}
	responsejson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responsejson)
}

func main() {

	flag.Parse()

	//starting number of workers
	StartDispatcher(*NumOfWorkers)

	fmt.Println("Registering the endpoints and associating handlers")
	http.HandleFunc("/work", ProcessRequest)
	http.HandleFunc("/status", ProcessRequestStatus)

	// Start the HTTP server
	fmt.Println("HTTP server listening on", *HTTPAddress)
	if err := http.ListenAndServe(*HTTPAddress, nil); err != nil {
		fmt.Println(err.Error())
	}

}
