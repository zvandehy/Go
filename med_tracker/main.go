package main

import (
	"time"
	"fmt"
	"github.com/graphql-go/graphql"
	_"github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"encoding/json"
)

//MedTracker keeps track of a specific medicine
type MedTracker struct {
	ID int `json:"id"`
	Name string `json:"name"`
	WaitTime int `json:"waitTime"`
	LastTime time.Time `json:"lastTime"`
	NextTime time.Time `json:"nextTime"`
}

//String() implements the fmt.Stringer interface
func (tracker *MedTracker) String() string {
	return fmt.Sprintf("You need to take %q at: %q", tracker.Name, tracker.LastTime.Format("July 2nd, 3:04 PM"))
}

func main() {
	//---Creating Graphql objects---

	
	//TODO: I don't understand the relation of GraphQL objects to SQL datatypes
	//  for example, if there was a Medicine table with "id" and "name" columns
	//  and the MedTracker referenced a MedicineType object
	//  then what would that MedicineType column be in the SQL database? (Is it a foreign key?)

	//MedTracker is the only table item we currently read
	medTrackerType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "MedTracker", //object name
			Fields: graphql.Fields{
				"id":&graphql.Field{
					Type:graphql.Int,
				},
				"name": &graphql.Field{ //pill/medicine field
					Type: graphql.String,
				},
				"waitTime": &graphql.Field{ //pill/medicine field
					Type: graphql.Int,
				},
				"lastTime": &graphql.Field{ //lastTime field
					Type:graphql.DateTime,
				},
				"nextTime": &graphql.Field{ //nextTime field
					Type:graphql.DateTime,
				},
			},
		},
	)

	//TODO: I don't understand the importance of "MutationType"

	//---Endpoints---

	//creating endpoints for mutations
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			//create a new MedTracker (no lastTime or nextTime initialized)
			"create": &graphql.Field{
				Type: medTrackerType,
				Description: "Add a new Medicine to the System",
				Args: graphql.FieldConfigArgument{
					"id":&graphql.ArgumentConfig{
						Type:graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{ //pill/medicine field
						Type: graphql.NewNonNull(graphql.String),
					},
					"waitTime": &graphql.ArgumentConfig{ //pill/medicine field
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db, err := sql.Open("sqlite3", "./pilltracker.db")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()
					_, err = db.Exec("INSERT INTO medtrackers VALUES(?,?,?,?,?)",params.Args["id"].(int), params.Args["name"].(string), params.Args["waitTime"].(int), "NULL","NULL")
					if err != nil {
						fmt.Println(err)
					}
					return nil, nil
				},
			},
			//Update endpoint for taking pill right now
			"takepillnow": &graphql.Field{
				Type: medTrackerType,
				Description: "Update the Times for the 'Named' Medicine",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{ //pill/medicine field
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db, err := sql.Open("sqlite3", "./pilltracker.db")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()
					var waitTime int
					err = db.QueryRow("SELECT waitTime FROM medtrackers WHERE name = ?", params.Args["name"].(string)).Scan(&waitTime)
					if err != nil {
						fmt.Println(err)
					}
					lastTime := time.Now().Local()
					nextTime := lastTime.Add(time.Hour * time.Duration(waitTime))
					_, err = db.Exec("UPDATE medtrackers SET lastTime = ?, nextTime = ? WHERE name = ?", lastTime, nextTime, params.Args["name"].(string))
					if err != nil {
						fmt.Println(err)
					}
					return nil, nil
				},
			},
			//update endpoint for a pill taken earlier
			"takepillwhen": &graphql.Field{
				Type: medTrackerType,
				Description: "Update the Times for the 'Named' Medicine",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{ //pill/medicine field
						Type: graphql.NewNonNull(graphql.String),
					},
					"time": &graphql.ArgumentConfig{ //pill/medicine field
						Type: graphql.NewNonNull(graphql.DateTime),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db, err := sql.Open("sqlite3", "./pilltracker.db")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()
					var waitTime int
					err = db.QueryRow("SELECT waitTime FROM medtrackers WHERE name = ?", params.Args["name"].(string)).Scan(&waitTime)
					if err != nil {
						fmt.Println(err)
					}
					lastTime := params.Args["time"].(time.Time)
					nextTime := lastTime.Add(time.Hour * time.Duration(waitTime))
					_, err = db.Exec("UPDATE medtrackers SET lastTime = ?, nextTime = ? WHERE name = ?", lastTime, nextTime, params.Args["name"].(string))
					if err != nil {
						fmt.Println(err)
					}
					return nil, nil
				},
			},
		},
	})

	//read endpoints
	fields := graphql.Fields{
		//endpoint to get a list of the medicine that I am taking
		"medtrackers": &graphql.Field {
			Type: graphql.NewList(medTrackerType),
			Description: "Get List of All MedTrackers",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, err := sql.Open("sqlite3", "./pilltracker.db")
				var meds []MedTracker
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()
				results, err := db.Query("SELECT * FROM medtrackers")
				if err != nil {
					fmt.Println(err)
				}
				//every call to Scan() must be preceded by Next()
				//Next() prepares the "next" row
				for results.Next() {
					var med MedTracker
					//copy the row data into the provided variables (Medicine)
					err = results.Scan(&med.ID, &med.Name, &med.WaitTime, &med.LastTime, &med.NextTime)
					if err != nil {
						fmt.Println(err)
					}
					meds = append(meds, med)
				} 
				return meds, nil
			},
		},
		//get a specific medtracker
		"medtracker": &graphql.Field{
			Type: medTrackerType,
			Description: "Get MedTracker By Medicine Name",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, ok := p.Args["name"]
				if ok {
					db, err := sql.Open("sqlite3", "./pilltracker.db")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()
					var tracker MedTracker
					err = db.QueryRow("SELECT * FROM medtrackers where NAME = ?", name).Scan(&tracker.ID, &tracker.Name, &tracker.WaitTime, &tracker.LastTime, &tracker.NextTime)
					if err != nil {
						fmt.Println(err)
           			}
					return tracker, nil
				}
				return nil, nil
			},
		},
	}

	//TODO: I don't completely understand what these do, but they are needed (from online tutorial)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields:fields}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: mutationType,}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create new GraphQL Schema, err %v", err)
	}

	
	
	//get medtrackers from db
	query := `
		{
			medtrackers {
				id
				name
				waitTime
				lastTime
				nextTime
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute GraphQL operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}

