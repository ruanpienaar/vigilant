package Vigilant

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/prometheus/alertmanager/template"
)

//type AlertLabels struct {
//	AlertName string `json:"alertname"`
//	Instance  string `json:"instance"`
//	Job       string `json:"job"`
//	Severity  string `json:"severity"`
//	Type      string `json:"type"`
//}
//
//type AlertAnnotations struct {
//	Summary string `json:"summary"`
//}
//
//type Alert struct {
//		Status string `json:"status"`
//		Labels AlertLabels `json:"labels"`
//		Annotations AlertAnnotations `json:"annotations"`
//		StartsAt     time.Time `json:"startsAt"`
//		EndsAt       time.Time `json:"endsAt"`
//		GeneratorURL string    `json:"generatorURL"`
//		Fingerprint  string    `json:"fingerprint"`
//}
//
//type GroupLabels struct {
//	AlertName string `json:"alertname"`
//}
//
//type CommonLabels struct {
//	AlertName string `json:"alertname"`
//	Service   string `json:"service"`
//	Severity  string `json:"severity"`
//}
//
//type CommonAnnotations struct {
//	Summary string `json:"summary"`
//}
//
//type AlertMsg struct {
//	Alerts   []Alert `json:"alerts"`
//	GroupLabels GroupLabels `json:"groupLabels"`
//	CommonLabels CommonLabels `json:"commonLabels"`
//	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
//	ExternalURL string `json:"externalURL"`
//	Version     string `json:"version"`
//	GroupKey    string `json:"groupKey"`
//}

func main() {

	// entry type
	type Person struct {
		Email string
		Name  string
		Age   int
	}

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	// Create write transaction
	txn := db.Txn(true)

	// Insert some people
	people := []*Person{
		&Person{"joe@aol.com", "Joe", 30},
		&Person{"lucy@aol.com", "Lucy", 35},
		&Person{"tariq@aol.com", "Tariq", 21},
		&Person{"dorothy@aol.com", "Dorothy", 53},
	}
	for _, p := range people {
		if err := txn.Insert("person", p); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()

	// web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("--- Handle POST ---")
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
			alertJson := GetPostJson(body)
			if alertJson.Status == "firing" {
				for _, v := range alertJson.Alerts {
					//if v.Labels.Severity == "warning" {
					//	fmt.Println("The service " + v.Labels.Job + " is broken")
					//}
					//fmt.Fprintf("The service XXX status is %s", v.Status)
					fmt.Println(v)

				}
			}
		} else {
			filepath := "www/" + r.URL.Path[1:]
			_, err := os.Open(filepath)
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("handle command "+ r.RequestURI)
				fmt.Fprintf(w, "")
			} else {
				fmt.Printf("serving file %s\n", filepath)
				http.ServeFile(w, r, filepath)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8801", nil))
}

func GetPostJson(bodyBytes []byte) template.Data {
	alertJson := template.Data{}
	err2 := json.Unmarshal([]byte(bodyBytes), &alertJson)
	if err2 != nil {
		fmt.Println(err2)
	}
	return alertJson
}