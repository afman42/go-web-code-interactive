package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/afman42/go-web-code-interactive/utils"
)

//go:embed web/dist
var WebContent embed.FS

const (
	ModeDev     = "dev"
	ModeProd    = "prod"
	ModePreview = "preview"
)

var (
	IpCors string
	Port   string
	Mode   string
)

func main() {
	flag.StringVar(&IpCors, "ip-cors", "http://localhost:5173", "ip or domain")
	flag.StringVar(&Port, "port", "8000", "server port")
	flag.Func("mode", "mode:dev,preview,prod", func(s string) error {
		Mode = ModeDev
		if s == ModeProd {
			Mode = s
		}
		if s == ModePreview {
			Mode = s
			IpCors = "http://localhost:" + Port
		}
		return nil
	})
	flag.Parse()
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
	dist, err := fs.Sub(WebContent, "web/dist")
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	if Mode == ModePreview {
		mux.Handle("/assets/", http.FileServer(http.FS(dist)))
	}
	fmt.Println("Server starting in localhost:" + Port)
	err = http.ListenAndServe(":"+Port, mux)
	if err != nil {
		log.Fatal("Something went wrong", err)
		os.Exit(1)
	}
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
	w.Header().Set("Access-Control-Allow-Origin", IpCors)
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodGet:
		var tmp, err = template.ParseFS(WebContent, "web/dist/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "text/html")
		err = tmp.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		w.WriteHeader(http.StatusOK)
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
