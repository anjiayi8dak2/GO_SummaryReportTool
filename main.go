package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBList struct {
	DbName string
}

func getDBList(db *sql.DB) []string {
	// Query the DB
	var row string
	rows, err := db.Query(`SHOW DATABASES;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var listSlice []string
	for rows.Next() {
		err := rows.Scan(&row)
		if err != nil {
			log.Fatal(err)
		}
		listSlice = append(listSlice, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB list is here :", listSlice)
	return listSlice
}

func main() {
	a := app.New()
	w := a.NewWindow("Select entry widget, drop down")
	w.Resize(fyne.NewSize(400, 400))

	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "moves:moves@/")
	defer db.Close()

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	// get DB List
	var list []string
	list = getDBList(db)
	// use list to update dropdown box
	// lets show our selected entry in label
	label1 := widget.NewLabel("...")
	// dropdown/ select entry
	//[]string{} all our option goes in slice
	// s is the variable to get the selected value
	dd := widget.NewSelect(
		list,
		func(s string) {
			fmt.Printf("I selected %s as my input DB..", s)
			label1.Text = s
			label1.Refresh()
		})
	// more than one widget. so use container
	c := container.NewVBox(dd, label1)
	w.SetContent(c)
	//show and run
	w.ShowAndRun()
}
