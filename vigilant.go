package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/prometheus/alertmanager/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	// entry type
	//type Person struct {
	//	Email string
	//	Name  string
	//	Age   int
	//}

	// entry Alert
	type Alert struct {
		Job string
		Status string
	}


	// Create the DB schema
	//schema := &memdb.DBSchema{
	//	Tables: map[string]*memdb.TableSchema{
	//		"person": &memdb.TableSchema{
	//			Name: "person",
	//			Indexes: map[string]*memdb.IndexSchema{
	//				"id": &memdb.IndexSchema{
	//					Name:    "id",
	//					Unique:  true,
	//					Indexer: &memdb.StringFieldIndex{Field: "Email"},
	//				},
	//				"age": &memdb.IndexSchema{
	//					Name:    "age",
	//					Unique:  false,
	//					Indexer: &memdb.IntFieldIndex{Field: "Age"},
	//				},
	//			},
	//		},
	//	},
	//}

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"alert": &memdb.TableSchema{
				Name: "alert",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name: "id",
						Unique: true,
						Indexer: &memdb.StringFieldIndex{Field: "Job"},
					},
					"status": &memdb.IndexSchema{
						Name: "status",
						Unique: false,
						Indexer: &memdb.StringFieldIndex{Field: "Status"},
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
	//txn := db.Txn(true)

	// Insert some people
	//people := []*Person{
	//	&Person{"joe@aol.com", "Joe", 30},
	//	&Person{"lucy@aol.com", "Lucy", 35},
	//	&Person{"tariq@aol.com", "Tariq", 21},
	//	&Person{"dorothy@aol.com", "Dorothy", 53},
	//}
	//for _, p := range people {
	//	if err := txn.Insert("person", p); err != nil {
	//		panic(err)
	//	}
	//}

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
			txn := db.Txn(true)
			//if alertJson.Status == "firing" {
				for _, v := range alertJson.Alerts {
					//fmt.Println(v)
					// fmt.Println(v.Labels["job"])
					// &Person{"joe@aol.com", "Joe", 30},
					aaa := Alert{
						Job: v.Labels["job"],
						Status: v.Status,
					}
					if err := txn.Insert("alert", aaa); err != nil {
						panic(err)
					}
					//if v.Labels.Severity == "warning" {
					//	fmt.Println("The service " + v.Labels.Job + " is broken")
					//}
					//fmt.Fprintf("The service XXX status is %s", v.Status)
				}
				// Commit the transaction
				txn.Commit()
			//}
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