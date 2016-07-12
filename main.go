package main

import (
  fdb "github.com/icecrasher321/usermanage" // fake database (.txt file)
)

 func main() {
   fdb.CreateRecord("libero101", "Vikhyath", "Mondreti", 16, []int{9741719089}, []string{"x@y.com", "y@z.com"})
}
