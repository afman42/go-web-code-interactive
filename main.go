package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("tmp", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	fmt.Println("Server starting in localhost:3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

// https://www.alexedwards.net/blog/which-go-router-should-i-use
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var tmpl, err = template.ParseFiles("./index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]string{"out": "", "errout": ""}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		var tmpl, err = template.ParseFiles("./index.html")
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var txt = r.FormValue("txt")
		filename := "index-" + StringWithCharset(5) + ".js"
		err = os.WriteFile(filename, []byte(txt), 0755)
		if err != nil {
			fmt.Printf("unable to write file: %w", err)
		}
		err = MoveFile(filename, "./tmp/"+filename)
		if err != nil {
			fmt.Println("error movefile: ", err)
		}
		out, errout, err := Shellout("./node/node.exe", ".\\tmp\\"+filename)
		if err != nil {
			log.Printf("error shell: %v\n", err)
		}
		fmt.Println("--- stdout ---")
		fmt.Println(out)
		fmt.Println("--- stderr ---")
		fmt.Println(errout)
		var data = map[string]string{"out": out, "errout": errout}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
