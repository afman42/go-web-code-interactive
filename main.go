package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/afman42/go-web-code-interactive/utils"
)

func main() {
	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("tmp", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}

	files, err := filepath.Glob(utils.PathFileTemp("*.js"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	fmt.Println("Server starting in localhost:8000")
	err = http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}

type Data struct {
	Txt        string `json:"txt"`
	Stdout     string `json:"out"`
	Stderr     string `json:"errout"`
	StatusCode int    `json:"statusCode"`
}

// https://www.alexedwards.net/blog/which-go-router-should-i-use
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodGet:
		// var tmpl, err = template.ParseFiles("./views/index.html")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		//
		// var data = map[string]string{"out": "", "errout": ""}
		// err = tmpl.Execute(w, data)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		return

	case http.MethodPost:
		var data Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		filename := "index-" + utils.StringWithCharset(5) + ".js"
		err = os.WriteFile(filename, []byte(data.Txt), 0755)
		if err != nil {
			fmt.Printf("unable to write file: %w", err)
		}
		err = utils.MoveFile(filename, utils.PathFileTemp(filename))
		if err != nil {
			fmt.Println("error movefile: ", err)
		}
		out, errout, err := utils.Shellout("node", utils.PathFileTemp(filename))
		if err != nil {
			log.Printf("error shell: %v\n", err)
		}
		fmt.Println("--- stdout ---")
		fmt.Println(out)
		fmt.Println("--- stderr ---")
		fmt.Println(errout)
		data.Stdout = out
		data.Stderr = errout
		data.StatusCode = http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		http.StatusText(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
