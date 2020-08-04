package main

import (
	"database/sql"
	"html/template"
	"fmt"
	_"github.com/lib/pq" //postgres driver
	"net/http"
	"time"
	"github.com/shopspring/decimal" // safer way to store currency
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

// todo Would a map work better?

// DashboardData encapsulates the data sent to the dashboard template
type DashboardData struct {
	Today time.Time
	Monthly []TrackerData
	Weekly []TrackerData
}

// TrackerData encapsulates expense totals for a specific item
type TrackerData struct {
	Item BudgetItem
	Total decimal.Decimal
	FrequencyTotals []float32
}

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
		fmt.Println("Error getting budget items: ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	//create list of budget items
	budgetItems := make([]BudgetItem, 0)
	for rows.Next() {
		item := BudgetItem{}
		err := rows.Scan(&item.ID, &item.Category, &item.Limit, &item.Frequency) // scan in order of query
		if err != nil {
			fmt.Println("Error reading budget items: ", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		budgetItems = append(budgetItems, item)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error returned during iteration of reading budget items: ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	//append weekly and monthy totals to budgetItems
	budgetItems = appendTotals(budgetItems)

	//----- Expenses -----
	// get all expenses
	expenseRows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		fmt.Println("Error getting expenses: ", err)
		http.Error(w, http.StatusText(500),500)
		return
	}
	defer expenseRows.Close()

	// create list of expenses
	expenses := make([]Expense, 0)
	for expenseRows.Next() {
		var expense Expense
		var itemID int32
		err := expenseRows.Scan(&expense.ID, &itemID, &expense.Description, &expense.ExpenseDate, &expense.Amount)
		if err != nil {
			fmt.Println("Error reading expenses: ", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// assign item to expense
		for _, item := range budgetItems {
			if item.ID == itemID {
				expense.Item = item
			}
		}

		expenses = append(expenses, expense)
	}
	if err = expenseRows.Err(); err != nil {
		fmt.Println("Error returned during iteration of reading expenses: ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := DashboardData{
		Today: time.Now(),
		Monthly: getTrackersWithFrequency(&budgetItems, &expenses, "Monthly"),
		Weekly: getTrackersWithFrequency(&budgetItems, &expenses, "Weekly"),
	}

	tpl.ExecuteTemplate(w, "dashboard.gohtml", data)
}

// Append weeklyTotal and monthlyTotal to budgetItems list
func appendTotals(budgetItems []BudgetItem) []BudgetItem {
	weeklyTotal := BudgetItem{
		Category: "Weekly Total", //only budget items with "weekly" frequency
		Frequency: "Weekly",
	}
	monthlyTotal := BudgetItem{
		Category: "Monthly Total", // only budget items with "monthly" frequency
		Frequency: "Monthly",
	}
	// set limits
	for _, item := range budgetItems {
		if item.Frequency == "Weekly" {
			weeklyTotal.Limit = weeklyTotal.Limit.Add(item.Limit)
		}
		if item.Frequency == "Monthly" {
			monthlyTotal.Limit = monthlyTotal.Limit.Add(item.Limit)
		}
	}
	budgetItems = append(budgetItems,weeklyTotal,monthlyTotal)
	return budgetItems
}

func getTrackersWithFrequency(budgetItems *[]BudgetItem, expenses *[]Expense, frequency string) []TrackerData {
	var trackers []TrackerData
	// create a tracker for each budget item
	for _, item := range *budgetItems {
		// if item matches frequency
		if item.Frequency == frequency {
			// create tracker data with item
			tracker := TrackerData{
				Item: item,
			}
			// add tracker to list of trackers
			trackers = append(trackers, tracker)
		}
	}

	// for all expenses
	for _, expense := range *expenses {
		// for all trackers
		for i, tracker := range trackers {
			// add expense Amount to tracker total
			if expense.Item == tracker.Item {
				trackers[i].Total = trackers[i].Total.Add(expense.Amount)
			}
			// add expense Amount to appropriate frequency total
			if frequency == "Weekly" && expense.Item.Frequency == "Weekly" && tracker.Item.Category == "Weekly Total" {
				trackers[i].Total = trackers[i].Total.Add(expense.Amount)
			}
			if frequency == "Monthly" && expense.Item.Frequency == "Monthly" && tracker.Item.Category == "Monthly Total" {
				trackers[i].Total = trackers[i].Total.Add(expense.Amount)
			}

			// todo: add totals for frequency intervals
		}
	}


	
	return trackers
}