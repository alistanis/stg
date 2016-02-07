package mysql

import(
	"bytes"
	"errors"
	"fmt"
)
// City is a generated struct from a mysql database. See github.com/alistanis/stg for more information.
type City struct {
	ID          int    `db:"ID"`
	Name        string `db:"Name"`
	CountryCode string `db:"CountryCode"`
	District    string `db:"District"`
	Population  int    `db:"Population"`
}

var (
	Cities = &City{}
)

// Find is a generated function. See github.com/alistanis/stg for more information.
func (c *City) Find(id int) (*City, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &City{}
	err := db.Select(c, "select * from City where id = ?;", id)
	if err != nil {
		return "", err
	}
	return c, nil
}

// Where is a generated function. See github.com/alistanis/stg for more information.
func (c *City) Where(column string, arg interface{}) ([]*City, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*City, 0)
	rows, err := db.Queryx("select * from City where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &City{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereFirst is a generated function. See github.com/alistanis/stg for more information.
func (c *City) WhereFirst(column string, arg interface{}) (*City, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &City{}
	err := db.Select(c, "select * from City where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	return c, nil
}

// WhereIn is a generated function. See github.com/alistanis/stg for more information.
func (c *City) WhereIn(column string, args ...interface{}) ([]*City, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	if len(args) == 0 {
		return nil, errors.New("Must provide at least one argument")
	}
	cs := make([]*City, 0)
	queryString := ""
	placeholder := "?"
	buffer := bytes.NewBuffer([]byte{})
	if placeholder == "$1" {
		queryString = "select * from City where $1 in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("$" + fmt.Sprintf("%d,", i+2))
			} else {
				buffer.WriteString("$" + fmt.Sprintf("%d);", i+2))
			}
		}
	} else {
		queryString = "select * from City where ? in ("
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
		c := &City{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereWithMap is a generated function. See github.com/alistanis/stg for more information.
func (c *City) WhereWithMap(args map[string]interface{}) ([]*City, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*City, 0)
	placeholder := "?"
	queryString := "select * from City where "

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
		c := &City{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
	