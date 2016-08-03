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
	"github.com/ivpusic/grpool"
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

// TODO: Make all functions  return  errors

var numOfErrors = 0 // should remove redundant code related to this
var db = []string{"db/usernames.txt", "db/firstnames.txt", "db/lastnames.txt", "db/ages.txt", "db/mobilenums.txt", "db/emailids.txt"}
var LRU = []string{"db/LRU/usernames.txt", "db/LRU/firstnames.txt", "db/LRU/lastnames.txt", "db/LRU/ages.txt", "db/LRU/mobilenums.txt", "db/LRU/emailids.txt"}

// only add new filenames to the end of the array

// Application functionalities

func DbReset() bool { // Add warning system
	return DbCreate()
}

func DbCreate() bool {
	os.Mkdir("db", os.ModePerm)
	os.Mkdir("db/LRU", os.ModePerm)
	for _, filename := range db {
		// os.Remove(filename)
		_, err := os.Create(filename)
		checkErrorWithCount(err, &numOfErrors)

	}
	for _, filename := range LRU {
		// os.Remove(filename)
		_, err := os.Create(filename)
		checkErrorWithCount(err, &numOfErrors)
	}
	return checkTotalErrors()
}

func DbDrop() bool {
	for _, filename := range LRU {
		err := os.Remove(filename)
		checkErrorWithCount(err, &numOfErrors)
	}
	err := os.Remove("db/LRU")
	checkErrorWithCount(err, &numOfErrors)
	for _, filename := range db {
		err := os.Remove(filename)
		checkErrorWithCount(err, &numOfErrors)
	}
	err = os.Remove("db")
	checkErrorWithCount(err, &numOfErrors)
	return checkTotalErrors()
}

