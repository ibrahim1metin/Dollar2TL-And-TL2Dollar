package main

import (
	"encoding/json"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type currency string

var DolarValue currency

type currSet struct {
	Unit1     float32
	Unit10    float32
	Unit100   float32
	Unit1000  float32
	Unit10000 float32
}
type PostData struct {
	Amount    float32
	ModAsBool bool
	LastValue float32
	mod       string
}

var modBool bool = true
var inpFloat32 float32

func (val *currency) TurnIntoFloat() float32 {
	res := ""
	result := &res
	for i, _ := range *val {
		if (*val)[i:i+1] == "," {
			*result += "."
		} else {
			*result += string((*val)[i : i+1])
		}
	}
	numRes, _ := strconv.ParseFloat(*result, 32)
	return float32(numRes)
}
func GetDolarVal(node *html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "span" && child.FirstChild != nil {
			if child.FirstChild.Data == "DOLAR" {
				DolarValue = currency(child.NextSibling.NextSibling.FirstChild.Data)
			}
		}
		GetDolarVal(child)
	}
}
func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp := template.Must(template.ParseFiles("RequestProject/templates/index.html"))
		data := PostData{ModAsBool: modBool}
		temp.Execute(w, data)
	} else if r.Method == "POST" {
		inp := r.FormValue("DolarValue")
		mod := r.FormValue("currency")
		if mod != "" {
			modBool, _ = strconv.ParseBool(mod)
		}
		if strings.Contains(inp, ",") {
			inp = strings.Replace(inp, ",", ".", 1)
		}
		if !(r.FormValue("DolarValue") == "") {
			inpFloat, _ := strconv.ParseFloat(inp, 32)
			inpFloat32 = float32(inpFloat)
		}
		var outFloat32 float32
		if modBool == true {
			outFloat32 = Renew() * inpFloat32
		} else {
			outFloat32 = inpFloat32 / Renew()
		}
		temp := template.Must(template.ParseFiles("RequestProject/templates/index.html"))
		data := PostData{ModAsBool: modBool, Amount: outFloat32, LastValue: inpFloat32}
		temp.Execute(w, data)
	}
}
func api(w http.ResponseWriter, r *http.Request) {
	dolar := currSet{Unit1: Renew(), Unit10: 10 * Renew(), Unit100: 100 * Renew(), Unit1000: 1000 * Renew(), Unit10000: 10000 * Renew()}
	tl := currSet{Unit1: 1 / Renew(), Unit10: 10 / Renew(), Unit100: 100 / Renew(), Unit1000: 1000 / Renew(), Unit10000: 10000 / Renew()}
	json.NewEncoder(w).Encode(map[string]currSet{"dolar2TL": dolar, "TL2dolar": tl})
}
func Renew() float32 {
	resp, _ := http.Get("https://www.doviz.com/")
	doc, _ := html.Parse(resp.Body)
	GetDolarVal(doc)
	return DolarValue.TurnIntoFloat()
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/API", api)
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}
	server.ListenAndServe()
}
