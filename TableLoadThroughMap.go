/**
 * This GO code will move the data from CSV/DB to Final-Project's Database, Postgres.
 * A JSON is used to map-&-copy all the required tables.
 * JSON will have the column level mapping details.
 * 
 * In this Postgres is used as Final's database. So the fastest way of data copy technique of PG's "COPY" command is used.
 * COPY command with source and target are formed with the mapping JSON.
 * 
 * All tables are copied first to a temporary table, which starts with "in_" prefix.
 * Then all the data are moved to core table, with/without JOINs, casting, formatting or other techniques.
 * 
 * 
 * Copyright (C) <2017> IIT-Madras
 * 
 * License:
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 **/

package main

import (
    "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
    _ "github.com/lib/pq"
)

// our Database's Functions, Variables, Constants
// as the Postgres is used for Data-Analytic, PG functions and variables are given below.
const today_12am = "CURRENT_DATE"
const current_datetime = "now()"

/* Structure of the JSON to be formed [ */
// DB connection properties. This can be for both Source & Final's (Target) databases.
type DataBaseConfiguration struct {
    DB_UserName string `json:"db_user_name"`
    DB_Password string `json:"db_password"`
	DB_Name string `json:"db_name"`
	DB_HostName string `json:"db_hostname"`
	DB_Port string `json:"db_port"`
	PG_Bin_Path string `json:"pg_bin"`
}

// If data is given in CSV format then, following CSV-File parameters are required.
type CSVMode struct {
	FilePath string `json:"filePath"`
	Delimiter string `json:"delimiter"`
	NullString string `json:"null_value"`
	HasHeader bool `json:"has_header"`
	IncrementalCondition string `json:"incremental_condition"`
}

// If data is given directly from a DB then, following DB details are required.
type DBMode struct {
	TableName string `json:"table_name"`
	IncrementalCondition string `json:"incremental_condition"`
}

// Source (CSV/DB) & Final-DB's table column mapping.
// If source column-name & final-db column name are same, "from_to" can be used and "from" & "to" could be blank.
// Any specific conversion/casting's formula can be given in "Import_Query"
type ColumnMapping struct {
	From string `json:"from"`
	To string `json:"to"`
	From_To string `json:"from_to"`
	Import_Query string `json:"import_query"`
}

// Table's porting details like source type CSV/DB-direct-read.
// And Column-Mapping for each and every column in the table.
type ImportTable struct {
	Intermediate_TableName string `json:"intermediate_tableName"`
	Core_TableName string `json:"core_tableName"`
	CSV_Mode CSVMode `json:"csv_mode"`
	DB_Mode DBMode `json:"db_mode"`
	Column_Mapping []ColumnMapping `json:"mapping"`
}

// This is the master struct, which is parent for source(if available) and final-db connection parameters;
// And tables to be ported.
type Configuration struct {
	TargetDBConfig DataBaseConfiguration `json:"target"`
	SourceDB DataBaseConfiguration `json:"source"`
	ImportTables []ImportTable `json:"import"`
}
/* Structure of the JSON to be formed ] */


/* Load the mapping configuration json-file.
 * Load it into a srtuct object.
 */
func loadConfig() []Configuration {
    file, e := ioutil.ReadFile("./config_TableMap_Bank1.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
	
    confJson := make([]Configuration, 1)
    json.Unmarshal(file, &confJson)
	
	return confJson
}

/* Format the PSQL's COPY command with CSV import options
 * Available options are: delimiter, header and NULL-value-notation
 */
func runCopyCommand(pg_bin_path, pg_user_name, pg_password, pg_database_name, pg_hostname, pg_port, tableName, csv_filePath, csv_delimiter, csv_nullString string, csv_hasHeader bool) {
	csvDelimiter := ""
	if len(csv_delimiter) > 0 {
		csvDelimiter = "DELIMITER '"+csv_delimiter+"'"
	}
	
	csvHasHeader := ""
	if csv_hasHeader {
		csvHasHeader = "HEADER"
	} else {
		csvHasHeader = ""
	}
	
	csvNullString := "";
	if len(csv_nullString) > 0 {
		csvNullString = "NULL '"+csv_nullString+"'"
	}
	
	fmt.Println("Executing..", pg_bin_path+"psql", "-d", "postgresql://"+pg_user_name+":"+pg_password+"@"+pg_hostname+":"+pg_port+"/"+pg_database_name, "-c", "COPY "+tableName+" FROM '"+csv_filePath+"' csv "+csvDelimiter+" "+csvHasHeader+" "+csvNullString+";")
	
	cmd := exec.Command(pg_bin_path+"psql", "-d", "postgresql://"+pg_user_name+":"+pg_password+"@"+pg_hostname+":"+pg_port+"/"+pg_database_name, "-c", "COPY "+tableName+" FROM '"+csv_filePath+"' csv "+csvDelimiter+" "+csvHasHeader+" "+csvNullString+";")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		
		if strings.Contains(err.Error(), "\"psql\"") && strings.Contains(err.Error(), "not found") && strings.Contains(err.Error(), "PATH" ) {
			fmt.Println("SET Postgres's /bin under environmental variable PATH")
		}
	}
	fmt.Printf("%s\n", stdoutStderr)
}

