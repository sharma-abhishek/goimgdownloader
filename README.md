# goimgdownloader

This project is created as an exercise to explore concept of bufferedchannel in golang. 

#### Supported Endpoints
This program has two endpoints:
 - GET API to request downloading n number of images. (This API schedules a worker and quickly returns the RequestID for tracking)
```
/work?n=<number_of_images>
```

- GET API to know the request status (This returns JSON containing fields with status such as queued, downloading, downloaded and failed)
```
/status?requestid=<requestid>
```

#### Building and Running

- Clone this repository on your local machine
- This program has one dependency for generating UUID (i.e used as request id). Install the following dependency
```
$ go get github.com/satori/go.uuid
```
- Once dependency is installed, use following command to build:
```
$ go build -o goimgdownloader .
```
- After build is successful, use below command to run:
```
$ ./goimgdownloader
```

By default, it runs on port 8090 with 5 workers. However, these settings are configurable.
