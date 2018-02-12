package main

type RequestIdResponse struct {
	RequestID string
}

type WorkRequestStatus struct {
	Downloaded  int `json:"downloaded"`
	Downloading int `json:"downloading"`
	Queued      int `json:"queued"`
	Failed      int `json:"failed"`
}

type WorkRequest struct {
	ID                       string
	Call                     func(w WorkRequest, worker Worker) error
	NumberOfImagesToDownload int
	status                   WorkRequestStatus
}

type Worker struct {
	ID         int
	Work       chan WorkRequest
	WorkerPool chan chan WorkRequest
	QuitChan   chan bool
}
