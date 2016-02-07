package sqlite

import(
	"bytes"
	"errors"
	"fmt"
)
// CountryLanguage is a generated struct from a sqlite3 database. See github.com/alistanis/stg for more information.
type CountryLanguage struct {
	Id          int    `db:"id"`
	Countrycode string `db:"countrycode"`
	Language    string `db:"language"`
}

var (
	CountryLanguages = &CountryLanguage{}
)

// Find is a generated function. See github.com/alistanis/stg for more information.
func (c *CountryLanguage) Find(id int) (*CountryLanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &CountryLanguage{}
	err := db.Select(c, "select * from countryLanguage where id = ?;", id)
	if err != nil {
		return "", err
	}
	return c, nil
}

// Where is a generated function. See github.com/alistanis/stg for more information.
func (c *CountryLanguage) Where(column string, arg interface{}) ([]*CountryLanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*CountryLanguage, 0)
	rows, err := db.Queryx("select * from countryLanguage where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &CountryLanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereFirst is a generated function. See github.com/alistanis/stg for more information.
func (c *CountryLanguage) WhereFirst(column string, arg interface{}) (*CountryLanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &CountryLanguage{}
	err := db.Select(c, "select * from countryLanguage where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	return c, nil
}

// WhereIn is a generated function. See github.com/alistanis/stg for more information.
func (c *CountryLanguage) WhereIn(column string, args ...interface{}) ([]*CountryLanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	if len(args) == 0 {
		return nil, errors.New("Must provide at least one argument")
	}
	cs := make([]*CountryLanguage, 0)
	queryString := ""
	placeholder := "?"
	buffer := bytes.NewBuffer([]byte{})
	if placeholder == "$1" {
		queryString = "select * from countryLanguage where $1 in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("$" + fmt.Sprintf("%d,", i+2))
			} else {
				buffer.WriteString("$" + fmt.Sprintf("%d);", i+2))
			}
		}
	} else {
		queryString = "select * from countryLanguage where ? in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("?,")
			} else {
				buffer.WriteString("?);")
			}
		}
	}

	queryString += buffer.String()

	rows, err := db.Queryx(queryString, args)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &CountryLanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereWithMap is a generated function. See github.com/alistanis/stg for more information.
func (c *CountryLanguage) WhereWithMap(args map[string]interface{}) ([]*CountryLanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*CountryLanguage, 0)
	placeholder := "?"
	queryString := "select * from countryLanguage where "

	argsSlice := make([]interface{}, 0)
	curNum := 1
	if placeholder == "$1" {
		for k, v := range args {
			if curNum/2 <= len(args) {
				queryString += fmt.Sprintf("$%d\" = \"$%d\" and ", curNum, curNum+1)
			} else {
				queryString += fmt.Sprintf("$%d\" = \"$%d\";", curNum, curNum+1)
			}
			argsSlice = append(argsSlice, k, v)
			curNum += 2
		}
	} else {
		for k, v := range args {
			if curNum/2 <= len(args) {
				queryString += "? = ? and "
			} else {
				queryString += "? = ?;"
			}
			argsSlice = append(argsSlice, k, v)
			curNum += 2
		}
	}
	rows, err := db.Queryx(queryString, args)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &CountryLanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
	