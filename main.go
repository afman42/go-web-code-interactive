package main

import (
	"bufio"
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
	"strings"

	"github.com/afman42/go-web-code-interactive/utils"
)

//go:embed web/dist
var WebContent embed.FS

var IpCors string

const (
	ModeDev     = "dev"
	ModeProd    = "prod"
	ModePreview = "preview"
)

func main() {
	Mode := ModeDev
	env := ".env.local"
	flag.Func("mode", "mode:dev,preview,prod", func(s string) error {
		if s == ModePreview {
			Mode = ModePreview
			env = ".env.preview"
		}
		if s == ModeProd {
			Mode = ModeProd
			env = ".env.prod"
		}
		return nil
	})
	flag.Parse()
	err := loadEnvFile(env)
	if err != nil {
		log.Fatal("Error loading " + env + " file")
	}

	IpCors = os.Getenv("CORS_DOMAIN")
	Port := os.Getenv("APP_PORT")
	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("tmp", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}

	files, err := filepath.Glob(utils.PathFileTemp("*"))
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
	handler := utils.WrapHandlerWithLogging(http.HandlerFunc(index))
	mux.Handle("/", handler)
	if Mode == ModePreview || Mode == ModeProd {
		mux.Handle("/assets/", http.FileServer(http.FS(dist)))
		// Static Folder web/public
		mux.Handle("/vite.svg", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := WebContent.ReadFile("web/dist/vite.svg")
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "image/svg+xml")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}))
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
	Language   string `json:"lang"`
	Tipe       string `json:"type"`
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
		var (
			data Data
			args string = "-"
		)
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data.Txt = strings.TrimSpace(data.Txt)
		data.Language = strings.TrimSpace(data.Language)
		data.Tipe = strings.TrimSpace(data.Tipe)
		if data.Language == "" || data.Txt == "" || data.Tipe == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				StatusCode int    `json:"statusCode"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Something Went Wrong, Check Field Empty",
			})
			return
		}
		if !utils.CheckIsNotData([]string{"php", "node", "go"}, data.Language) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				StatusCode int    `json:"statusCode"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Something Went Wrong, Check Language is exist",
			})
			return
		}
		if !utils.CheckIsNotData([]string{"repl", "stq"}, data.Tipe) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				StatusCode int    `json:"statusCode"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Something Went Wrong, Check Type is exist",
			})
			return
		}
		filename := "index-" + utils.StringWithCharset(5) + ".js"
		if data.Language == "php" {
			filename = "index-" + utils.StringWithCharset(5) + ".php"
		}
		if data.Language == "go" {
			args = "run"
			filename = "main-" + utils.StringWithCharset(5) + ".go"
		}
		if data.Tipe == "stq" {
			if data.Language == "go" {
				data.Txt = data.Txt + utils.TxtGo
			}
			if data.Language == "php" {
				data.Txt = data.Txt + utils.TxtPHP
			}
			if data.Language == "node" {
				data.Txt = data.Txt + utils.TxtJS
			}
		}
		err = os.WriteFile(filename, []byte(data.Txt), 0755)
		if err != nil {
			fmt.Printf("unable to write file: %w", err)
		}
		err = utils.MoveFile(filename, utils.PathFileTemp(filename))
		if err != nil {
			fmt.Println("error movefile: ", err)
		}

		out, errout, err := utils.Shellout(data.Language, utils.PathFileTemp(filename))
		if data.Language == "go" {
			out, errout, err = utils.Shellout(data.Language, args, utils.PathFileTemp(filename))
		}
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

// ref: https://grok.com/share/bGVnYWN5_6e37d0ea-cc9a-4ab2-8cad-ad24fe59e00d
func loadEnvFile(filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // Remove leading/trailing whitespace

		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable
		os.Setenv(key, value)
	}

	// Check for scanning errors
	return scanner.Err()
}
