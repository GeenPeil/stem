package bapi

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
)

// very simple
var regexpDateString = regexp.MustCompile(`^([0-9]{4})-([01][0-9])-([0-3][0-9])$`)

// Date implements the sql.Scanner and sql.Valuer interfaces. It converts from date strings following the format `yyyy-mm-dd`.
type Date struct {
	Year  uint64 `json:"year"`
	Month uint64 `json:"month"`
	Day   uint64 `json:"day"`
}

// Scan reads a date string into the Date format. The scanner does not validate that given datestring contains a valid date.
func (d *Date) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Date", value)
	}

	matches := regexpDateString.FindStringSubmatch(str)
	if len(matches) != 4 {
		return fmt.Errorf("invalid format for date string \"%s\"", str)
	}

	var err error
	d.Year, err = strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return err
	}
	d.Month, err = strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		return err
	}
	d.Day, err = strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return err
	}
	return nil
}

// Value transforms Date into a date-string which can be easily converted by postgress to a date type.
func (d *Date) Value() (driver.Value, error) {
	return fmt.Sprintf("%d-%d-%d", d.Year, d.Month, d.Day), nil
}
