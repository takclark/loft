package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultRootDir = "/Users/tclark/Personal/code/go/src/github.com/takclark/loft"
	defaultPort    = 8080
)

type PhotoStreamServer struct {
	Port    int
	RootDir string
}

type ImageListResponse struct {
	Images []*ImageFileInfo `json:"images"`
}

type ImageFileInfo struct {
	CreatedAt string `json:"created_at"`
	Filename  string `json:"filename"`
}

func NewPhotoStreamServer() *PhotoStreamServer {
	return &PhotoStreamServer{
		RootDir: defaultRootDir,
		Port:    defaultPort,
	}
}

func (s *PhotoStreamServer) UploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UploadHandler with method:", r.Method)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, mpfh, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "problem reading upload\n")
		return
	}

	detectedContentType := http.DetectContentType(buff)
	fmt.Printf("received upload and detected content type: %v\n", detectedContentType)
	if !strings.HasPrefix(detectedContentType, "image") {
		// not an image uplaod
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "that doesn't look like an image")
		return
	}

	// reset reader
	file.Seek(0, 0)

	defer file.Close()
	extension := filepath.Ext(mpfh.Filename)
	newFilename := fmt.Sprintf("%v%s", time.Now().UnixNano(), extension)
	f, err := os.OpenFile("./assets/images/"+newFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	fmt.Printf("received upload and wrote to assets/images/%s\n", newFilename)
	fmt.Fprintf(w, "received upload. thanks!")
}

// ListHandler returns a JSON-formatted list of files in the image directory
func (s *PhotoStreamServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received list request")

	imageDir := s.RootDir + "/assets/images"
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imageListResponse := &ImageListResponse{Images: make([]*ImageFileInfo, len(files))}
	for i, f := range files {
		info := &ImageFileInfo{
			Filename:  f.Name(),
			CreatedAt: fmt.Sprintf("%v", f.ModTime()),
		}
		imageListResponse.Images[i] = info
	}

	res, err := json.Marshal(imageListResponse)
	if err != nil {
		fmt.Println("problem marshalling json list response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func (s *PhotoStreamServer) Start() {
	http.HandleFunc("/upload", s.UploadHandler)
	http.HandleFunc("/list", s.ListHandler)
	http.Handle("/", http.FileServer(http.Dir(s.RootDir+"/assets")))
	addr := fmt.Sprintf(":%d", s.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
