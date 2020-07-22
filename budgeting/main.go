package main

import (
	// "time"
	"database/sql"
	"html/template"
	"fmt"
	_"github.com/lib/pq" //postgres driver
	"net/http"
	"strconv"
)

//database type
var db *sql.DB
//html templates type
var tpl *template.Template

//BudgetItem sets limit for spending on a specific category
type BudgetItem struct {
	ID int32
	Category string
	Limit float32
	Frequency string
}

// type Expense struct {
// 	ID int
// 	Item BudgetItem
// 	Description string
// 	ExpenseDate time.Time
// 	Charge float32
// }

const (
	host = "localhost"
	port = 5432
	user = "zeke"
	password = "hunter2"
	dbname = "budgeting"
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	var err error //must initialize here so that the global "db" variable isn't shadowed 
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", index) //redirect to budgetItemsIndex
	http.HandleFunc("/budgetItems", budgetItemsIndex) //show all budget items
	http.HandleFunc("/budgetItems/create", budgetItemsCreateForm) //display form for creating a new budget item
	http.HandleFunc("/budgetItems/create/process", budgetItemsCreateProcess) //handle the "create" POST to the db and redirect to index
	http.HandleFunc("/budgetItems/update", budgetItemsUpdateForm) //display the form (with provided info) to update the budget item 
	http.HandleFunc("/budgetItems/update/process", budgetItemsUpdateProcess) //handle the "update" POST (should it be PUT?) to the db and redirect to index
	http.HandleFunc("/budgetItems/delete/process", budgetItemsDeleteProcess) //handle the request to delete the budget item
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/budgetItems", http.StatusSeeOther)
}

func budgetItemsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM budget_items")
	if err != nil {
		fmt.Println("Internal Server Error: 1")
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	budgetItems := make([]BudgetItem, 0)
	for rows.Next() {
		item := BudgetItem{}
		err := rows.Scan(&item.ID, &item.Category, &item.Limit, &item.Frequency) // order matters
		if err != nil {
			fmt.Println("Internal Server Error: 2")
			http.Error(w, http.StatusText(500), 500)
			return
		}
		budgetItems = append(budgetItems, item)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Internal Server Error: 3")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "budgetItems.gohtml", budgetItems)

}

func budgetItemsCreateForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "create.gohtml", nil)
}

func budgetItemsCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	var item BudgetItem
	item.Category = r.FormValue("category")
	item.Frequency = r.FormValue("frequency")
	limit :=r.FormValue("expense_limit")

	// validate form values
	if item.Category == "" || limit == "" {
		fmt.Println("Bad Request: ", item, limit)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert form values
	f64, err := strconv.ParseFloat(limit, 32)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter a limit", http.StatusNotAcceptable)
		return
	}
	item.Limit = float32(f64)

	// insert values
	_, err = db.Exec("INSERT INTO budget_items (category, expense_limit, frequency) VALUES ($1, $2, $3)", item.Category, item.Limit, item.Frequency)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//redirect to home
	http.Redirect(w, r, "/budgetItems", http.StatusSeeOther)
}

func budgetItemsUpdateForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	var item BudgetItem
	err := db.QueryRow("SELECT category,expense_limit,frequency FROM budget_items WHERE item_id=$1", id).Scan(&item.Category,&item.Limit,&item.Frequency)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		fmt.Println("here", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	id32, err := strconv.ParseInt(id,0,32)
	item.ID = int32(id32)

	tpl.ExecuteTemplate(w, "update.gohtml", item)
}

func budgetItemsUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	var item BudgetItem

	item.Category = r.FormValue("category")
	item.Frequency = r.FormValue("frequency")
	limit := r.FormValue("expense_limit")
	id := r.FormValue("id")
	

	// validate form values
	if item.Frequency == "" || item.Category == "" || limit == "" || id == "" {
		fmt.Printf("Bad Request! item: %v, limit: %v, id: %v", item, limit, id)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert form values
	f64, err := strconv.ParseFloat(limit, 32)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the price", http.StatusNotAcceptable)
		return
	}
	item.Limit = float32(f64)
	
	id32, err := strconv.ParseInt(id,0,32)
	item.ID = int32(id32)


	// insert values
	_, err = db.Exec("UPDATE budget_items SET category=$2, expense_limit=$3, frequency=$4 WHERE item_id=$1;", item.ID, item.Category, item.Limit, item.Frequency)
	if err != nil {
		fmt.Println("here2", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w,r,"/budgetItems", http.StatusSeeOther)
}

func budgetItemsDeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	

	// delete book
	_, err := db.Exec("DELETE FROM budget_items WHERE item_id=$1;", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/budgetItems", http.StatusSeeOther)
}