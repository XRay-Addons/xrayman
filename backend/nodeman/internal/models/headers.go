package models

type HeaderID = int

type Header struct {
	ID    HeaderID
	Key   string
	Value string
}

type Headers []Header
