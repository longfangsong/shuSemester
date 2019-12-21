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
	var start, end string
	row := infrastructure.DB.QueryRow(`
	SELECT id, name, lower(dateRange), upper(dateRange)
	FROM semester
	where ($1)::date <@ dateRange;
	`, date)
	err := row.Scan(&result.Id, &result.Name, &start, &end)
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	result.Start, err = time.ParseInLocation("2006-01-02", start[:len("2006-01-02")], zone)
	if err != nil {
		return result, err
	}
	result.End, err = time.ParseInLocation("2006-01-02", end[:len("2006-01-02")], zone)
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
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	shift.WorkDate = shift.WorkDate.In(zone)
	shift.RestDate = shift.RestDate.In(zone)
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO shift(fromholiday, restdate, workdate) 
	VALUES ($1, $2, $3);
	`, holidayId, shift.RestDate, shift.WorkDate)
}

func saveHoliday(holiday Holiday, semesterId int64) {
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	holiday.Start = holiday.Start.In(zone)
	holiday.End = holiday.End.In(zone)
	if holiday.Id == 0 {
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
	} else {
		_, _ = infrastructure.DB.Exec(`
		UPDATE holiday
		SET name=$2, daterange=daterange($3,$4)
		WHERE id=$1;
		`, holiday.Id, holiday.Name, holiday.Start, holiday.End)
	}
}

func deleteHoliday(holidayId int64) {
	infrastructure.DB.Exec(`
	DELETE FROM holiday
	WHERE id=$1;
	`, holidayId)
}

func Save(semester Semester) {
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	semester.Start = semester.Start.In(zone)
	semester.End = semester.End.In(zone)
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

		rows, _ := infrastructure.DB.Query(`
			SELECT id
			FROM holiday
			where belongTo = $1;
			`, semester.Id)
		holidayIdExists := map[int64]bool{}
		for rows.Next() {
			var id int64
			err := rows.Scan(&id)
			if err != nil {
				return
			}
			holidayIdExists[id] = false
		}
		for _, holiday := range semester.Holidays {
			if holiday.Id != 0 {
				holidayIdExists[holiday.Id] = true
			}
			saveHoliday(holiday, semester.Id)
		}
		for holidayId, existed := range holidayIdExists {
			if !existed {
				deleteHoliday(holidayId)
			}
		}
	}
}
