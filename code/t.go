package main

//(relational-types
type Row []string
type Relation struct {
	colNames []string
	rows     []Row
} //)

//(node-interface
type Node interface {
	// Start is called to initialize any state that this node needs to execute.
	Start()

	// Next returns the next row in the Node's result set. If there are no more
	// rows to return, the second return value will be false, otherwise, it will
	// be true.
	Next() (Row, bool)
} //)
