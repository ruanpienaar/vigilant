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
					"datetime": &memdb.IndexSchema{
						Name:         "datetime",
						AllowMissing: false,
						Unique:       false,
						Indexer:      &memdb.StringFieldIndex{Field: "Status"},
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
				fmt.Println("Found next ")
				break
			} else {
				fmt.Println(next)
			}
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
				aRec := &Alert{v.GeneratorURL, "job", v.Status, v.StartsAt}
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
			filepath := "www/" + r.URL.Path[1:]
			_, err := os.Open(filepath)
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("handle command " + r.RequestURI)
				commandResponse := HandleURICommand(db, r.RequestURI)
				// fmt.Fprintf(w, commandResponse)
				w.Write(commandResponse)
			} else {
				fmt.Printf("serving file %s\n", filepath)
				http.ServeFile(w, r, filepath)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8801", nil))

} // -end of main

func GetPostJson(bodyBytes []byte) template.Data {
	alertJson := template.Data{}
	err2 := json.Unmarshal([]byte(bodyBytes), &alertJson)
	if err2 != nil {
		fmt.Println(err2)
	}
	return alertJson
}

func HandleURICommand(db *memdb.MemDB, RequestURI string) []byte {
	if RequestURI == "/api/list/all-alerts" {
		fmt.Println("/api/list/all-alerts")
		txn := db.Txn(false)
		defer txn.Abort()
		it, err := txn.Get("alert", "id")
		if err != nil {
			panic(err)
		}
		var responseJsonAlerts []AlertJson
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
		b, err := syscall.ByteSliceFromString("null")
		if err != nil {
			panic(err)
		}
		return b
	}
}
