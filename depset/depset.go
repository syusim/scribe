package main

type Object string

type DepSet struct {
	deps map[Object][]Object
}

func New() *DepSet {
	return &DepSet{
		deps: make(map[Object][]Object),
	}
}
