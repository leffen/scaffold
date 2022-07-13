package db

import (
	"database/sql"
	"fmt"
	"github.com/boourns/dblib"
)

func sqlFieldsForDeparture() string {
	return "Departure.DepartureTime,Departure.ArrivalTime,Departure.VesselCode,Departure.DepartureCode,Departure.JourneyCode,Departure.DepStatus,Departure.EstimatedDepartureTime,Departure.EstimatedArrivalTime,Departure.Products,Departure.CarresRoutes,Departure.PortFrom,Departure.ToPorts" // ADD FIELD HERE
}

func loadDeparture(rows *sql.Rows) (*Departure, error) {
	ret := Departure{}

	err := rows.Scan(&ret.DepartureTime, &ret.ArrivalTime, &ret.VesselCode, &ret.DepartureCode, &ret.JourneyCode, &ret.DepStatus, &ret.EstimatedDepartureTime, &ret.EstimatedArrivalTime, &ret.Products, &ret.CarresRoutes, &ret.PortFrom, &ret.ToPorts) // ADD FIELD HERE
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func SelectDeparture(tx dblib.DBLike, cond string, condFields ...interface{}) ([]*Departure, error) {
	ret := []*Departure{}
	sql := fmt.Sprintf("SELECT %s from Departure %s", sqlFieldsForDeparture(), cond)
	rows, err := tx.Query(sql, condFields...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item, err := loadDeparture(rows)
		if err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
	rows.Close()
	return ret, nil
}

func (s *Departure) Update(tx dblib.DBLike) error {
	stmt, err := tx.Prepare(fmt.Sprintf("UPDATE Departure SET DepartureTime=?,ArrivalTime=?,VesselCode=?,DepartureCode=?,JourneyCode=?,DepStatus=?,EstimatedDepartureTime=?,EstimatedArrivalTime=?,Products=?,CarresRoutes=?,PortFrom=?,ToPorts=? WHERE Departure.ID = ?")) // ADD FIELD HERE

	if err != nil {
		return err
	}

	params := []interface{}{s.DepartureTime, s.ArrivalTime, s.VesselCode, s.DepartureCode, s.JourneyCode, s.DepStatus, s.EstimatedDepartureTime, s.EstimatedArrivalTime, s.Products, s.CarresRoutes, s.PortFrom, s.ToPorts} // ADD FIELD HERE
	params = append(params, s.ID)

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Departure) Insert(tx dblib.DBLike) error {
	stmt, err := tx.Prepare("INSERT INTO Departure(DepartureTime,ArrivalTime,VesselCode,DepartureCode,JourneyCode,DepStatus,EstimatedDepartureTime,EstimatedArrivalTime,Products,CarresRoutes,PortFrom,ToPorts) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)") // ADD FIELD HERE
	if err != nil {
		return err
	}

	result, err := stmt.Exec(s.DepartureTime, s.ArrivalTime, s.VesselCode, s.DepartureCode, s.JourneyCode, s.DepStatus, s.EstimatedDepartureTime, s.EstimatedArrivalTime, s.Products, s.CarresRoutes, s.PortFrom, s.ToPorts) // ADD FIELD HERE
	if err != nil {
		return err
	}

	s.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (s *Departure) Delete(tx dblib.DBLike) error {
	stmt, err := tx.Prepare("DELETE FROM Departure WHERE ID = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.ID)
	if err != nil {
		return err
	}

	return nil
}

func CreateDepartureTable(tx dblib.DBLike) error {
	stmt, err := tx.Prepare(`



CREATE TABLE Departure (
  
    DepartureTime DATETIME,
  
    ArrivalTime DATETIME,
  
    VesselCode VARCHAR(255),
  
    DepartureCode VARCHAR(255),
  
    JourneyCode VARCHAR(255),
  
    DepStatus VARCHAR(255),
  
    EstimatedDepartureTime DATETIME,
  
    EstimatedArrivalTime DATETIME,
  
    Products VARCHAR(255),
  
    CarresRoutes VARCHAR(255),
  
    PortFrom VARCHAR(255),
  
    ToPorts VARCHAR(255)
  
);

`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
