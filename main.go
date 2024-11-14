package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type ParkingApp struct {
	DB *sql.DB
}

type Car struct {
	RegistrationNumber string
	SlotNumber         int
}

func NewParkingApp() (*ParkingApp, error) {
	db, err := sql.Open("sqlite3", "parkingapp.db")
	if err != nil {
		return nil, err
	}

	sql := `
	CREATE TABLE IF NOT EXISTS parking_lot (
		slot_number INTEGER NOT NULL PRIMARY KEY,
		occupied BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS cars (
		slot_number INTEGER NOT NULL,	
		registration_number TEXT NOT NULL PRIMARY KEY
	);
	`

	_, err = db.Exec(sql)
	if err != nil {
		return nil, err
	}

	return &ParkingApp{DB: db}, nil
}

func (parkingApp *ParkingApp) CreateParkingLot(capacity string) string {
	capacityInt, err := strconv.Atoi(capacity)
	if err != nil {
		return ""
	}

	if _, err := parkingApp.DB.Exec("DELETE FROM parking_lot; DELETE FROM cars; VACUUM;"); err != nil {
		return ""
	}

	if capacityInt > 0 {
		for i := 1; i <= capacityInt; i++ {
			if _, err := parkingApp.DB.Exec("INSERT INTO parking_lot (slot_number, occupied) VALUES (?, ?)", i, false); err != nil {
				return ""
			}
		}
	}

	return ""
}

func (parkingApp *ParkingApp) Park(registrationNumber string) string {
	var slotNumber int
	if err := parkingApp.DB.QueryRow(`
	SELECT p.slot_number FROM parking_lot AS p
		LEFT JOIN cars AS c ON p.slot_number = c.slot_number
	WHERE (c.registration_number = ? AND p.occupied = 1)
		OR (c.registration_number IS NULL AND p.occupied = 0)
	ORDER BY p.slot_number ASC LIMIT 1`, registrationNumber).Scan(&slotNumber); err == sql.ErrNoRows {
		return "Sorry, parking lot is full"
	} else if err != nil {
		return ""
	}

	if _, err := parkingApp.DB.Exec("INSERT OR IGNORE INTO cars (slot_number, registration_number) VALUES (?, ?)", slotNumber, registrationNumber); err != nil {
		return ""
	}

	if _, err := parkingApp.DB.Exec("UPDATE parking_lot SET occupied = ? WHERE slot_number = ?", true, slotNumber); err != nil {
		return ""
	}

	return fmt.Sprintf("Allocated slot number: %d", slotNumber)
}

func (parkingApp *ParkingApp) Leave(registrationNumber string, duration string) string {
	var slotNumber int
	if err := parkingApp.DB.QueryRow("SELECT slot_number FROM cars WHERE registration_number = ?", registrationNumber).Scan(&slotNumber); err == sql.ErrNoRows {
		return fmt.Sprintf("Registration number %s not found", registrationNumber)
	}

	if _, err := parkingApp.DB.Exec("DELETE FROM cars WHERE registration_number = ?", registrationNumber); err != nil {
		return ""
	}

	if _, err := parkingApp.DB.Exec("UPDATE parking_lot SET occupied = ? WHERE slot_number = ?", false, slotNumber); err != nil {
		return ""
	}

	hours, err := strconv.Atoi(duration)
	if err != nil {
		return ""
	}

	charge := 10
	if hours > 2 {
		charge += (hours - 2) * 10
	}

	return fmt.Sprintf("Registration number %s with Slot Number %d is free with Charge $%d", registrationNumber, slotNumber, charge)
}

func (parkingApp *ParkingApp) Status() string {
	rows, err := parkingApp.DB.Query("SELECT p.slot_number, COALESCE(c.registration_number, '') FROM cars AS c LEFT JOIN parking_lot AS p ON c.slot_number = p.slot_number ORDER BY c.slot_number ASC")
	if err != nil {
		return ""
	}
	defer rows.Close()

	fmt.Println("Slot No. Registration No.")

	for rows.Next() {
		var slotNumber int
		var registrationNumber string

		if err := rows.Scan(&slotNumber, &registrationNumber); err != nil {
			return ""
		}

		fmt.Printf("%d. %s\n", slotNumber, registrationNumber)
	}

	return ""
}

func main() {
	var command, param1, param2 string = "", "", ""

	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 {
		param1 = os.Args[2]
	}
	if len(os.Args) > 3 {
		param2 = os.Args[3]
	}

	parkingApp, err := NewParkingApp()
	if err != nil {
		return
	}

	switch command {
	case "create_parking_lot":
		parkingApp.CreateParkingLot(param1)
	case "park":
		fmt.Println(parkingApp.Park(param1))
	case "leave":
		fmt.Println(parkingApp.Leave(param1, param2))
	case "status":
		fmt.Println(parkingApp.Status())
	}
}
