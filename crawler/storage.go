package main

type Storer interface {
	StorePlaceholder(any) error
}
