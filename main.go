package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	excelize "github.com/xuri/excelize/v2"
)

const dataFile string = "data.xlsx"

var (
	//go:embed main.html
	htmlMain string
	//go:embed find.html
	htmlFind string
	//go:embed incorrect.html
	htmlIncorrect string
)

type Data struct {
	Family string
	Ball   string
}

// Template of Find html
var tml = template.Must(template.New("Templat").Parse(htmlFind))

func getDataFromExcel(family string) Data {
	var found Data
	f, err := excelize.OpenFile(dataFile)
	if err != nil {
		fmt.Println(err)
		return found
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return found
	}

	for _, row := range rows {
		if len(row) == 2 && row[0] == family {
			found.Family = family
			found.Ball = row[1]
		}
	}
	return found
}

func getMain(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlMain)
}

func getFind(w http.ResponseWriter, r *http.Request) {
	hasFamily := r.URL.Query().Has("family")
	if !hasFamily {
		io.WriteString(w, htmlIncorrect)
		return
	}
	family := r.URL.Query().Get("family")
	data := getDataFromExcel(family)
	// Parse Find template to output
	tml.Execute(w, data)
}

func main() {
	http.HandleFunc("/", getMain)
	http.HandleFunc("/find", getFind)
	http.ListenAndServe(":3333", nil)
}
