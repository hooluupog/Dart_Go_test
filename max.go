package main

import (
	"comparable"
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// ByAge implements comparable.Interface for []Person based on
// the Age field.
type ByAge []Person

func (a ByAge) Compare(i, j int) (value int) {
	switch {
	case a[i].Age < a[j].Age:
		value = -1
	case a[i].Age == a[j].Age:
		value = 0
	case a[i].Age > a[j].Age:
		value = 1
	}
	return value
}

func (a ByAge) Len() int { return len(a) }

func main() {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	fmt.Println(people)
	max := comparable.Max(ByAge(people))
	fmt.Println(people[max].String())

}
