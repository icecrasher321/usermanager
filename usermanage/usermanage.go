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
	// "unsafe"
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

// only add new filenames to the end of the array

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
	if validateAll(&newUser) && validateUserName(&newUser) == nil {
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
		} else {
			fmt.Println("DATABASE does NOT exist(Check if all the required files exist)")
		}
	}
}

func UpdateRecord(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) { // STILL NEEDS DEVELOPMENT
	if checkIfRequiredFilesExist(db) {
		if queryString(usrName, "db/usernames.txt") == false {
			fmt.Println("The following user does not exist: ", usrName)
		} else {
			user := User{usrName, firstName, lastName, age, mobileNos, emailIds}
			if validateAll(&user) {
				DeleteRecord(usrName)
				CreateRecord(usrName, firstName, lastName, age, mobileNos, emailIds) // DO NOT DISPLAY MESSAGES  (Ex - CREATE message during update)
			}
		}
	}
}

func DeleteRecord(usrName string) {
	if checkIfRequiredFilesExist(db) {
		if queryString(usrName, "db/usernames.txt") == false {
			fmt.Println("The following user does not exist: ", usrName)
		} else {
			usernames := fileToArray("db/usernames.txt")
			line := lineNum(usrName, usernames)
			simpleDB := []string{"db/usernames.txt", "db/ages.txt", "db/firstnames.txt", "db/lastnames.txt"} // simple single strings
			removeStringsFromDb(line, simpleDB)
			complexDB := []string{"db/emailids.txt", "db/mobilenums.txt"} // array files
			removeArraysFromDb(line, complexDB)
			fmt.Println("The following user has been removed:", usrName)
		}
	}
}

// VALIDATIONS

func validateAll(user *User) bool {

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
	fmted := fmt.Sprintf("%s\n", text)
	_, err = file.WriteString(fmted)
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
	for i := 0; i < len(ints); i++ {
		cons[i] = strconv.Itoa(ints[i])
	}
	return cons
}

func arrStrToarrInt(strs []string) (cons []int) {
	cons = make([]int, len(strs))
	for i := 0; i < len(strs); i++ {
		cons[i], _ = strconv.Atoi(strs[i])
	}
	return cons
}

func readLine(filename string, lineNum int) (line string, lastLine int) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	checkErrorWithCount(err, &numOfErrors)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine
		}
		if lineNum == len(fileToArray("usernames.txt")) {
			break
		}
	}
	return line, lastLine
}

func stripSpaces(str string) string { // STACKOVERFLOW credited for this func
	return strings.Map(func(r rune) rune { // mapping function allows for character-wise modification of a string
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func lineNum(text string, fileArr []string) (num int) {
	for n := 0; n < len(fileArr); n++ {
		if text == fileArr[n] {
			return n
		}
	}
	return 0
}

func fetchUser(usrName string) User {
	if checkIfRequiredFilesExist(db) {
		if queryString(usrName, "db/usernames.txt") == false {
			fmt.Println("The following user does not exist: ", usrName)
		} else {
			usernames := fileToArray("db/usernames.txt")
			line := lineNum(usrName, usernames)
			userName := usernames[line]
			firstName := fetchStringFromDb(line, "db/firstnames.txt")
			lastName := fetchStringFromDb(line, "db/lastnames.txt")
			age, _ := strconv.Atoi(fetchStringFromDb(line, "db/ages.txt"))
			emailids := fetchArrayFromDB(line, "db/emailids.txt")
			mobilenums := fetchArrayFromDB(line, "db/mobilenums.txt")
			mobilenumsInts := arrStrToarrInt(mobilenums) // test
			user := User{userName, firstName, lastName, age, mobilenumsInts, emailids}
			return user
		}
	}
	return User{}
}

// Programmable arrays - ORM type method

func fileToArray(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	checkErrorWithCount(err, &numOfErrors)
	words := strings.Fields(string(content))
	return words
}

// FOR DeleteRecord

func removeStringsFromDb(lnnum int, database []string) {
	for _, field := range database {
		fileArr := fileToArray(field)
		fileArr = append(fileArr[:lnnum], fileArr[(lnnum+1):]...)
		os.Create(field)
		fileArrStr := strings.Join(fileArr, "\n")
		writeToFile(field, fileArrStr)
	}
}

func removeArraysFromDb(lnnum int, database []string) {
	for _, field := range database {
		fileArr := fileToArray(field)
		fileArrStr := strings.Join(fileArr, " ")
		fileArrSplit := strings.Split(fileArrStr, "] ")
		fileArr = append(fileArrSplit[:lnnum], fileArrSplit[(lnnum+1):]...)
		fileArrStr = strings.Join(fileArr, "]\n")
		os.Create(field)
		writeToFile(field, fileArrStr)
	}
}

// FOR UpdateRecord

func fetchStringFromDb(line int, dbFile string) string {
	fileArr := fileToArray(dbFile)
	res := fileArr[line]
	fmt.Println(res)
	return res
}

func fetchArrayFromDB(line int, dbFile string) []string {
	fileArr := fileToArray(dbFile)
	fileArrStr := strings.Join(fileArr, " ")
	fileArrStr = strings.Replace(fileArrStr, " [", "", -1)
	fileArrStr = strings.Replace(fileArrStr, " ", ", ", -1)
	fileArrSplit := strings.Split(fileArrStr, "]")

	fmt.Println(fileArrSplit[3])
	return fileArrSplit
}
