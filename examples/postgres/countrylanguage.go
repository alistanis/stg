package postgres

import(
	"bytes"
	"errors"
	"fmt"
)
// Countrylanguage is a generated struct from a postgres database. See github.com/alistanis/stg for more information.
type Countrylanguage struct {
	Countrycode string  `db:"countrycode"`
	Language    string  `db:"language"`
	Isofficial  bool    `db:"isofficial"`
	Percentage  float64 `db:"percentage"`
}

var (
	Countrylanguages = &Countrylanguage{}
)

// Where is a generated function. See github.com/alistanis/stg for more information.
func (c *Countrylanguage) Where(column string, arg interface{}) ([]*Countrylanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*Countrylanguage, 0)
	rows, err := db.Queryx("select * from countrylanguage where $1 = $2;", column, arg)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &Countrylanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereFirst is a generated function. See github.com/alistanis/stg for more information.
func (c *Countrylanguage) WhereFirst(column string, arg interface{}) (*Countrylanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &Countrylanguage{}
	err := db.Select(c, "select * from countrylanguage where $1 = $2;", column, arg)
	if err != nil {
		return "", err
	}
	return c, nil
}

// WhereIn is a generated function. See github.com/alistanis/stg for more information.
func (c *Countrylanguage) WhereIn(column string, args ...interface{}) ([]*Countrylanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	if len(args) == 0 {
		return nil, errors.New("Must provide at least one argument")
	}
	cs := make([]*Countrylanguage, 0)
	queryString := ""
	placeholder := "$1"
	buffer := bytes.NewBuffer([]byte{})
	if placeholder == "$1" {
		queryString = "select * from countrylanguage where $1 in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("$" + fmt.Sprintf("%d,", i+2))
			} else {
				buffer.WriteString("$" + fmt.Sprintf("%d);", i+2))
			}
		}
	} else {
		queryString = "select * from countrylanguage where ? in ("
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
		c := &Countrylanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereWithMap is a generated function. See github.com/alistanis/stg for more information.
func (c *Countrylanguage) WhereWithMap(args map[string]interface{}) ([]*Countrylanguage, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*Countrylanguage, 0)
	placeholder := "$1"
	queryString := "select * from countrylanguage where "

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
		c := &Countrylanguage{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
	