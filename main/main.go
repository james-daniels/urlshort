package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/james-daniels/gophersizes/urlshort"
)

var yamlFile string
var jsonFile string

// Build the MapHandler using the mux as the fallback
var pathsToUrls = map[string]string{}

func init(){
	flag.StringVar(&yamlFile, "yaml", "", "read a yaml file?")
	flag.StringVar(&jsonFile, "json", "", "read a json file?")
}

func main() {
	flag.Parse()

	mux := defaultMux()
// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if yamlFile == "" && jsonFile == "" {
		fmt.Println("need to use either a yaml or json file")
		os.Exit(0)
	} else if yamlFile != "" && jsonFile != "" {
		fmt.Println("there can only be one!")
		os.Exit(0)
	}

	if jsonFile != "" {
		execJSON(mapHandler)
	} else if yamlFile != "" {
		execYAML(mapHandler)
	}

}


func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func parseFile(file string )[]byte {
	fbyte, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return fbyte
}


func execJSON(mapHandler http.HandlerFunc){
	jsonInput := parseFile(jsonFile)
	jsonHandler, err := urlshort.JSONHandler(jsonInput, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func execYAML(mapHandler http.HandlerFunc){
	yamlInput := parseFile(yamlFile)
	yamlHandler, err := urlshort.YAMLHandler(yamlInput, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}
