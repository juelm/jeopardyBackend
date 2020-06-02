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

func getClues(categories []CategoryHeader) [] map[string] interface{}{
	clues := [] map[string] interface{}{}
	for _, category := range categories {
		clueChannel := <- asyncGetClueData(category.Id)
		if clueChannel.err != nil{
			panic(clueChannel.err)
		}
		clueResponse := clueChannel.response.Body
		defer clueResponse.Close()
		clueJSON, ioerr := ioutil.ReadAll(clueResponse)
		if ioerr != nil{
			panic(ioerr)
		}
		newClue := unMarshalClues(clueJSON)
		clues = append(clues, newClue)
	}
	return clues
}

func asyncGetClueData(id int) <- chan *HTTPResponseErrorBundle{
	ch := make(chan *HTTPResponseErrorBundle)
	go func() {
		resp, err := http.Get(fmt.Sprintf("http://jservice.io/api/category?id=%d", id))
		if err != nil{
			panic(err)
		}
		currentResponseBundle := HTTPResponseErrorBundle{resp, err}
		ch <- &currentResponseBundle
	}()
	return ch
}

func unMarshalClues(data []byte) map[string] interface{} {
	clue := make(map[string] interface{})
	err := json.Unmarshal(data, &clue)
	if err != nil {
		panic(err)
	}
	return clue
}

func GetNewBoard(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	currentCategories := GetCategories()
	currentClues := getClues(currentCategories)
	prettyClues, err := json.MarshalIndent(currentClues,"","   ")
	if err != nil {
		panic(err)
	}
	fmt.Println(currentCategories)
	fmt.Println("\n\n")
	fmt.Printf("%s\n", string(prettyClues))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentClues)
	//fmt.Fprintf(w, string(currentCategories))

}

func enableCORS(w *http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func marshalJSON(data []map[string] interface{}) []byte {
	jsonClues, err := json.Marshal(map[string] interface{}{})
	if err != nil{
		panic(err)
	}
	return jsonClues
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/newgame", GetNewBoard)
	//http.HandleFunc("/getcategories", GetCategories)
	fmt.Println("Server Starting...")
	http.ListenAndServe(":8080", nil)
}
