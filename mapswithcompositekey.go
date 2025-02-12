package main

import "fmt"

type Ckey struct {
	name string
	age  int
}

func main() {
	Mymap := make(map[Ckey]string)
	Mymap[Ckey{"Vinamra Yadav", 21}] = "TCW"
	Mymap[Ckey{"Vinamra Yadav", 22}] = "PD"
	Mymap[Ckey{"Vinamra Yadav", 23}] = "SME"
	Mymap[Ckey{"Vinamra Yadav", 24}] = "TGT"

	for key, value := range Mymap {
		fmt.Printf("%s at the age of %d worked as %s.\n", key.name, key.age, value)
	}
}
