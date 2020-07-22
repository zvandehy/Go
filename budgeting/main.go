package main

import (
	// "time"
	"database/sql"
	"html/template"
	"fmt"
	_"github.com/lib/pq" //postgres driver
	"net/http"
)

//database type
var db *sql.DB
//html templates type
var tpl *template.Template

//BudgetItem sets limit for spending on a specific category
type BudgetItem struct {
	ID int
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
	http.HandleFunc("/", index)
	http.HandleFunc("/budgetitems", budgetItemsIndex)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/budgetitems", http.StatusSeeOther)
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

	tpl.ExecuteTemplate(w, "budgetitems.gohtml", budgetItems)

}