package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	QListTables   = "SELECT table_name FROM information_schema.tables WHERE table_catalog=$1 AND table_schema=$2 ORDER BY table_name"
	QListColumns  = "SELECT column_name FROM information_schema.columns WHERE table_catalog=$1 AND table_schema=$2 AND table_name=$3 ORDER BY ordinal_position"
	QPKConstraint = "SELECT constraint_name FROM information_schema.table_constraints WHERE constraint_catalog=$1 AND constraint_schema=$2 AND constraint_type='PRIMARY KEY' AND table_name=$3"
	QPKColumns    = "SELECT column_name FROM information_schema.key_column_usage WHERE constraint_catalog=$1 AND constraint_schema=$2 AND table_name=$3 ORDER BY ordinal_position"
)

func main() {
	dbName := "p2p_dev"
	dbSchema := "alamisharia"
	tableName := "alami_info"

	db, err := sql.Open("pgx", "postgres://@localhost:5432/p2p_dev")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	tableRows, err := db.Query(QListTables, dbName, dbSchema)
	if err != nil {
		log.Fatalln(err)
	}
	defer tableRows.Close()
	tableNames := make([]string, 0)
	i := 0
	for tableRows.Next() {
		i++
		var name string
		err = tableRows.Scan(&name)
		if err != nil {
			log.Fatalln(err)
		}
		tableNames = append(tableNames, "\""+name+"\"")
		fmt.Println(i, ":", name)
	}

	fmt.Println(tableName + ".*")
	columnRows, err := db.Query(QListColumns, dbName, dbSchema, tableName)
	if err != nil {
		log.Fatalln(err)
	}
	defer columnRows.Close()
	columNames := make([]string, 0)
	i = 0
	for columnRows.Next() {
		i++
		var name string
		err = columnRows.Scan(&name)
		if err != nil {
			log.Fatalln(err)
		}
		columNames = append(columNames, "\""+name+"\"")
		fmt.Println(i, ":", name)
	}

	fmt.Println(tableName + " PK constraint name")
	pkRows, err := db.Query(QPKConstraint, dbName, dbSchema, tableName)
	if err != nil || !pkRows.Next() {
		log.Fatalln(err)
	}
	defer pkRows.Close()
	var pkName string
	err = pkRows.Scan(&pkName)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(pkName)

	fmt.Println(tableName + " PK column(s) name")
	pkColumnRows, err := db.Query(QPKColumns, dbName, dbSchema, tableName)
	if err != nil {
		log.Fatalln(err)
	}
	defer pkColumnRows.Close()
	pkColumns := make([]string, 0)
	for pkColumnRows.Next() {
		var name string
		err = pkColumnRows.Scan(&name)
		if err != nil {
			log.Fatalln(err)
		}
		pkColumns = append(pkColumns, "\""+name+"\"")
	}
	fmt.Println(pkColumns)

	fmt.Println(tableName + "data")
	qData := fmt.Sprintf("SELECT %s FROM %s.%s ORDER BY %s", strings.Join(columNames, ", "), dbSchema, tableName, strings.Join(pkColumns, ", "))
	fmt.Println("query:", qData)
	data := make([][]any, 0)
	dataRows, err := db.Query(qData)
	if err != nil {
		log.Fatalln(err)
	}
	for dataRows.Next() {
		row := make([]any, len(columNames))
		rowPtr := make([]any, len(columNames))
		for j := range row {
			rowPtr[j] = &row[j]
		}
		err = dataRows.Scan(rowPtr...)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, row)
		//fmt.Println(row)
	}

	row := data[0]
	dataTypes := make([]string, len(columNames))
	dataTypesMap := make(map[string]string)
	for j, key := range columNames {
		//fmt.Printf("%v, %T\n", row[j], row[j])
		if row[j] != nil {
			dataTypes[j] = reflect.TypeOf(row[j]).String()
		} else {
			dataTypes[j] = "nil"
		}
		dataTypesMap[key[1:len(key)-1]] = dataTypes[j]
	}
	fmt.Println(dataTypesMap)
	dataMap := make(map[string]any)
	for j, key := range columNames {
		dataMap[key[1:len(key)-1]] = row[j]
	}
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("json:", string(jsonData))

}
