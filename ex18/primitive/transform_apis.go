package primitive

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

// Modify image endpoint.
func Modify(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	ext := filepath.Ext(f.Name())[1:]
	modeStr := r.FormValue("mode")
	if modeStr == "" {
		renderModeChoices(w, r, f, ext)
		return
	}
	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nStr := r.FormValue("n")
	if nStr == "" {
		renderNumShapeChoices(w, r, f, ext, Mode(mode))
		return
	}
	http.Redirect(w, r, "/img/"+filepath.Base(f.Name()), http.StatusFound)
}

var mocktempfile = tempfile

// Upload image endpoint.
func Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)[1:]
	onDisk, err := mocktempfile("", ext, "./img/")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	defer onDisk.Close()
	_, err = io.Copy(onDisk, file)
	http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)
}

func renderNumShapeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode Mode) {
	opts := []genOpts{
		{N: 10, M: mode},
		{N: 20, M: mode},
		{N: 30, M: mode},
		{N: 40, M: mode},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumShapes}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name      string
		Mode      Mode
		NumShapes int
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name:      filepath.Base(img),
			Mode:      opts[i].M,
			NumShapes: opts[i].N,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	opts := []genOpts{
		{N: 10, M: ModeCircle},
		{N: 10, M: ModeBeziers},
		{N: 10, M: ModePolygon},
		{N: 10, M: ModeCombo},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name string
		Mode Mode
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name: filepath.Base(img),
			Mode: opts[i].M,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type genOpts struct {
	N int
	M Mode
}

func genImages(rs io.ReadSeeker, ext string, opts ...genOpts) ([]string, error) {
	var ret []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		f, err := genImage(rs, ext, opt.N, opt.M)
		if err != nil {
			return nil, err
		}
		ret = append(ret, f)
	}
	return ret, nil
}

var mockTransform = Transform

func genImage(r io.Reader, ext string, numShapes int, mode Mode) (string, error) {
	out, err := mockTransform(r, ext, numShapes, WithMode(mode))
	if err != nil {
		return "", err
	}
	outFile, err := mocktempfile("", ext, "./img/")
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}
