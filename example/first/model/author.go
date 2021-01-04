package model

type Author struct {
	ID  int
	Bio string
}

type AuthorCount struct {
	Author
	Count int
}

type SomeType struct {
	ID     int
	Values []int
	Prop1  string
	Prop2  int
}
