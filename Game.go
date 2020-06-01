package main

//import (
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)
//
//type Game struct {
//	int points
//
//}
//
//type HTTPResponseErrorBundle struct {
//	response *http.Response
//	err error
//}
//
//func getCluesAsync(url string) <-chan *HTTPResponseErrorBundle{
//	ch := make(chan *HTTPResponseErrorBundle)
//
//	go func() {
//		defer close(ch)
//		resp, err := http.Get(url)
//		if err != nil{
//			fmt.Println("Error Fetching Clues...")
//		}
//		currentResponseBundle := &HTTPResponseErrorBundle{resp, err}
//		ch <- currentResponseBundle
//	}()
//
//	return ch
//}
//
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