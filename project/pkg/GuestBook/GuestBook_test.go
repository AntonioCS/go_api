package GuestBook

import (
	"reflect"
	"testing"
)

var list = []Guest{
	{"Test1", 4, 1},
	{"Test2", 10, 2},
	{"Test3", 11, 3},
}

var emptySeats int = 10

type MockDataManager struct {
}

func (mdm MockDataManager) AddGuest(Guest) bool {
	return true
}
func (mdm MockDataManager)ListGuests() []Guest{
	return list
}
func (mdm MockDataManager)RemoveGuest(string) bool{
	return true
}
func (mdm MockDataManager)Arrived(string, int) bool{
	return true
}
func (mdm MockDataManager)FetchEmptySeatCount(string) int{
	return emptySeats
}

var guestBook = NewGuestBook(MockDataManager{}, "testVenue")

func createGuestBook() {
	NewGuestBook(MockDataManager{}, "testVenue")
}

func TestGuestBook_Add(t *testing.T) {
	if guestBook.Add("TestUser", 10, 1) != true {
		t.Errorf("Unable to add User")
	}
}

func TestGuestBook_Remove(t *testing.T) {
	if guestBook.Remove("TestUser") != true {
		t.Errorf("Unable to remove")
	}
}

func TestGuestBook_Arrived(t *testing.T) {
	if guestBook.Arrived("TestUser", 10) != true {
		t.Errorf("Unable mark arrived")
	}
}

func TestGuestBook_EmptySeats(t *testing.T) {
	if guestBook.EmptySeats() != emptySeats {
		t.Errorf("Wrong empty seats")
	}
}

func TestGuestBook_List(t *testing.T) {
	if reflect.DeepEqual(guestBook.List(),list) != true {
		t.Errorf("Wrong list of users")
	}
}


