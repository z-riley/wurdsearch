package main

type Storer interface {
	StorePlaceholder(any) error
}

type Store struct {
	Storer
}

func newStore() *Store {
	return &Store{}
}

func (s *Store) StorePlaceholder(any) error {
	return nil
}
