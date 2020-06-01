package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type HTTPResponseErrorBundle struct {
	response *http.Response
	err error
}

type CategoryHeader struct {
	Id int
	Title string
	Count int
}

func getCluesAsync(url string) <-chan *HTTPResponseErrorBundle{
	ch := make(chan *HTTPResponseErrorBundle)

	go func() {
		defer close(ch)
		resp, err := http.Get(url)
		if err != nil{
			fmt.Println("Error Fetching Clues...")
		}
		currentResponseBundle := &HTTPResponseErrorBundle{resp, err}
		ch <- currentResponseBundle
	}()

	return ch
}

//func index(w http.ResponseWriter, r *http.Request) {
//	responseBundle := <- getCluesAsync("http://jservice.io/api/random")
//	if responseBundle.err != nil{
//		fmt.Println("Error Fetching Clues...")
//	}
//	w.WriteHeader(http.StatusOK)
//	resp := responseBundle.response.Body
//	defer responseBundle.response.Body.Close()
//	respString, ioerr := ioutil.ReadAll(resp)
//	if ioerr != nil{
//		fmt.Println("Error Fetching extracting body...")
//	}
//	fmt.Println(respString)
//	fmt.Fprintf(w, "Check your console")
//	//fmt.Fprintf(w, string(respString))
//}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Call an endpoint to see the jeopardy magic")
}

func GetRandomClue(w http.ResponseWriter, r *http.Request) {
	responseBundle := <- getCluesAsync("http://jservice.io/api/random")
	if responseBundle.err != nil{
		fmt.Println("Error Fetching Clues...")
	}
	w.WriteHeader(http.StatusOK)
	resp := responseBundle.response.Body
	defer responseBundle.response.Body.Close()
	respString, ioerr := ioutil.ReadAll(resp)
	if ioerr != nil{
		fmt.Println("Error Fetching extracting body...")
	}
	fmt.Println(respString)
	fmt.Fprintf(w, "Check your console")
	//fmt.Fprintf(w, string(respString))
}

func GetCategories() []CategoryHeader {
	randomOffset := rand.Int31n(300) + 1
	responseBundle := <- getCluesAsync(fmt.Sprintf("http://jservice.io/api/categories?count=6&offset=%d", randomOffset))
	if responseBundle.err != nil{
		fmt.Println("Error Fetching Clues...")
	}
	resp := responseBundle.response.Body
	defer responseBundle.response.Body.Close()
	respString, ioerr := ioutil.ReadAll(resp)
	if ioerr != nil{
		fmt.Println("Error extracting body...")
	}
	categories := unMarshalCategories(respString)
	return categories
}

func unMarshalCategories(data []byte) []CategoryHeader{
	CurrentCategories := make([]CategoryHeader, 0)
	err := json.Unmarshal(data, &CurrentCategories)
	if err != nil{
		log.Fatal(err)
	}
	return CurrentCategories
}

func GetNewBoard(w http.ResponseWriter, r *http.Request) {
	currentCategories := GetCategories()
	fmt.Println(currentCategories)
}

//func GetNewBoard(w http.ResponseWriter, r *http.Request) {
//	responseBundle := <- getCluesAsync("http://jservice.io/api/random")
//	if responseBundle.err != nil{
//		fmt.Println("Error Fetching Clues...")
//	}
//	w.WriteHeader(http.StatusOK)
//	resp := responseBundle.response.Body
//	defer responseBundle.response.Body.Close()
//	respString, ioerr := ioutil.ReadAll(resp)
//	if ioerr != nil{
//		fmt.Println("Error Fetching extracting body...")
//	}
//	fmt.Println(respString)
//	fmt.Fprintf(w, "Check your console")
//	//fmt.Fprintf(w, string(respString))
//}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/newgame", GetNewBoard)
	//http.HandleFunc("/getcategories", GetCategories)
	fmt.Println("Server Starting...")
	http.ListenAndServe(":8080", nil)
}
