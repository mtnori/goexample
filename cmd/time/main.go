package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println(t)
	tUtc := t.UTC()
	fmt.Println(tUtc)

	loc := time.FixedZone("Asia/Tokyo", 9*60*60)
	t2 := tUtc.In(loc)
	fmt.Println(t2)

	s := t2.Format("2006_01_02_150405")
	fmt.Println(s)
}
