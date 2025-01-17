package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/waltcow/sql2pb/core"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbType := flag.String("db", "mysql", "the database type")
	host := flag.String("host", "localhost", "the database host")
	port := flag.Int("port", 3306, "the database port")
	user := flag.String("user", "root", "the database user")
	password := flag.String("password", "root", "the database password")
	schema := flag.String("schema", "", "the database schema")
	table := flag.String("table", "*", "the table schema，multiple tables ',' split. ")
	serviceName := flag.String("service_name", *schema, "the protobuf service name , defaults to the database schema.")
	packageName := flag.String("package", *schema, "the protocol buffer package. defaults to the database schema.")
	goPackageName := flag.String("go_package", "", "the protocol buffer go_package. defaults to the database schema.")
	ignoreTableStr := flag.String("ignore_tables", "", "a comma spaced list of tables to ignore")

	emitConverterCode := flag.Bool("converter", false, "emit converter code")
	emitConverterPbPath := flag.String("converter_deps_pb", "", "emit converter package dependencies pb path")
	emitConverterModelPath := flag.String("converter_deps_model", "", "emit converter package dependencies model path")
	emitConverterGeneratePath := flag.String("converter_gen_path", "", "emit converter generate path")

	flag.Parse()

	if *schema == "" {
		fmt.Println(" - please input the database schema ")
		return
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *password, *host, *port, *schema)
	db, err := sql.Open(*dbType, connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ignoreTables := strings.Split(*ignoreTableStr, ",")

	s, err := core.GenerateSchema(db, *table, ignoreTables, *serviceName, *goPackageName, *packageName)

	if nil != err {
		log.Fatal(err)
	}

	if nil != s {
		var outFile string
		if *emitConverterCode {
			outFile = s.StringForConverter(*emitConverterPbPath, *emitConverterModelPath, *emitConverterGeneratePath)
		} else {
			outFile = s.String()
		}
		fmt.Println(outFile)
	}
}
