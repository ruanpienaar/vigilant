package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/madflojo/tasks"
	"github.com/prometheus/alertmanager/template"
)

// Alert used in db
type Alert struct {
	GeneratorURL string
	Job          string
	Status       string
	StartsAt     time.Time
	EndsAt       time.Time
}

// Json used for web api calls
type AlertJson struct {
	GeneratorURL string
	Job          string
	Status       string
	StartsAt     time.Time
}

type ListAlerts struct {
	Alerts []AlertJson
}

func main() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"alert": &memdb.TableSchema{
				Name: "alert",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:         "id",
						AllowMissing: false,
						Unique:       true,
						Indexer:      &memdb.StringFieldIndex{Field: "GeneratorURL"},
					},
					"job": &memdb.IndexSchema{
						Name:         "job",
						AllowMissing: false,
						Unique:       false,
						Indexer:      &memdb.StringFieldIndex{Field: "Job"},
					},
					"status": &memdb.IndexSchema{
						Name:         "status",
						AllowMissing: false,
						Unique:       false,
						Indexer:      &memdb.StringFieldIndex{Field: "Status"},
					},
					"StartsAt": &memdb.IndexSchema{
						Name:         "StartsAt",
						AllowMissing: false,
						Unique:       false,
						Indexer:      &memdb.StringFieldIndex{Field: "StartsAt"},
					},
					"EndsAt": &memdb.IndexSchema{
						Name:         "EndsAt",
						AllowMissing: false,
						Unique:       false,
						Indexer:      &memdb.StringFieldIndex{Field: "EndsAt"},
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

	// Start the Scheduler
	scheduler := tasks.New()
	defer scheduler.Stop()

	// used for periodically checking expired alerts, and deleting
	taskFunc := func() error {
		//fmt.Println(db)
		read_txn := db.Txn(false)
		iter, iterErr := read_txn.Get("alert", "job")
		if iterErr != nil {
			panic(iterErr)
		}
		for {
			next := iter.Next()
			if next == nil {
				break
			}
			//fmt.Println(next)
		}
		//fmt.Println("Inside anonymous function")
		return nil
	}

	// re-occurring task to cleanup
	id, err := scheduler.Add(
		&tasks.Task{
			Interval: time.Duration(1 * time.Minute),
			//StartAfter: time.Now().Add(5 * time.Minute),
			TaskFunc: taskFunc,
		},
	)
	// TODO: keep track of id
	fmt.Printf("task id %s \n", id)
	if err != nil {
		panic(err)
	}

	// TODO: json parsing, makes it fail, and get stuck
	// nothing after this block - TODO: refactor all these blocks of logic :D
	// web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// POST - used for alert manager webhook.
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
				fmt.Println(v)
				//aaa := &Alert{v.Labels["job"], v.Status}
				// TODO: get job?
				fmt.Println(v.StartsAt)
				aRec := &Alert{v.GeneratorURL, "job", v.Status, v.StartsAt, v.EndsAt}
				fmt.Println(aRec)
				if err := txn.Insert("alert", aRec); err != nil {
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
			// TODO: reverse this logic, check if the cmd exists first, if not file-serve.
			filepath := "www/" + r.URL.Path[1:]
			_, err := os.Open(filepath)
			if errors.Is(err, fs.ErrNotExist) {
				//fmt.Println("handle command " + r.RequestURI)
				commandResponse := HandleURICommand(db, r.URL.Path, r.URL.Query)
				// fmt.Fprintf(w, commandResponse)
				w.Write(commandResponse)
			} else {
				//fmt.Printf("serving file %s\n", filepath)
				http.ServeFile(w, r, filepath)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8801", nil))

} // -end of main

func GetPostJson(bodyBytes []byte) template.Data {
	alertJson := template.Data{}
	err := json.Unmarshal([]byte(bodyBytes), &alertJson)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return alertJson
}

func HandleURICommand(db *memdb.MemDB, path string, qsMap map[string][]string) []byte {
	// TODO: set application/json MIME response HEADER
	if path == "/api/list/all-alerts" {
		//fmt.Println("/api/list/all-alerts")
		txn := db.Txn(false)
		defer txn.Abort()
		it, err := txn.Get("alert", "id")
		if err != nil {
			panic(err)
		}
		var responseJsonAlerts []AlertJson

		//TODO: now get the args from qsMap, and tailor the query...

		for obj := it.Next(); obj != nil; obj = it.Next() {
			DbAlert := obj.(*Alert)
			responseJsonAlerts = append(responseJsonAlerts, AlertJson{
				GeneratorURL: DbAlert.GeneratorURL,
				Job:          DbAlert.Job,
				Status:       DbAlert.Status,
			})
		}
		responseJson := ListAlerts{
			Alerts: responseJsonAlerts,
		}
		b, jsonMarshalErr := json.Marshal(responseJson)
		if jsonMarshalErr != nil {
			panic(err)
		}
		return b
	} else {
		//fmt.Printf("Unhandled URI Command %s\n", RequestURI)
		b, err := syscall.ByteSliceFromString("null")
		if err != nil {
			panic(err)
		}
		return b
	}
}
