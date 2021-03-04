package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	//globalDebug ... is a bool, defined by the return of the nameless function, invoked with () at the end.
	globalDebug = true
	httpClient  = http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}
)

//dp ... stands for Debug Printing: Will print to terminal only if the var globalDebug is true.
func dp(str string) {
	if globalDebug {
		fmt.Println(str)
	}
}

func pause(w http.ResponseWriter, r *http.Request) {
	var portInterface interface{}
	fullBody, err := ioutil.ReadAll(r.Body)
	shouldPanic(err)
	err = json.Unmarshal(fullBody, &portInterface)
	containerPort := portInterface.(map[string]string)["port"]
	fmt.Println(containerPort)
}

func getContainerName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Container Name Called")
	myContainerName := os.Getenv("HOSTNAME")
	fmt.Printf("My Container Name is: %s\n", myContainerName)
	response := fmt.Sprintf(`
		{containername: %s}
	`, myContainerName)
	_, err := w.Write([]byte(response))
	shouldPanic(err)
}

func shouldPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func touchServerByContainerName(w http.ResponseWriter, r *http.Request) {
	fullBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	containerNameReceived := stripContainerName(fullBody)
	fmt.Printf("Calling get function on container: %s\n", containerNameReceived)
	myContainerNameJSON := fmt.Sprintf(`
	{"containername":"%s"}
	`, containerNameReceived)
	req, err := http.NewRequest(http.MethodGet, "http://"+containerNameReceived+":9001/api/v1/getcontainername", strings.NewReader(myContainerNameJSON))
	shouldPanic(err)
	res, err := httpClient.Do(req)
	shouldPanic(err)
	resBytes, err := ioutil.ReadAll(res.Body)
	shouldPanic(err)
	fmt.Printf("\n%s\n", resBytes)
	defer res.Body.Close()

}

func stripContainerName(fullBody []byte) string {
	var myJSONInterface interface{}
	err := json.Unmarshal(fullBody, &myJSONInterface)
	shouldPanic(err)
	containername := myJSONInterface.(map[string]interface{})["containername"].(string)
	fmt.Printf("ContainerName: %v\n", containername)
	return containername
}

func printJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Print JSON called.")
	fullBody, err := ioutil.ReadAll(r.Body)
	shouldPanic(err)
	var myJSONInterface interface{}
	err = json.Unmarshal(fullBody, &myJSONInterface)
	fmt.Printf("\n%v\n", myJSONInterface)
	containername := myJSONInterface.(map[string]interface{})["containername"].(string)
	fmt.Printf("ContainerName: %v\n", containername)
}

func testServ(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Online")
}

//*********************************************** API HANDLERS *******************************************

func main() {
	dp("Server Started")
	r := mux.NewRouter()
	myName := os.Getenv("HOSTNAME")
	fmt.Println("My name is: " + myName)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/pause", pause).Methods(http.MethodPost)
	api.HandleFunc("/getcontainername", getContainerName).Methods(http.MethodGet)
	api.HandleFunc("/touchserver", touchServerByContainerName).Methods(http.MethodGet)
	api.HandleFunc("", testServ).Methods(http.MethodGet)
	api.HandleFunc("/printjson", printJSON).Methods(http.MethodGet)

	//api.HandleFunc("/batchcollect", batchCollect).Methods(http.MethodPost)
	/* 	api.HandleFunc("/started", collectionStarted).Methods(http.MethodPost)
	   	api.HandleFunc("/isalive", collectionIsAlive).Methods(http.MethodPost)
	   	api.HandleFunc("/finished", collectionFinished).Methods(http.MethodPost) */

	//api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe("0.0.0.0:9001", r))
}