func CreateRecord(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) bool {
	newUser := User{usrName, firstName, lastName, age, mobileNos, emailIds}
	if checkIfRequiredFilesExist(db) && checkIfRequiredFilesExist(LRU) {
		if validateAll(&newUser) && validateUserName(&newUser) == nil {
			if createRecordFunc(usrName, firstName, lastName, age, mobileNos, emailIds, db) {
				if updateLRU(&newUser) {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

func UpdateRecord(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) bool { // STILL NEEDS DEVELOPMENT - NOT SYNCHRONISEd - DeleteRecord AND CreateRecord
	_ = setupValid(func() bool {
		user := User{usrName, firstName, lastName, age, mobileNos, emailIds}
		if validateAll(&user) {
			_ = DeleteRecord(usrName)
			status := CreateRecord(usrName, firstName, lastName, age, mobileNos, emailIds)
			return status
		}
		return false
	}, usrName)
	return false
}

func DeleteRecord(usrName string) bool {
	_ = setupValid(func() bool {
		if deleteRecordFunc(usrName, db) {
			_ = deleteRecordFunc(usrName, LRU)
			return true
		}
		return false
	}, usrName)
	return false
}

func FetchUser(usrName string) User {
	if checkIfRequiredFilesExist(db) && checkIfRequiredFilesExist(LRU) {
		if queryString(usrName, "db/LRU/usernames.txt") == false {
			if queryString(usrName, "db/usernames.txt") == false {
				fmt.Println("The following user does not exist: ", usrName)
				return User{}
			} else {
				return fetchUserFunc(usrName, db)
			}
		} else {
			return fetchUserFunc(usrName, LRU)
		}
	}
	return User{}
}

// For terminal OUTPUT - DEVELOPMENT stages
func FetchUserUI(usrName string) {
	if checkIfRequiredFilesExist(db) {
		if queryString(usrName, "db/usernames.txt") == false {
			fmt.Println("The following user does not exist: ", usrName)
		} else {
			user := FetchUser(usrName)
			user.printInfo()
		}
	}
}

func CreateRecordUI(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) {
	if CreateRecord(usrName, firstName, lastName, age, mobileNos, emailIds) == true {
		fmt.Println("User created!")
	}
}

func DeleteRecordUI(usrName string) {
	if DeleteRecord(usrName) == true {
		fmt.Println("The following user has been removed:", usrName)
	}
}

func UpdateRecordUI(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string) {
	if UpdateRecord(usrName, firstName, lastName, age, mobileNos, emailIds) == true {
		fmt.Println("User updated: ", usrName)
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

func validateUserName(user *User) error { // TRY REMOVE REDUNDACY BELOW
	user.UsrName = stripSpaces(strings.ToLower(user.UsrName))
	if queryString(user.UsrName, "db/LRU/usernames.txt") == false {
		if queryString(user.UsrName, "db/usernames.txt") {
			errMsg := fmt.Sprintf("%s (Username) already taken", user.UsrName)
			err := errors.New(errMsg)
			checkErrorWithCount(err, &numOfErrors)
			return err
		}
	} else {
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

func validateEmailId(email string, user *User) error { // Make email id
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
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		scanner.Text()
	}
	fmted := fmt.Sprintf("%s", text)
	_, err = file.WriteString(fmted)
	err = file.Sync()
	checkErrorWithCount(err, &numOfErrors)
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
		w, err := os.OpenFile(file, os.O_RDWR, 0644)
		checkErrorWithCount(err, &numOfErrors)
		defer w.Close()
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
	defer file.Close()
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

func setupValid(function func() bool, usrName string) bool {
	if checkIfRequiredFilesExist(db) && checkIfRequiredFilesExist(LRU) {
		if queryString(usrName, "db/LRU/usernames.txt") == false {
			if queryString(usrName, "db/usernames.txt") == false {
				fmt.Println("The following user does not exist: ", usrName)
				return false
			} else {
				function()
				return true
			}
		} else {
			function()
			return true
		}
	} else {
		fmt.Println("DATABASE does NOT exist(Check if all the required files exist)")
	}
	return false
}

// Raw functionality

func createRecordFunc(usrName string, firstName string, lastName string, age int, mobileNos []int, emailIds []string, fileArr []string) bool {
	newUser := User{usrName, firstName, lastName, age, mobileNos, emailIds}
	if checkIfRequiredFilesExist(fileArr) {

		pool := grpool.NewPool(12, 6)
		defer pool.Release()

		pool.WaitCount(6)

		pool.JobQueue <- func() {
			newUser.UsrName = fmt.Sprintf("%s\n", newUser.UsrName)
			writeToFile(fileArr[0], newUser.UsrName)
			defer pool.JobDone()
		}
		pool.JobQueue <- func() {
			newUser.FirstName = fmt.Sprintf("%s\n", newUser.FirstName)
			writeToFile(fileArr[1], newUser.FirstName)
			defer pool.JobDone()
		}
		pool.JobQueue <- func() {
			newUser.LastName = fmt.Sprintf("%s\n", newUser.LastName)
			writeToFile(fileArr[2], newUser.LastName)
			defer pool.JobDone()
		}
		pool.JobQueue <- func() {
			age := fmt.Sprintf("%s\n", strconv.Itoa(newUser.Age))
			writeToFile(fileArr[3], age)
			defer pool.JobDone()
		}
		pool.JobQueue <- func() {
			mobilenums := fmt.Sprintf("%s\n", arrIntToarrStr(newUser.MobileNos))
			writeToFile(fileArr[4], mobilenums)
			defer pool.JobDone()
		}
		pool.JobQueue <- func() {
			emailids := fmt.Sprintf("%s\n", newUser.EmailIds)
			writeToFile(fileArr[5], emailids)
			defer pool.JobDone()
		}
		pool.WaitAll()
		return true
	} else {
		fmt.Println("DATABASE does NOT exist(Check if all the required files exist)")
		return false
	}
	return false
}

func deleteRecordFunc(usrName string, fileArr []string) bool {
	if checkIfRequiredFilesExist(fileArr) {
		if queryString(usrName, fileArr[0]) {
			usernames := fileToArray(fileArr[0])
			line := lineNum(usrName, usernames)
			simpleDB := []string{fileArr[0], fileArr[1], fileArr[2], fileArr[3]} // simple single strings
			removeStringsFromDb(line, simpleDB)
			complexDB := []string{fileArr[4], fileArr[5]} // array files
			removeArraysFromDb(line, complexDB)
			return true
		}
	}
	return false
}

var mpUser = User{} // Multi-Purpose user - Used as a global variable for fetchUserFunc

func fetchUserFunc(usrName string, fileArr []string) User {
	pool := grpool.NewPool(12, 6)
	defer pool.Release()

	pool.WaitCount(6)
	usernames := fileToArray(fileArr[0])
	line := lineNum(usrName, usernames)

	pool.JobQueue <- func() {
		userName := usernames[line]
		mpUser.UsrName = userName
		defer pool.JobDone()
	}
	pool.JobQueue <- func() {
		firstName := fetchStringFromDb(line, fileArr[1])
		mpUser.FirstName = firstName
		defer pool.JobDone()
	}
	pool.JobQueue <- func() {
		lastName := fetchStringFromDb(line, fileArr[2])
		mpUser.LastName = lastName
		defer pool.JobDone()
	}
	pool.JobQueue <- func() {
		age, _ := strconv.Atoi(fetchStringFromDb(line, fileArr[3]))
		mpUser.Age = age
		defer pool.JobDone()
	}
	pool.JobQueue <- func() {
		emailids := fetchArrayFromDB(line, fileArr[5])
		mpUser.EmailIds = emailids
		defer pool.JobDone()
	}
	pool.JobQueue <- func() {
		mobilenums := fetchArrayFromDB(line, fileArr[4])
		mobilenumsInts := arrStrToarrInt(mobilenums)
		mpUser.MobileNos = mobilenumsInts
		defer pool.JobDone()
	}
	pool.WaitAll()

	updateLRU(&mpUser)
	return mpUser
}

// Programmable arrays - ORM type method

func fileToArray(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	checkErrorWithCount(err, &numOfErrors)
	words := strings.Split(string(content), "\n")
	return words
}

// FOR DeleteRecord

func removeStringsFromDb(lnnum int, database []string) { // Try remove a certain line instead of rewriting the whole file
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

// FOR Querying

func fetchStringFromDb(line int, dbFile string) string {
	fileArr := fileToArray(dbFile)
	res := fileArr[line]
	return res
}

func fetchArrayFromDB(line int, dbFile string) []string {
	fileArr := fileToArray(dbFile)
	fileArrStr := fileArr[line]
	fileArrStr = strings.Replace(fileArrStr, "[", "", -1)
	fileArrStr = strings.Replace(fileArrStr, "]", "", -1)
	fileArrSplit := strings.Fields(fileArrStr)
	return fileArrSplit
}

func (user *User) printInfo() { // Eventually, try to convert CRUD funtions to methods
	fmt.Println("Username: ", user.UsrName)
	fmt.Println("Firstname: ", user.FirstName)
	fmt.Println("Lastname: ", user.LastName)
	fmt.Println("Age: ", user.Age)
	fmt.Println("Email-Id(s): ", user.EmailIds)
	fmt.Println("Mobile-No.(s): ", user.MobileNos)
}

// Internal LRU
func updateLRU(user *User) bool {
	if checkIfRequiredFilesExist(LRU) {
		if queryString(user.UsrName, "db/LRU/usernames.txt") {
			_ = deleteRecordFunc(user.UsrName, LRU)
			_ = createRecordFunc(user.UsrName, user.FirstName, user.LastName, user.Age, user.MobileNos, user.EmailIds, LRU)
			LRUarr := fileToArray("db/LRU/usernames.txt")
			controlLRUlength(&LRUarr)
			return true
		} else {
			_ = createRecordFunc(user.UsrName, user.FirstName, user.LastName, user.Age, user.MobileNos, user.EmailIds, LRU)
			LRUarr := fileToArray("db/LRU/usernames.txt")
			controlLRUlength(&LRUarr)
			return true
		}
	}
	return false
}

func controlLRUlength(LRUarr *[]string) {
	if len(*LRUarr) > 6 {
		x := *LRUarr
		deleteRecordFunc(x[0], LRU)
	}
}
