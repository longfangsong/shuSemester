package model

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"shuSemester/infrastructure"
	"shuSemester/tools"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	var dbmock sqlmock.Sqlmock
	var err error
	infrastructure.DB, dbmock, err = sqlmock.New()
	tools.CheckErr(err, "cannot create Mock")
	jsonString := `{
  "start": "2019-11-23T16:00:00.000Z",
  "end": "2020-03-20T16:00:00.000Z",
  "name": "2019-2020冬季学期",
  "holidays": [
    {
      "start": "2020-01-04T16:00:00.000Z",
      "end": "2020-02-07T16:00:00.000Z",
      "name": "寒假",
      "shifts": []
    },
    {
      "start": "2019-12-31T16:00:00.000Z",
      "end": "2020-01-01T16:00:00.000Z",
      "name": "元旦",
      "shifts": []
    }
  ]
}`
	var semester Semester
	_ = json.Unmarshal([]byte(jsonString), &semester)
	zone, _ := time.LoadLocation("Asia/Shanghai")
	start := time.Date(2019, 11, 24, 0, 0, 0, 0, zone)
	end := time.Date(2020, 3, 21, 0, 0, 0, 0, zone)
	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)
	dbmock.ExpectQuery(`INSERT INTO semester`).
		WithArgs("2019-2020冬季学期", start, end).
		WillReturnRows(rows)
	start = time.Date(2020, 1, 5, 0, 0, 0, 0, zone)
	end = time.Date(2020, 2, 8, 0, 0, 0, 0, zone)
	rows = sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)
	dbmock.ExpectQuery(`INSERT INTO holiday`).
		WithArgs("寒假", 1, start, end).
		WillReturnRows(rows)
	start = time.Date(2020, 1, 1, 0, 0, 0, 0, zone)
	end = time.Date(2020, 1, 2, 0, 0, 0, 0, zone)
	rows = sqlmock.NewRows([]string{"id"})
	rows.AddRow(2)
	dbmock.ExpectQuery(`INSERT INTO holiday`).
		WithArgs("元旦", 1, start, end).
		WillReturnRows(rows)

	Save(semester)

	date := time.Date(2020, 12, 20, 0, 0, 0, 0, zone)
	rows = sqlmock.NewRows([]string{"id", "name", "lower", "upper"})
	rows.AddRow(1, "2019-2020冬季学期", "2019-11-24T00:00:00+08:00", "2020-03-21T00:00:00+08:00")
	start = time.Date(2019, 11, 24, 0, 0, 0, 0, zone)
	end = time.Date(2020, 3, 21, 0, 0, 0, 0, zone)
	dbmock.ExpectQuery(`SELECT id, name, *`).
		WithArgs(date).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "name", "lower", "upper"})
	start = time.Date(2019, 1, 5, 0, 0, 0, 0, zone)
	end = time.Date(2020, 2, 7, 0, 0, 0, 0, zone)
	rows.AddRow(1, "寒假", start, end)
	start = time.Date(2019, 1, 1, 0, 0, 0, 0, zone)
	end = time.Date(2020, 2, 2, 0, 0, 0, 0, zone)
	rows.AddRow(2, "元旦", start, end)

	dbmock.ExpectQuery(`SELECT id, name, *`).
		WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "fromholiday", "restdate", "workdate"})
	dbmock.ExpectQuery(`SELECT id, restDate, *`).
		WithArgs(1).WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "fromholiday", "restdate", "workdate"})
	dbmock.ExpectQuery(`SELECT id, restDate, *`).
		WithArgs(2).WillReturnRows(rows)
	result, err := GetByDate(date)
	if err != nil {
		t.Error(err)
	}
	if result.Name != "2019-2020冬季学期" {
		t.Error("Invalid name")
	}
	if len(result.Holidays) != 2 {
		t.Error("Holidays not fully scanned!")
	}
}
