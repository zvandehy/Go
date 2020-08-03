package main

import (
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

const (
	host = "localhost"
	port = 5432
	user = "zeke"
	password = "hunter2"
	dbname = "budgeting"
)

const (
	readableDate = "Jan 02, 2006"
	userInputDate = "2006-01-02"
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

	//----- Budget Items -----
	http.HandleFunc("/budgetItems", budgetItemsIndex) //show all budget items
	http.HandleFunc("/budgetItems/create", budgetItemsCreateForm) //display form for creating a new budget item
	http.HandleFunc("/budgetItems/create/process", budgetItemsCreateProcess) //handle the "create" POST to the db and redirect to index
	http.HandleFunc("/budgetItems/update", budgetItemsUpdateForm) //display the form (with provided info) to update the budget item 
	http.HandleFunc("/budgetItems/update/process", budgetItemsUpdateProcess) //handle the "update" POST (should it be PUT?) to the db and redirect to index
	http.HandleFunc("/budgetItems/delete/process", budgetItemsDeleteProcess) //handle the request to delete the budget item

	//----- Expenses -----
	http.HandleFunc("/expenses", expensesIndex) //show all expenses
	http.HandleFunc("/expenses/create", expensesCreateForm) //display form for creating a new expense
	http.HandleFunc("/expenses/create/process", expensesCreateProcess) //handle the "create" POST to the db and redirect to index for expenses
	http.HandleFunc("/expenses/update", expensesUpdateForm) //display the form (with provided info) to update the expense
	http.HandleFunc("/expenses/update/process", expensesUpdateProcess) //handle the "update" POST (should it be PUT?) to the db and redirect to index for expenses
	http.HandleFunc("/expenses/delete/process", expensesDeleteProcess) //handle the request to delete the expense

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	//get home page
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	//----- BUDGET ITEMS -----
	//get all budget items
	rows, err := db.Query("SELECT * FROM budget_items")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	//create list of budget items
	budgetItems := make([]BudgetItem, 0)
	for rows.Next() {
		item := BudgetItem{}
		err := rows.Scan(&item.ID, &item.Category, &item.Limit, &item.Frequency) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		budgetItems = append(budgetItems, item)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	//append weekly and monthy totals to budgetItems
	budgetItems = appendTotals(budgetItems)

	//----- Expenses -----
	// get all expenses
	expenseRows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		http.Error(w, http.StatusText(500),500)
		return
	}
	defer expenseRows.Close()

	// create list of expenses
	expenses := make([]Expense, 0)
	for expenseRows.Next() {
		var expense Expense
		var itemID int
		err := rows.Scan(&expense.ID, &itemID, &expense.Description, &expense.ExpenseDate, &expense.Amount)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		expenses = append(expenses, expense)
	}


	

	tpl.ExecuteTemplate(w, "dashboard.gohtml", budgetItems)
}

// APPEND weeklyTotal and monthlyTotal to budgetItems list
func appendTotals(budgetItems []BudgetItem) []BudgetItem {
	weeklyTotal := BudgetItem{
		Category: "Weekly Total", //only budget items with "weekly" frequency
		Frequency: "Weekly",
	}
	monthlyTotal := BudgetItem{
		Category: "Monthly Total", // only budget items with "monthly" frequency
		Frequency: "Monthly",
	}
	for _, item := range budgetItems {
		if item.Frequency == "Weekly" {
			weeklyTotal.Limit = weeklyTotal.Limit + item.Limit
		}
		if item.Frequency == "Monthly" {
			monthlyTotal.Limit = monthlyTotal.Limit + item.Limit
		}
	}
	budgetItems = append(budgetItems,weeklyTotal,monthlyTotal)
	return budgetItems
}