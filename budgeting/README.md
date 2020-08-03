This can be a long term project that helps me (and if deployed, others) create and keep track of my budget. 
# User Stories
1. User can create a budget item with a limit and frequency ✅
2. User can view all of their budget items and the total budget
    * How should total budget be displayed with frequency? Should there be a weekly and a monthly? Should monthly budget items only be calculated in a total monthly frequency, or should I divide that amount by 5 for a weekly estimate?
3. User can update any budget items ✅
4. User can delete any budget items ✅
5. User can add expenses with a date, description, amount, and category
6. User can view a list of all of their expenses
7. User can sort, filter, or search the list of expenses based on any expense information
8. User can edit any expense
9. User can delete any expense
10. User can view a dashboard of their weekly and monthly expenses within the context of their budget items
    * The dashboard is organized like a table where each row is a budget item and each column is a week
    * (TBD) Columns could also contain a 'monthly' column
    * Each cell should be color coded so that green = below-budget, red = over-budget
    * There should be "start date" and (optional) "end date" fields that the user can change
    * There should be a "total row" that has the total amount spent per week (and total limit for the column of budget item limits)
    * There should be a "running total" column that shows how much has been spent (since 'start date')