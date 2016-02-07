package generators

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"

	"path/filepath"

	"github.com/alistanis/stg/models"
	"github.com/gedex/inflector"
)

const (
	// Postgres is the postgres string
	Postgres = "postgres"
	// MySql is the mysql string
	MySql = "mysql"
	// SQLite is the sqlite3 string
	SQLite   = "sqlite3"
	sharedDb = `
import "github.com/jmoiron/sqlx"

var (
	db *sqlx.DB
)

func SetDb(d *sqlx.DB) {
	db = d
}
`
)

func Write(s *models.Schema, dir string) error {
	srcMap, err := GenerateModels(s)
	if err != nil {
		return err
	}
	err = writeSharedDb(dir)
	if err != nil {
		return err
	}
	return writeModels(srcMap, dir)
}

func writeSharedDb(dir string) error {
	base := filepath.Base(dir)
	return ioutil.WriteFile(dir+"/shared_db.go", []byte("package "+base+"\n"+sharedDb), 0664)
}

func writeModels(models map[string][]byte, dir string) error {
	base := filepath.Base(dir)
	for k, v := range models {
		err := ioutil.WriteFile(dir+"/"+k+".go", append([]byte("package "+base+`

import(
	"bytes"
	"errors"
	"fmt"
)`), v...), 0664)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateModels(s *models.Schema) (map[string][]byte, error) {
	modelSrc := make(map[string][]byte)
	for _, table := range s.Tables {
		buffer := bytes.NewBuffer([]byte{})
		structName := inflector.Singularize(strings.Title(table.Name))
		buffer.WriteString(fmt.Sprintf("\n// %s is a generated struct from a %s database. See github.com/alistanis/stg for more information.", structName, s.Adapter))
		buffer.WriteString(fmt.Sprintf("\ntype %s struct {\n", structName))
		for _, row := range table.Rows {
			buffer.WriteString(fmt.Sprintf(" %s", strings.Title(row.Field)))
			t := row.Type
			switch {
			case strings.Contains(t, "int"), strings.Contains(t, "numeric"):
				buffer.WriteString(fmt.Sprintf(" int `db:\"%s\"` \n", row.Field))
			case strings.Contains(t, "varchar"), strings.Contains(t, "text"), strings.Contains(t, "char"):
				buffer.WriteString(fmt.Sprintf(" string `db:\"%s\"` \n", row.Field))
			case strings.Contains(t, "real"), strings.Contains(t, "float"):
				buffer.WriteString(fmt.Sprintf(" float64 `db:\"%s\"` \n", row.Field))
			case strings.Contains(t, "bool"):
				buffer.WriteString(fmt.Sprintf(" bool `db:\"%s\"` \n", row.Field))
			case strings.Contains(t, "enum"):
				buffer.WriteString(fmt.Sprintf(" interface{} `db:\"%s\"` \n", row.Field))
			case strings.Contains(t, "blob"), strings.Contains(t, "binary"):
				buffer.WriteString(fmt.Sprintf(" []byte `db:\"%s\"` \n", row.Field))
			default:
				return nil, errors.New("Unsupported type")
			}

		}
		buffer.WriteString("\n}\n")

		buffer.WriteString(fmt.Sprintf(`var (
		%s = &%s{}
		)
		`, inflector.Pluralize(strings.Title(table.Name)), strings.Title(table.Name)))
		methods, err := GenerateMethods(table, s.Adapter)
		if err != nil {
			return nil, err
		}
		for _, m := range methods {
			buffer.WriteString(m)
		}
		src, err := format.Source(buffer.Bytes())
		if err != nil {
			return nil, err
		}
		buffer.Reset()
		modelSrc[strings.ToLower(table.Name)] = src
	}
	return modelSrc, nil
}

func GenerateMethods(t *models.Table, adapter string) ([]string, error) {
	var methods []string
	for _, row := range t.Rows {
		switch row.Field {
		case "id", "Id", "ID":
			find, err := generateFind(t, adapter)
			if err != nil {
				return nil, err
			}
			methods = append(methods, find)
		}
	}
	where, err := generateWhere(t, adapter)
	if err != nil {
		return nil, err
	}
	methods = append(methods, where)
	whereFirst, err := generateWhereFirst(t, adapter)
	if err != nil {
		return nil, err
	}
	methods = append(methods, whereFirst)
	whereIn, err := generateWherein(t, adapter)
	if err != nil {
		return nil, err
	}
	methods = append(methods, whereIn)
	whereWithMap, err := generateWhereWithMap(t, adapter)
	if err != nil {
		return nil, err
	}
	methods = append(methods, whereWithMap)
	return methods, nil
}

func generateFind(t *models.Table, adapter string) (string, error) {
	placeholder := getPlaceholder(adapter)
	title := strings.Title(t.Name)
	firstChar := strings.ToLower(string(t.Name[0]))
	funcPlaceholderText := FormatFunction(title, "Find", "id int", fmt.Sprintf("(*%s, error)", title))
	queryText := fmt.Sprintf(`"select * from %s where id = %s;"`, t.Name, placeholder)
	typeName := inflector.Singularize(strings.Title(t.Name))
	srcText := fmt.Sprintf(`	%s := &%s{}
	err := db.Select(%s, %s, id)
	if err != nil {
		return "", err
	}
	return %s, nil`, firstChar, typeName, firstChar, queryText, firstChar)
	functionText := fmt.Sprintf(funcPlaceholderText, srcText)
	return functionText, nil
}

func generateWhereFirst(t *models.Table, adapter string) (string, error) {
	placeholder := getPlaceholder(adapter)
	title := strings.Title(t.Name)
	firstChar := strings.ToLower(string(t.Name[0]))
	funcPlaceholderText := FormatFunction(title, "WhereFirst", "column string, arg interface{}", fmt.Sprintf("(*%s, error)", title))
	queryText := `"select * from %s where %s = %s;"`
	if placeholder == "?" {
		queryText = fmt.Sprintf(queryText, t.Name, placeholder, placeholder)
	} else {
		queryText = fmt.Sprintf(queryText, t.Name, placeholder, "$2")
	}
	typeName := inflector.Singularize(strings.Title(t.Name))
	srcText := fmt.Sprintf(`	%s := &%s{}
	err := db.Select(%s, %s, column, arg)
	if err != nil {
		return "", err
	}
	return %s, nil`, firstChar, typeName, firstChar, queryText, firstChar)
	functionText := fmt.Sprintf(funcPlaceholderText, srcText)
	return functionText, nil
}

func generateWhere(t *models.Table, adapter string) (string, error) {
	placeholder := getPlaceholder(adapter)
	title := strings.Title(t.Name)
	firstChar := strings.ToLower(string(t.Name[0]))
	funcPlaceholderText := FormatFunction(title, "Where", "column string, arg interface{}", fmt.Sprintf("([]*%s, error)", title))
	queryText := `"select * from %s where %s = %s;"`
	if placeholder == "?" {
		queryText = fmt.Sprintf(queryText, t.Name, placeholder, placeholder)
	} else {
		queryText = fmt.Sprintf(queryText, t.Name, placeholder, "$2")
	}
	typeName := inflector.Singularize(strings.Title(t.Name))
	srcText := fmt.Sprintf(`	%s := make([]*%s, 0)
	rows, err := db.Queryx(%s, column, arg)
	if err != nil {
		return "", err
	}
	for rows.Next() {
			%s := &%s{}
			err = rows.StructScan(%s)
			if err != nil {
				return nil, err
			}
			%s = append(%s, %s)
		}
	return %s, nil`, firstChar+"s", typeName, queryText, firstChar, typeName, firstChar, firstChar+"s", firstChar+"s", firstChar, firstChar+"s")

	functionText := fmt.Sprintf(funcPlaceholderText, srcText)
	return functionText, nil
}

func generateWherein(t *models.Table, adapter string) (string, error) {
	placeholder := getPlaceholder(adapter)
	title := strings.Title(t.Name)
	firstChar := strings.ToLower(string(t.Name[0]))
	funcPlaceholderText := FormatFunction(title, "WhereIn", "column string, args ...interface{}", fmt.Sprintf("([]*%s, error)", title))
	typeName := inflector.Singularize(strings.Title(t.Name))
	srcText := strings.Replace(fmt.Sprintf(`	if len(args) == 0 {
		return nil, errors.New("Must provide at least one argument")
	}
	%s := make([]*%s, 0)
	queryString := ""
	placeholder := "%s"
	buffer := bytes.NewBuffer([]byte{})
	if placeholder == "$1" {
		queryString = "select * from %s where $1 in ("
		for i, _ := range args {
			if i < len(args) {
				buffer.WriteString("$" + fmt.Sprintf("#{d},", i+2))
			} else {
				buffer.WriteString("$" + fmt.Sprintf("#{d});", i+2))
			}
		}
	} else {
		queryString = "select * from %s where ? in ("
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
			%s := &%s{}
			err = rows.StructScan(%s)
			if err != nil {
				return nil, err
			}
			%s = append(%s, %s)
		}
	return %s, nil`, firstChar+"s", typeName, placeholder, t.Name, t.Name, firstChar, title, firstChar, firstChar+"s", firstChar+"s", firstChar, firstChar+"s"), "#{d}", "%d", -1)
	functionText := fmt.Sprintf(funcPlaceholderText, srcText)
	return functionText, nil
}

func generateWhereWithMap(t *models.Table, adapter string) (string, error) {
	placeholder := getPlaceholder(adapter)
	title := strings.Title(t.Name)
	firstChar := strings.ToLower(string(t.Name[0]))
	funcPlaceholderText := FormatFunction(title, "WhereWithMap", "args map[string]interface{}", fmt.Sprintf("([]*%s, error)", title))

	typeName := inflector.Singularize(strings.Title(t.Name))
	srcText := strings.Replace(
		fmt.Sprintf(`	%s := make([]*%s, 0)
	placeholder := "%s"
	queryString := "select * from %s where "

	argsSlice := make([]interface{}, 0)
	curNum := 1
	if placeholder == "$1" {
		for k, v := range args {
			if curNum / 2 <= len(args) {
				queryString += fmt.Sprintf("$#{d}\" = \"$#{d}\" and ", curNum, curNum +1)
			} else {
				queryString += fmt.Sprintf("$#{d}\" = \"$#{d}\";", curNum, curNum +1)
			}
			argsSlice = append(argsSlice, k, v)
			curNum +=2
		}
	} else {
		for k, v := range args {
			if curNum / 2 <= len(args) {
				queryString += "? = ? and "
			} else {
				queryString += "? = ?;"
			}
			argsSlice = append(argsSlice, k, v)
			curNum +=2
		}
	}
	rows, err := db.Queryx(queryString, args)
	if err != nil {
		return "", err
	}
	for rows.Next() {
			%s := &%s{}
			err = rows.StructScan(%s)
			if err != nil {
				return nil, err
			}
			%s = append(%s, %s)
		}
	return %s, nil`, firstChar+"s", typeName, placeholder, t.Name, firstChar, typeName, firstChar, firstChar+"s", firstChar+"s", firstChar, firstChar+"s"),
		"#{d}", "%d", -1)

	functionText := fmt.Sprintf(funcPlaceholderText, srcText)
	return functionText, nil
}

func getPlaceholder(adapter string) string {
	placeholder := ""
	switch adapter {
	case MySql, SQLite:
		placeholder = "?"
	case Postgres:
		placeholder = "$1"
	}
	return placeholder
}

func FormatFunction(typeName, functionName, params, retVals string) string {
	firstChar := strings.ToLower(string(typeName[0]))
	return fmt.Sprintf(`// %s is a generated function. See github.com/alistanis/stg for more information.
func (%s *%s)%s (%s) %s {
if db == nil {
	panic("Must call SetDb() in order to initialize the database")
}
%s
}
	`, functionName, firstChar, typeName, functionName, params, retVals, "%s")
}
