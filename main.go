package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

//go:embed index.tmpl
var indexHTML embed.FS

func main() {

	port := getPort()

	startingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Impossible de récupérer le repertoire courant")
		return
	}

	fmt.Printf("Start serve : %v\n", startingDirectory)

	// Download
	http.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {

		fname := r.URL.Query().Get("name")
		if fname == "" {
			http.Error(w, "Aucun nom fourni", 500)
		}

		file, err := os.Open(fname)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	})

	// Index Of
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(indexHTML, "index.tmpl"))
		files, err := ReadDirectory(startingDirectory)
		if err != nil {
			http.Error(w, "Impossible de lister les fichiers", 500)
		}

		buffer := new(bytes.Buffer)
		tmpl.Execute(buffer, files)
		buffer.WriteTo(w)
	})

	addr := "0.0.0.0:" + port
	fmt.Println("(http://localhost:"+port+") - Serve on port ", port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getPort() string {
	var port string
	if len(os.Args) < 2 {
		port = "8080"
	} else {
		port = os.Args[1]
		if port == "" {
			port = "8080"
		}
	}
	return port
}

func ReadDirectory(dir string) ([]os.FileInfo, error) {

	f, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, v := range files {
		fmt.Printf("-- %v : %v \n", v.Name(), v.IsDir())
	}

	return files, nil
}
