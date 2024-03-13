package service

import (
	"encoding/json"
)

// delete
func SelectTest() {
	res, err := exec_select("select * from guilds limit 3")

	if err != nil {
		LogPrint("red", "SelectTest", err)
	}

	FmtPrint("", res)
}

// select
func Select(args *[]string) (string, error) {
	//create select query
	query := StrJoin_2(512, "select", *args...)

	//execute
	res, err := exec_select(query)

	if err != nil {
		FmtPrintln("red", query)
		LogPrint("red", "Select", err)

		return "", err
	}

	table := make([]interface{}, len(res))

	//display string generation process
	for i, records := range res {
		record := map[string]string{}

		for col, val := range records {
			val_s := string(val)
			//row_pointer is exeception
			if col == "row_pointer" {
				val_s = "・・・"
			}
			record[col] = val_s
		}

		//convert map to struct
		table[i] = ConvertMapToStruct(record)
	}

	//json string
	b, err := json.Marshal(table)

	if err != nil {
		LogPrintln("red", "Marshal", err)

		return "", err
	}

	return string(b), nil
}

// Query Execution Core
func ExecQuery(query string) error {
	//query exection
	_, err := DbEngine.Exec(query)

	return err
}

// Select Execution Core
func exec_select(query string) ([]map[string][]byte, error) {
	//select exection
	res, err := DbEngine.Query(query)

	return res, err
}
