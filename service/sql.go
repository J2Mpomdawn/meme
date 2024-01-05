package service

// test if sql works
func SqlTest() {
	query := `show tables`
	err := ExecQuery(query)

	if err != nil {
		LogPrint("red", "exec_query", err)
	}
}

func SelectTest() {
	res, err := SelectQuery("select * from guilds limit 3")

	if err != nil {
		LogPrint("red", "SelectTest", err)
	}

	FmtPrint("", res)
}

func Select(args *[]string) {
	//
	query := StrJoin_2(512, "select", *args...)

	res, err := SelectQuery(query)

	if err != nil {
		LogPrint("red", "Select", err)
		FmtPrint("red", query)
	}

	/*
		guild_id404400524099guild_name引きちぎられたヒレcreate_date2023-12-12T16:17:03+09:00update_date2023-12-12T23:20:00+09:00world_id1099
		create_date2023-12-12T16:17:03+09:00update_date2023-12-24T07:40:03+09:00world_id1099guild_id768981549099guild_nameカルデア
	*/
	for _, records := range res {
		for col, row := range records {
			FmtPrint("blue", col)
			FmtPrint("green", string(row))
		}
		FmtPrint("", "\n")
	}
}

// Query Execution Core
func ExecQuery(query string) error {
	//query exection
	_, err := DbEngine.Exec(query)

	return err
}

// Select Execution Core
func SelectQuery(query string) ([]map[string][]byte, error) {
	//select exection
	res, err := DbEngine.Query(query)

	return res, err
}
