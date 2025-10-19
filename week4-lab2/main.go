package main

import(
	"fmt"
)

// var email string = ""

func main() {
	// var name string = "teeraphong"
	var age int = 20

	email := "teeraphonglurwongsa81@gmail.com"
	gpa := 2.93

	firstname, lastname := "teeraphong", "lurwongsa"

	fmt.Printf("Name %s %s, age %d, email %s, gpa %.2f\n", firstname, lastname, age, email, gpa)
}