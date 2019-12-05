package store

import "testing"

func testStore(t *testing.T) *Store {
	store, err := New(nil)
	if err != nil {
		t.Fatalf("unable to connect to MySQL: %v", err)

	}
	return store
}

func TestConnect(t *testing.T) {
	testStore(t)
}
