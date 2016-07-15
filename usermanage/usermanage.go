package userManage

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"text/scanner"
	"bufio"
	"errors"
	invalid "github.com/asaskevich/govalidator"
	"strconv"
	"unicode"
)

type User struct {
	UsrName   string
	FirstName string
	LastName  string
	Age       int
	MobileNos []int
	EmailIds  []string
}

// TODO: Add db does not exist error message

var numOfErrors = 0 // should remove redundant code related to this
var db = []string{"db/usernames.txt", "db/ages.txt", "db/emailids.txt", "db/firstnames.txt", "db/lastnames.txt", "db/mobilenums.txt"}

// Application functionalities

func DbReset() { // Add warning system
	for _, filename := range db {
		// os.Remove(filename)
		os.Create(filename)
	}
	fmt.Println("DATABASE RESET!")
}

func CreateRecord(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) {
	newUser := User{usrName, firstName, lastName, age, mobileNos, emailIds}
	if validateAll(&newUser) {
		if checkIfRequiredFilesExist(db) {
			newUser.UsrName = fmt.Sprintf("%s\n", newUser.UsrName)
			writeToFile("db/usernames.txt", newUser.UsrName)

			newUser.FirstName = fmt.Sprintf("%s\n", newUser.FirstName)
			writeToFile("db/firstnames.txt", newUser.FirstName)

			newUser.LastName = fmt.Sprintf("%s\n", newUser.LastName)
			writeToFile("db/lastnames.txt", newUser.LastName)

			age := fmt.Sprintf("%s\n", strconv.Itoa(newUser.Age))
			writeToFile("db/ages.txt", age)

			mobilenums := fmt.Sprintf("%s\n", arrIntToarrStr(newUser.MobileNos))
			writeToFile("db/mobilenums.txt", mobilenums)

			emailids := fmt.Sprintf("%s\n", newUser.EmailIds)
			writeToFile("db/emailids.txt", emailids)

			fmt.Println("User Created!")
		}
	}
}

func UpdateRecord() {
	fmt.Println("Not Ready")
}

func DeleteRecord() {
	fmt.Println("Not Ready")
}

// VALIDATIONS

func validateAll(user *User) bool {
	_ = validateUserName(user)

	_ = validateAge(user)

	for _, num := range user.MobileNos {
		_ = validateMobileNo(num, user)
	}

	for _, id := range user.EmailIds {
		_ = validateEmailId(id, user)
	}
	return checkTotalErrors()

}

func validateAge(user *User) error {
	if !(invalid.IsNatural(float64(user.Age))) {
		errMsg := fmt.Sprintf("The age has to be above 0 for the user %s", user.UsrName)
		err := errors.New(errMsg)
		checkErrorWithCount(err, &numOfErrors)
		return err
	}
	return nil
}

func validateUserName(user *User) error {
	user.UsrName = stripSpaces(strings.ToLower(user.UsrName))
	if queryString((user.UsrName), "db/usernames.txt") {
		errMsg := fmt.Sprintf("%s (Username) already taken", user.UsrName)
		err := errors.New(errMsg)
		checkErrorWithCount(err, &numOfErrors)
		return err
	}
	return nil
}

func validateMobileNo(number int, user *User) error {
	numberString := strconv.Itoa(number)
	if len(numberString) != 10 {
		errMsg := fmt.Sprintf("Please enter a valid mobile number (without country code) for the user: %s", user.UsrName)
		err := errors.New(errMsg)
		checkErrorWithCount(err, &numOfErrors)
		return err
	}
	return nil
}

func validateEmailId(email string, user *User) error {
	if !(invalid.IsEmail(email)) {
		errMsg := fmt.Sprintf("Please enter a valid email id for the user: %s", user.UsrName)
		err := errors.New(errMsg)
		checkErrorWithCount(err, &numOfErrors)
		return err
	}
	return nil
}

// Useful functions

func queryString(str string, filename string) bool {
	content, err := ioutil.ReadFile(filename)
	checkErrorWithCount(err, &numOfErrors)
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
	checkErrorWithCount(err, &numOfErrors)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		scanner.Text()
	}
	_, err = file.WriteString(text)
	err = file.Sync()
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func checkTotalErrors() bool {
	return (numOfErrors == 0)
}

func checkErrorWithCount(e error, errCount *int) { // Change the function to support global var numOfErrors
	if e != nil {
		fmt.Println(e)
		*errCount += 1
	}
}

func checkIfRequiredFilesExist(files []string) bool {
	for _, file := range files {
		_, err := os.OpenFile(file, os.O_RDWR, 0644)
		checkErrorWithCount(err, &numOfErrors)
	}
	return checkTotalErrors()
}

func arrIntToarrStr(ints []int) (cons []string) {
	cons = make([]string, len(ints))
	for _, num := range ints {
		for i := 0; i < len(ints); i++ {
			cons[i] = strconv.Itoa(num)
		}
	}
	return cons
}



func stripSpaces(str string) string {   // STACKOVERFLOW credited for this func
    return strings.Map(func(r rune) rune {  // mapping function allows for character-wise modification of a string
        if unicode.IsSpace(r) {
            // if the character is a space, drop it
            return -1
        }
        // else keep it in the string
        return r
    } , str)
}
