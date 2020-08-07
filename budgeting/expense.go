package main


import (
	"net/http"
	"fmt"
	"strconv"
	"time"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/araddon/dateparse"
)

//Expense is an individual record of a payment
type Expense struct {
	ID int
	Item BudgetItem
	Description string
	ExpenseDate time.Time
	Amount decimal.Decimal
}


//----- EXPENSES -----
func expensesIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	//get all expenses
	rows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	//create list of expenses
	expenses := make([]Expense, 0)
	for rows.Next() {
		var expense Expense
		var budgetItem BudgetItem
		err := rows.Scan(&expense.ID, &budgetItem.ID, &expense.Description, &expense.ExpenseDate, &expense.Amount) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		//get budget item info
		err = db.QueryRow("SELECT category,expense_limit,frequency FROM budget_items WHERE item_id=$1",budgetItem.ID).Scan(&budgetItem.Category, &budgetItem.Limit, &budgetItem.Frequency)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		expense.Item = budgetItem
		expenses = append(expenses, expense)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//display webpage
	tpl.ExecuteTemplate(w, "expenses.gohtml", expenses)

}

func expensesCreateForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "createExpenses.gohtml", nil)
}

func expensesCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	//get the budget item using the provided Category (text from form)
	var budgetItem BudgetItem
	err := db.QueryRow("SELECT * FROM budget_items WHERE category=$1", r.FormValue("category")).Scan(&budgetItem.ID, &budgetItem.Category, &budgetItem.Limit, &budgetItem.Frequency)
	switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
	}

	// get Expense values
	var expense Expense
	expense.Item = budgetItem
	expense.Description = r.FormValue("description")
	date := r.FormValue("date") //user input yyyy-mm-dd string
	amount :=r.FormValue("amount") //type string

	// validate form values
	if amount == "" || date == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert Amount from string to decimal
	a, err := decimal.NewFromString(amount)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter an amount", http.StatusNotAcceptable)
		return
	}
	expense.Amount = a

	//parse the user input into a time.Time struct
	expenseDate, err := dateparse.ParseAny(date)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(406)+" Please hit back and enter a valid date", http.StatusNotAcceptable)
		return
	}
	expense.ExpenseDate = expenseDate //yyyy-mm-dd 00:00:00 +0000 UTC

	// insert values
	_, err = db.Exec("INSERT INTO expenses (item_id, description, expense_date, amount) VALUES ($1, $2, $3, $4)", expense.Item.ID, r.FormValue("description"), expense.ExpenseDate, r.FormValue("amount"))
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//redirect to home
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func expensesUpdateForm(w http.ResponseWriter, r *http.Request) {
	//get the id from the link after the user presses "update" on an expense
	expenseID, err := strconv.ParseInt(r.FormValue("id"), 0, 64)
	if err != nil {
		http.Error(w, http.StatusText(406)+" Error with Expense ID", http.StatusNotAcceptable)
		return
	}
	//initialize variables
	var expense Expense
	expense.ID = int(expenseID)
	// var date string
	var amount string
	var budgetItem BudgetItem
	
	//get the data from the provided expense
	err = db.QueryRow("SELECT item_id, description, expense_date, amount FROM expenses WHERE expense_id=$1", expenseID).Scan(&budgetItem.ID, &expense.Description, &expense.ExpenseDate, &amount)

	//if there was a server error or the expenseID wasn't found in the database
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// convert amount from string to float32
	a, err := decimal.NewFromString(amount)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter an amount", http.StatusNotAcceptable)
		return
	}

	//get category string from budget_item id
	err = db.QueryRow("SELECT category FROM budget_items WHERE item_id=$1", budgetItem.ID).Scan(&budgetItem.Category)
	//if server error or item_id not found
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	
	//assign expense values
	expense.Item = budgetItem
	expense.Amount = a
	// expense.ExpenseDate = expenseDate
	
	tpl.ExecuteTemplate(w, "updateExpenses.gohtml", expense)
}

//todo
func expensesUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// BUDGET ITEM

	// get budget item from form 'category'
	var item BudgetItem
	item.Category = r.FormValue("category")
	var limit string
	err := db.QueryRow("SELECT item_id,expense_limit,frequency FROM budget_items WHERE category=$1",item.Category).Scan(&item.ID,&limit,&item.Frequency)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	//convert limit to float32
	l, err := decimal.NewFromString(limit)
	if err != nil {
		http.Error(w, http.StatusText(406)+" Error converting limit to float", http.StatusNotAcceptable)
		return
	}
	item.Limit = l

	fmt.Println("budget item: ", item)

	// EXPENSE

	//get expense values from form
	var expense Expense
	expense.Item = item
	expenseID, err := strconv.ParseInt(r.FormValue("id"), 0, 64)
	if err != nil {
		http.Error(w, http.StatusText(406)+" Error with Expense ID", http.StatusNotAcceptable)
		return
	}
	expense.ID = int(expenseID)
	expense.Description = r.FormValue("description")
	//parse the form value "date" into a time.Time object FROM the readable format
	expenseDate, err := dateparse.ParseAny(r.FormValue("date"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(406)+" Please hit back and enter a valid date", http.StatusNotAcceptable)
		return
	}
	expense.ExpenseDate = expenseDate
	a, err := decimal.NewFromString(r.FormValue("amount"))
	if err != nil {
		http.Error(w, http.StatusText(406)+" Please hit back and enter a number for the price", http.StatusNotAcceptable)
		return
	}
	expense.Amount = a
	
	// insert expense into db
	_, err = db.Exec("UPDATE expenses SET item_id=$1, description=$2, expense_date=$3, amount=$4 WHERE expense_id=$5;", item.ID, expense.Description, expense.ExpenseDate, expense.Amount, expense.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	fmt.Println("expense: ", expense)
	http.Redirect(w,r,"/expenses", http.StatusSeeOther)
}

//todo
func expensesDeleteProcess(w http.ResponseWriter, r *http.Request) {
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