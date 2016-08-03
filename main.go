package main

import (
	fdb "github.com/icecrasher321/usermanage"
	"github.com/icrowley/fake"
	"strconv"
	"fmt"
)

func main() {
	//fdb.DbReset()
	//seed()
	//fdb.FetchUserUI("6Jacobs39992")
	//fdb.UpdateRecordUI("iccrat", "flunatic3", "100", 64, []int{9741712134, 9741719085, 9741717865}, []string{"lo@tytytytyt.com", "10000@101.com"})
	//fdb.DeleteRecord("iccrat")
	//fdb.CreateRecord("watermeloyyyt", "Space", "X", 43, []int{9741712134, 9741719085, 9741719043}, []string{"elon@musk.com"})
	//fdb.DbCreate()
	//fdb.DbDrop()
}

func seed() { // 10,000 users per minute
	fdb.DbReset()
	fmt.Println("Seeding database ......... (This could take a while)")
	for i := 1; i < 100000; i++ {
		username := fmt.Sprintf(fake.UserName() + strconv.Itoa(i))
		fdb.CreateRecord(username, fake.FirstName(),  fake.LastName(), i, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com"})
	}
}
