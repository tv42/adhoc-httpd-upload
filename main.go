package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/tv42/adhoc-httpd-upload/internal"
)

var (
	host = flag.String("host", "0.0.0.0", "IP address to bind to")
	port = flag.Int("port", 8000, "TCP port to listen on")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [OPTS] DIR\n", os.Args[0])
	flag.PrintDefaults()
}

var page = template.Must(template.New("top").Parse(`
<html>
  <head>
    <title>Ad hoc file upload</title>
  </head>
  <body>
    <p>{{.}}</p>
    <p>Upload a file:</p>
    <form action="/" method="POST" enctype="multipart/form-data">
      <input type="file" name="f">
      <input type="submit" value="Upload">
    </form>
  </body>
</html>
`))

type uploadDir string

func (dir uploadDir) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		page.Execute(w, nil)
		return
	}
	f, hdr, err := req.FormFile("f")
	if err != nil {
		http.Error(w, "Need a file to upload.", 500)
		return
	}
	defer f.Close()

	if !internal.IsSafeName(hdr.Filename) {
		http.Error(w, "Unsafe filename.", 400)
		return
	}

	tmp, err := ioutil.TempFile(string(dir), ".tmp-")
	if err != nil {
		http.Error(w, "Cannot create temp file.", 500)
		return
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	_, err = io.Copy(tmp, f)
	if err != nil {
		http.Error(w, "Cannot write to temp file.", 500)
		return
	}

	err = os.Link(tmp.Name(), filepath.Join(string(dir), hdr.Filename))
	if err != nil {
		http.Error(w, "Cannot save file.", 500)
		return
	}

	log.Printf("Saved %q", hdr.Filename)
	page.Execute(w, "Thanks!")
}

func main() {
	prog := path.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	dir := flag.Arg(0)

	log.Printf("Receiving uploads to %q at http://%s:%d/", dir, *host, *port)
	http.Handle("/", uploadDir(dir))
	addr := fmt.Sprintf("%s:%d", *host, *port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
