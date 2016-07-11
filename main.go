package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"text/scanner"
	"strconv"
	"bufio"
)

var fullData = "data.txt"

var usrName string
var firstName string
var lastName string
var age int
var noOfMobileNumbers int
var noOfEmailIds int
var mobileNos []int
var emailIds []string

func main() {
	var function string
	flag.StringVar(&function, "op", "create", "create, update or delete") // operation taken through a flag
	flag.Parse()
	op := strings.ToLower(function) // to prevent case-sensitivity
	if op == "create" {
		CreateRecord()
	} else if op == "update" {
		UpdateRecord()
	} else if op == "delete" {
		DeleteRecord()
	} else {
		fmt.Println("Operation does not exist. Please only use the functions create, update or delete")
	}
}

func checkError(e error) { // function will only exist during development
	if e != nil {
		defer func() {
			str := recover()
			fmt.Println(str)
		}()
		panic(e)
	}
}

// Application functionalities

func CreateRecord() { // -op create
	fmt.Println("Please enter a username")
	fmt.Scanf("%s", &usrName)
	validateUserName(usrName) // Will read from individual detail files (ex - usernames.txt) and
	// create inner arrays(using strings.Fields) so that the program can work with the data(and later write it back)
	fmt.Println("Please enter the user's first name")
	fmt.Scanf("%s", &firstName)

	fmt.Println("Please enter the user's last name")
	fmt.Scanf("%s", &lastName)

	fmt.Println("Please enter the user's age")
	fmt.Scan(&age) // Using scan instead of scanf to prevent unexpected newlines
	validateAge(age)

	fmt.Println("Please enter the number of mobile numbers the user has")
	_, err := fmt.Scan(&noOfMobileNumbers)
	checkError(err)
	mobileNos = make([]int, noOfMobileNumbers)
	var curMobileNo int
	counter := 0
	for i := 0; i < noOfMobileNumbers; i++ {
		fmt.Println("Please enter a mobile number(India only- No country code)")
		fmt.Scan(&curMobileNo)
		validateMobileNo(&curMobileNo)
		mobileNos[counter] = curMobileNo
		counter += 1
	}

	fmt.Println("Please enter the number of email ids the user has")
	fmt.Scan(&noOfEmailIds)
	emailIds = make([]string, noOfEmailIds)
	var curEmailId string
	counter = 0
	for i := 0; i < noOfEmailIds; i++ {
		fmt.Println("Please enter an email id")
		fmt.Scanf("%s", &curEmailId)
		emailIds[counter] = curEmailId
		counter += 1
	}
// File Inputs
	usrDetails := fmt.Sprintf("Username: %s\n First Name: %s\n Last Name: %s\n Age: %d\n Mobile No.[s]: %d\n Email-id[s]: %s\n \n ", usrName, firstName, lastName, age, mobileNos, emailIds)
	if err == nil {
		writeToFile("data.txt", usrDetails)
		usrName = fmt.Sprintf("%s \n", usrName)
		writeToFile("usernames.txt", usrName)
	} else {
		fmt.Println(err)
	}
}

func UpdateRecord() {
	fmt.Println("Not Ready")
}

func DeleteRecord() {
	fmt.Println("Not Ready")
}

// VALIDATIONS

func validateAge(usrAge int) {
	if usrAge <= 0 {
		fmt.Println("The age has to be an integer(above 0)")
		fmt.Scan(&age)
		validateAge(age)
	}
}

func validateUserName(name string) {
	if queryString(name, "usernames.txt") {
		fmt.Println("Username already taken(Please enter another one)")
		fmt.Scanf("%s", &usrName)
		validateUserName(usrName)
	}
}

func validateMobileNo(numberPtr *int) {
  numberString := strconv.Itoa(*numberPtr)
	if len(numberString) != 10 {
		fmt.Println("Please enter a valid mobile number (without country code)")
		fmt.Scanf("%d", numberPtr)
		validateMobileNo(numberPtr)
	}
}

// Useful functions

func queryString(str string, filename string) bool {
	content, err := ioutil.ReadFile(filename)
	checkError(err)
	words := strings.Fields(string(content))
	for _, word := range words {
		if word == str {
			return true
			break
		}
	}
	return false
}

func writeToFile(filename string, text string) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	checkError(err)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		scanner.Text()
	}
	_, err = file.WriteString(text)
	err = file.Sync()
}