/* Move the data from temp table, which starts with extension "in_", into actual core table.
 * Before starting the operation, validate the mapping JSON data for blank and "from_to" field verification with "from" & "to".
 */
func importTempToCoreTable(targetDB *sql.DB, table ImportTable) {
	FromColumns := ""
	ToColumns := ""
	insertQuery := ""
	
	// move the data from "in_*" table into "core" table, only if "Core_TableName" present
	if len(table.Core_TableName) > 0 {
		for _, column := range table.Column_Mapping {
			// fmt.Printf("Column: %q %q %q %q\n", column.From, column.To, column.From_To, column.Import_Query)
			
			// Validate column-name presence
			if len(column.From) == 0 && len(column.From_To) == 0 {
				panic("Column-name is blank in both From & From_To.")
			}
			if len(column.To) == 0 && len(column.From_To) == 0 {
				panic("Column-name is blank in both To & From_To.")
			}
			
			if len(column.From) == 0 {
				column.From = column.From_To
			}
			if len(column.To) == 0 {
				column.To = column.From_To
			}
			
			FromColumns = FromColumns + column.From + ","
			ToColumns = ToColumns + column.To + ","
		}
		
		FromColumns = strings.Trim(FromColumns, ",")
		ToColumns = strings.Trim(ToColumns, ",")
		
		insertQuery = "INSERT INTO " + table.Core_TableName + "(" + ToColumns + ")" + " SELECT " + FromColumns + " FROM " + table.Intermediate_TableName
		
		if len(table.CSV_Mode.IncrementalCondition) > 0 {
			insertQuery = insertQuery+" WHERE " + table.CSV_Mode.IncrementalCondition
		}
		
		insertQuery = strings.Replace(insertQuery, "today_12am", today_12am, -1)
		
		fmt.Println("Query: ", insertQuery);
		
		_, err := targetDB.Exec(insertQuery)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error Query: ", insertQuery);
		}
	}
}

/*
 * The starting function of the Import operations.
 * Creates the DB-Connection of Final's Postgres-DB and passes it through-out the program.
 * Loops through each table's details given in JSON. And move it with Postgres's COPY command.
 */
func main() {
	configurations := loadConfig()
	
	fmt.Println(configurations)
	
    TargetDBConfig := configurations[0].TargetDBConfig
	
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", TargetDBConfig.DB_UserName, TargetDBConfig.DB_Password, TargetDBConfig.DB_Name)
    targetDB, err := sql.Open("postgres", dbinfo)
    if err != nil {
        fmt.Println(err)
    }
    defer targetDB.Close()
		
	for _, table := range configurations[0].ImportTables {
		if( len(table.Intermediate_TableName) > 0 ) {
			fmt.Println("Importing :: from File: ", table.CSV_Mode.FilePath, "; into Table: ", table.Intermediate_TableName)
		} else {
			fmt.Println("Importing :: from File: ", table.CSV_Mode.FilePath, "; into Table: ", table.Core_TableName)
		}
		// Validations ::
		// Validate presence of both CSV & Source DB
		if len(table.CSV_Mode.FilePath) > 0 && len(table.DB_Mode.TableName) > 0 {
			panic("Can't have both CSV & DB mode")
		}
		
		// Clear the "in_" table
		if( len(table.Intermediate_TableName) > 0 ) {
			_, err := targetDB.Exec("TRUNCATE TABLE "+table.Intermediate_TableName+" RESTART IDENTITY")
			if err != nil {
				fmt.Println(err)
			}
		}
		
		// Load CSV into IN process table
		Target_TableName := ""
		if( len(table.Intermediate_TableName) > 0 ) {
			Target_TableName = table.Intermediate_TableName
		} else {
			Target_TableName = table.Core_TableName
		}
		runCopyCommand(TargetDBConfig.PG_Bin_Path, TargetDBConfig.DB_UserName, TargetDBConfig.DB_Password, TargetDBConfig.DB_Name, TargetDBConfig.DB_HostName, TargetDBConfig.DB_Port, Target_TableName, table.CSV_Mode.FilePath, table.CSV_Mode.Delimiter, table.CSV_Mode.NullString, table.CSV_Mode.HasHeader)
		
		// Move the date from "in_" table to core table
		if( len(table.Intermediate_TableName) > 0 ) {
			importTempToCoreTable(targetDB, table)
		}
		
		// Execute the Move-to-Core Array
		// TODO : Move the data from multiple "in_" tables into particular core table.
		
		
		fmt.Println("completed.\n")
    }
}
