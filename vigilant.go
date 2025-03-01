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

type alertJSON struct {
	GeneratorURL string
	Job          string
	Status       string
	StartsAt     time.Time
	EndsAt       time.Time
}

type listAlerts struct {
	Alerts []alertJSON
}

func main() {

	fmt.Println("2022-04-19 22:40:51.245939463 +0000 UTC" > "2022-04-18 22:40:51.245939463 +0000 UTC")

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
		ReadTxn := db.Txn(false)
		iter, iterErr := ReadTxn.Get("alert", "job")
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

	// re-occurring task to clean up
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
			//fmt.Println("--- Handle POST ---")
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
			alertJSON := getPostJSON(body)
			txn := db.Txn(true)
			//if alertJson.Status == "firing" {
			for _, v := range alertJson.Alerts {
				//fmt.Println(v)
				//aaa := &Alert{v.Labels["job"], v.Status}
				// TODO: get job?
				//fmt.Println(v.StartsAt)

				aRec := &Alert{v.GeneratorURL, "job", v.Status, v.StartsAt, v.EndsAt}
				//fmt.Println(aRec)
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
				commandResponse := HandleURICommand(db, r.URL.Path, r.URL.Query())
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

func getPostJSON(bodyBytes []byte) template.Data {
	alertJSON := template.Data{}
	err := json.Unmarshal([]byte(bodyBytes), &alertJSON)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return alertJSON
}

func handleURICommand(db *memdb.MemDB, path string, qsMap map[string][]string) []byte {
	// TODO: set application/json MIME response HEADER
	fmt.Println(path)
	if path == "/api/list/all-alerts" {

		//TODO: now get the args from qsMap, and tailor the query...

		txn := db.Txn(false)
		defer txn.Abort()
		//it, err := txn.LowerBound("alert", "StartsAt", "2022-04-19 22:40:51.245939463 +0000 UTC")

		// t, _ := time.Parse("2006-01-02", "2020-01-29")

		// it, err := txn.ReverseLowerBound("alert", "StartsAt", t.String())
		it, err := txn.LowerBound("alert", "StartsAt", "2022-04-19 22:40:51.245939463 +0000 UTC")

		// it, err := txn.Get("alert", "id")
		if err != nil {
			panic(err)
		}
		var responseJsonAlerts []AlertJson
		for obj := it.Next(); obj != nil; obj = it.Next() {
			DbAlert := obj.(*Alert)
			fmt.Printf("date started %s\n", DbAlert.StartsAt)
			responseJsonAlerts = append(responseJsonAlerts, AlertJson{
				GeneratorURL: DbAlert.GeneratorURL,
				Job:          DbAlert.Job,
				Status:       DbAlert.Status,
			})
		}
		responseJSON := listAlerts{
			Alerts: responseJSONAlerts,
		}
		b, jsonMarshalErr := json.Marshal(responseJSON)
		if jsonMarshalErr != nil {
			panic(err)
		}
		// log.Println(b)
		return b

	} else if path == "/api/list/print-alerts" {
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
			fmt.Println(next)
		}
		b, err := syscall.ByteSliceFromString("null")
		if err != nil {
			panic(err)
		}
		return b
	} else {
		//fmt.Printf("Unhandled URI Command %s\n", RequestURI)
		b, err := syscall.ByteSliceFromString("undefined")
		if err != nil {
			panic(err)
		}
		return b
	}
}
