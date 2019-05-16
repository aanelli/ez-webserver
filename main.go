package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"io/ioutil"
	"errors"
	"mime/multipart"
)

type SubmissionInfo struct {
	Term string
}

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))
	mux  := http.NewServeMux()
	templateData := SubmissionInfo{
		Term: "Fall 2019",
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.Execute(w, templateData)
    })

	mux.HandleFunc("/upload/pdf", func(w http.ResponseWriter, r *http.Request) {
		err := UploadFile(w, r)
		if err != nil {
			log.Fatalln("err handling file upload: %v", err)
		}
    })

	fmt.Println("Listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}


func UploadFile(w http.ResponseWriter, r *http.Request) error {
	//if the request isn't a post request tell them to shoo
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        err := errors.New("request isn't post request")
		fmt.Printf("request isn't post request: %v", err)
        return err
    }

	//grab the first file in the form
    file, handle, err := r.FormFile("file")
    if err != nil {
		fmt.Printf("error grabbing the file from form: %v", err)
        return err
    }
    defer file.Close()

	//make sure it's a supported file type, I guess? (not sure if necessary)
    mimeType := handle.Header.Get("Content-Type")
    switch mimeType {
    case "application/pdf":
        err = saveFile(w, file, handle)
    default:
        jsonResponse(w, http.StatusBadRequest, "The format of the uploaded file is not valid.(*.pdf)")
    }

	if err != nil{
		fmt.FPrintf("error saving file: %v", err)
		return err
	}

	return nil
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) error {
    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Fprintf(w, "%v", err)
        return err
    }

	//write to disk just for testing purposes
    err = ioutil.WriteFile("./files/"+handle.Filename, data, 0666)
    if err != nil {
        fmt.Fprintf(w, "%v", err)
        return err
    }
    jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
	return nil
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    fmt.Fprint(w, message)
}
