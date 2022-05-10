package DataManager

import (
	"database/sql"
	"errors"
	"github.com/antoniocs/go_api/pkg/GuestBook"
	"log"
	"strconv"
)

type DBManager struct {
	conn *sql.DB
}

func NewDBManager(conn *sql.DB) DBManager {
	return DBManager{conn}
}

func (dbm DBManager) AddGuest(guest GuestBook.Guest) bool {

	sqlStatement := "INSERT INTO `go_api`.`guest` (`name`, `accompanying_guests`, `table_id`) VALUES (?, ?, ?)"
	_, err := dbm.conn.Exec(sqlStatement, guest.Name, guest.AccompanyingGuests, guest.TableId)
	if err != nil {
		log.Println("Unable to add guest: " + guest.Name)
		return false
	}

	return true
}

func (dbm DBManager) ListGuests() []GuestBook.Guest {
	var list []GuestBook.Guest

	sqlStatement := "SELECT name, accompanying_guests, table_id FROM guest"
	rows, err := dbm.conn.Query(sqlStatement)

	if err != nil {
		log.Printf("Error listing guests: %v\n", err)
		return list
	}
	defer rows.Close()

	for rows.Next() {
		var guest GuestBook.Guest
		if err := rows.Scan(&guest.Name, &guest.AccompanyingGuests, &guest.TableId); err != nil {
			log.Fatalf("Error retriving data: %v", err)
			return nil
		}
		list = append(list, guest)
	}

	return list
}

func (dbm DBManager) RemoveGuest(name string) bool {
	sqlStatement := "DELETE FROM `go_api`.`guest` WHERE  `name` = ?"
	_, err := dbm.conn.Exec(sqlStatement, name)
	if err != nil {
		log.Printf("Error deleting guest: %v", err)
		return false
	}

	return true
}

func (dbm DBManager) Arrived(name string, peeps int) bool {
	sqlStatement := "SELECT id,name, accompanying_guests, table_id FROM guest WHERE name = ?"
	row := dbm.conn.QueryRow(sqlStatement, name)

	var gId int
	var guestName string
	var thePeeps int
	var tableId int

	switch err := row.Scan(&gId, &guestName, &thePeeps, &tableId); err {
		case sql.ErrNoRows:
			return false
		case nil:
			size, err := dbm.fetchTableSize(tableId)
			if err != nil {
				return false
			}

			if size < thePeeps {
				return false
			}

			if thePeeps != peeps {
				dbm.updateGuestPeeps(gId, thePeeps)
			}

			dbm.guestHasArrived(gId)
			return true
	}

	return false
}

func (dbm DBManager) guestHasArrived(guestId int) {
	dbm.execQuery("UPDATE guest SET is_arrived = 1 WHERE id = ?", guestId)
}

func (dbm DBManager) updateGuestPeeps(guestId, newPeeps int) {
	dbm.execQuery("UPDATE guest SET accompanying_guests = ? WHERE id = ?", newPeeps, guestId)
}

func (dbm DBManager) execQuery(sqlStatement string, args ...any)  {
	_, err := dbm.conn.Exec(sqlStatement, args...)
	if err != nil {
		log.Printf("Error executing query: %s, error: %v", sqlStatement, err)
	}
}


func (dbm DBManager) fetchTableSize(tableId int) (int, error) {
	var size int
	sqlStatement := "SELECT `size` FROM `table` WHERE id = ?"
	row := dbm.conn.QueryRow(sqlStatement, tableId)

	switch err := row.Scan(&size); err {
		case sql.ErrNoRows:
			return 0, errors.New("no table with given id")
		case nil:
			return size, nil
		default:
			panic(err)
	}
}

func (dbm DBManager) FetchEmptySeatCount(venue string) int {
	sqlStatement := `
		SELECT
		(
			SELECT
				SUM(size)
			FROM
				` + "`table`" +
			`WHERE
			venue_id = (SELECT id FROM venue WHERE name = ?)
		)
		-
			IFNULL((
				SELECT
					SUM(guest.accompanying_guests)
				FROM
					guest
				WHERE
				is_arrived = 1
		), 0) 
	AS AvailableSeats
	`
	row := dbm.conn.QueryRow(sqlStatement, venue)

	var theSeats string
	switch err := row.Scan(&theSeats); err {
		case sql.ErrNoRows:
			log.Printf("Error executing query: %v", err)
			return 0
		case nil:
			seats, _ := strconv.Atoi(theSeats)
			return seats
		default:
			panic(err)
	}

}