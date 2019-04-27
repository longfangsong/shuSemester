package model

import (
	"database/sql"
	"shuSemester/infrastructure"
	"time"
)

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Shift struct {
	Id       int64     `json:"id"`
	RestDate time.Time `json:"rest_date"`
	WorkDate time.Time `json:"work_date"`
}

type Holiday struct {
	DateRange
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Shifts []Shift `json:"shifts"`
}

type Semester struct {
	DateRange
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Holidays []Holiday `json:"holidays"`
}

func fetchShift(rows *sql.Rows) (Shift, error) {
	result := Shift{}
	err := rows.Scan(&result.Id, &result.RestDate, &result.WorkDate)
	return result, err
}

func fetchHoliday(rows *sql.Rows) (Holiday, error) {
	result := Holiday{
		Shifts: []Shift{},
	}
	err := rows.Scan(&result.Id, &result.Name, &result.Start, &result.End)
	if err != nil {
		return result, err
	}
	shiftRows, err := infrastructure.DB.Query(`
	SELECT id, restDate, workDate
	FROM shift
	WHERE fromHoliday = $1
	`, result.Id)
	if err != nil {
		return result, err
	}
	for shiftRows.Next() {
		shift, err := fetchShift(shiftRows)
		if err != nil {
			return result, err
		}
		result.Shifts = append(result.Shifts, shift)
	}
	return result, nil
}

func GetByDate(date time.Time) (Semester, error) {
	result := Semester{
		Holidays: []Holiday{},
	}
	row := infrastructure.DB.QueryRow(`
	SELECT id, name, lower(dateRange), upper(dateRange)
	FROM semester
	where ($1)::date <@ dateRange;
	`, date)
	err := row.Scan(&result.Id, &result.Name, &result.Start, &result.End)
	if err != nil {
		return result, err
	}
	rows, _ := infrastructure.DB.Query(`
	SELECT id, name, lower(dateRange), upper(dateRange)
	FROM holiday
	where belongTo = $1;
	`, result.Id)
	for rows.Next() {
		holiday, err := fetchHoliday(rows)
		if err != nil {
			return result, err
		}
		result.Holidays = append(result.Holidays, holiday)
	}
	return result, nil
}

func saveShift(shift Shift, holidayId int64) {
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO shift(fromholiday, restdate, workdate) 
	VALUES ($1, $2, $3);
	`, holidayId, shift.RestDate, shift.WorkDate)
}

func saveHoliday(holiday Holiday, semesterId int64) {
	row := infrastructure.DB.QueryRow(`
	INSERT INTO holiday(name, belongTo, dateRange) 
	VALUES ($1,$2,daterange($3,$4))
	RETURNING id;
	`, holiday.Name, semesterId, holiday.Start, holiday.End)
	var id int64
	_ = row.Scan(&id)
	for _, shift := range holiday.Shifts {
		saveShift(shift, id)
	}
}

func Save(semester Semester) {
	if semester.Id == 0 {
		row := infrastructure.DB.QueryRow(`
		INSERT INTO semester(name, dateRange) 
		VALUES ($1,daterange($2,$3))
		RETURNING id;
		`, semester.Name, semester.Start, semester.End)
		var id int64
		_ = row.Scan(&id)
		for _, holiday := range semester.Holidays {
			saveHoliday(holiday, id)
		}
	} else {
		_, _ = infrastructure.DB.Exec(`
		UPDATE semester
		SET name=$2,
		    dateRange=daterange($3,$4)
		WHERE id=$1;
		`, semester.Id, semester.Name, semester.Start, semester.End)
		for _, holiday := range semester.Holidays {
			saveHoliday(holiday, semester.Id)
		}
	}
}
