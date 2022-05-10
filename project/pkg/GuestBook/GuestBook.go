package GuestBook

type DataManager interface {
	AddGuest(Guest) bool
	ListGuests() []Guest
	RemoveGuest(string) bool
	Arrived(string, int) bool
	FetchEmptySeatCount(string) int
}


type Guest struct {
	Name               	string 	`json:"name"`
	AccompanyingGuests 	int		`json:"accompanyingGuests"`
	TableId				int		`json:"tableId"`
}

type GuestBook struct {
	dataManager DataManager
	venueSelected string
}

func NewGuestBook(dataManager DataManager, venue string) GuestBook {
	return GuestBook{dataManager, venue}
}


func (gb GuestBook) List() []Guest {
	return gb.dataManager.ListGuests()
}

func (gb GuestBook) Add(name string, peeps int, tableId int) bool {

	return gb.dataManager.AddGuest(Guest{
		name,
		peeps,
		tableId,
	})
}

func (gb GuestBook) Remove(name string) bool {
	return gb.dataManager.RemoveGuest(name)
}

func (gb GuestBook) Arrived(name string, peeps int) bool {
	return gb.dataManager.Arrived(name, peeps)
}

func (gb GuestBook) EmptySeats() int {
	return gb.dataManager.FetchEmptySeatCount(gb.venueSelected)
}





