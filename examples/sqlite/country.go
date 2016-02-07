package sqlite

import(
	"bytes"
	"errors"
	"fmt"
)
// Country is a generated struct from a sqlite3 database. See github.com/alistanis/stg for more information.
type Country struct {
	Id          int     `db:"id"`
	Name        string  `db:"name"`
	Surfacearea float64 `db:"surfacearea"`
}

var (
	Countries = &Country{}
)

// Find is a generated function. See github.com/alistanis/stg for more information.
func (c *Country) Find(id int) (*Country, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &Country{}
	err := db.Select(c, "select * from country where id = ?;", id)
	if err != nil {
		return "", err
	}
	return c, nil
}

// Where is a generated function. See github.com/alistanis/stg for more information.
func (c *Country) Where(column string, arg interface{}) ([]*Country, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*Country, 0)
	rows, err := db.Queryx("select * from country where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		c := &Country{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereFirst is a generated function. See github.com/alistanis/stg for more information.
func (c *Country) WhereFirst(column string, arg interface{}) (*Country, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	c := &Country{}
	err := db.Select(c, "select * from country where ? = ?;", column, arg)
	if err != nil {
		return "", err
	}
	return c, nil
}

// WhereIn is a generated function. See github.com/alistanis/stg for more information.
func (c *Country) WhereIn(column string, args ...interface{}) ([]*Country, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	if len(args) == 0 {
		return nil, errors.New("Must provide at least one argument")
	}
	cs := make([]*Country, 0)
	queryString := ""
	placeholder := "?"
	buffer := bytes.NewBuffer([]byte{})
	if placeholder == "$1" {
		queryString = "select * from country where $1 in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("$" + fmt.Sprintf("%d,", i+2))
			} else {
				buffer.WriteString("$" + fmt.Sprintf("%d);", i+2))
			}
		}
	} else {
		queryString = "select * from country where ? in ("
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
		c := &Country{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

// WhereWithMap is a generated function. See github.com/alistanis/stg for more information.
func (c *Country) WhereWithMap(args map[string]interface{}) ([]*Country, error) {
	if db == nil {
		panic("Must call SetDb() in order to initialize the database")
	}
	cs := make([]*Country, 0)
	placeholder := "?"
	queryString := "select * from country where "

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
		c := &Country{}
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
	