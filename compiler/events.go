package compiler

type Event interface {
}

type FileAdded struct {
	Path string
}

type FileRemoved struct {
	path string
}

type stop struct{}

var Stop = stop{}
