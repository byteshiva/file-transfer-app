package main

import (
    "fmt"
    "html/template"
    "io"
    "log"
    "net/http"
    "os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// render the index page template
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, err := template.ParseFiles("index.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
    } else if r.Method == "POST" {
        // parse input, type multipart/form-data
        err := r.ParseMultipartForm(32 << 20) // limit your max input length!
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // retrieve the file from form data
        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        // create a new file in the uploads directory
        err = os.MkdirAll("./uploads", os.ModePerm)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        dst, err := os.Create("./uploads/" + header.Filename)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        // copy the uploaded file to the destination file
        if _, err := io.Copy(dst, file); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "File uploaded successfully!")
    } else {
        http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
    }
}

func main() {
	http.HandleFunc("/", indexHandler)

    http.HandleFunc("/upload", uploadHandler)
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
