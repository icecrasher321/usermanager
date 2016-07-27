package main

import (
	fdb "github.com/icecrasher321/usermanage"
)

func main() {
	//protocol()
	//fdb.UpdateRecord("iccraft", "flunatic3", "100", 42, []int{9741718043, 9741719085, 9741717865}, []string{"lo@tytytytyt.com", "10000@101.com"})
	//fdb.DeleteRecord("watermelo")
	//fdb.FindByUserName("icecraft")
	//fdb.CreateRecord("watermelo", "Space", "X", 43, []int{9741719086, 9741719085, 9741719043}, []string{"elon@musk.com"})
}

func protocol() {
	fdb.DbReset()
	fdb.CreateRecord("iecraf", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085, 9741719043}, []string{"1@2.com"})
	fdb.CreateRecord("icecrat", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com"})
	fdb.CreateRecord("icecrft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("icecaft", "lunatic", "1000", 3, []int{9741719086, 9741719085, 9741717865}, []string{"locker@hotmail.com", "3@4.com"})

	fdb.CreateRecord("iceraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com"})
	fdb.CreateRecord("iccraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("iecraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com"})
	fdb.CreateRecord("cecraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com", "4@5.com"})
	fdb.CreateRecord("icecft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("icraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("ecraft", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("icet", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com"})
	fdb.CreateRecord("iccrat", "Vikhyath", "Mondreti", 13, []int{9741719086, 9741719085}, []string{"1@2.com", "3@4.com"})
}
