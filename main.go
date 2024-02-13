package main

import (
	"bytes"
	"mime/multipart"
	"fmt"
	"log"
	"html/template"
	"net/http"
	"os"
	"io"
	
	"github.com/joho/godotenv"
)


func init() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}

func main() {
	http.HandleFunc("/", upload)
	fmt.Println("Starting Server")
	http.ListenAndServe(":5000",nil)
}


var tmpl = template.Must(template.ParseGlob("index.html"))
// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w,"index.html",nil)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
		fmt.Println(err)
		return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		filename := "./test/"+handler.Filename
		fmt.Println(filename)
		target_url := os.Getenv("TARGET_URL")
		fmt.Println(target_url)
		postFile(filename, target_url)
		os.Remove("./test/"+handler.Filename)
       }
}


func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    // this step is very important
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    // open file handle
    fh, err := os.Open(filename)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }
    defer fh.Close()

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}

//~ func sendImageAsHTML(w http.ResponseWriter, r *http.Request, a string) {
	//~ fmt.Fprintf(w,a)
//~ }

//~ func sendImageAsAttachment(w http.ResponseWriter, r *http.Request) {
    //~ buf, err := os.ReadFile("F46ZPQ0bQAACFYs.jpg")
    //~ if err != nil {
        //~ log.Fatal(err)
    //~ }
    //~ w.Header().Set("Content-Type", "image/jpeg")
    //~ w.Header().Set("Content-Disposition", `attachment;filename="F46ZPQ0bQAACFYs.jpg"`)
    //~ w.Write(buf)
//~ }

//~ func sendImageAsBytes(w http.ResponseWriter, r *http.Request, a string) {
    //~ buf, err := os.ReadFile("./uploads/"+a)
    //~ if err != nil {
        //~ log.Fatal(err)
    //~ }
    //~ w.Header().Set("Content-Type", "image/jpeg")
    //~ w.Write(buf)
//~ }


//~ func renderError(w http.ResponseWriter, message string, statusCode int) {
	//~ w.WriteHeader(statusCode)
	//~ w.Write([]byte(message))
//~ }
