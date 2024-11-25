package main

import (
	"campus_forum_cloud/config"
	"campus_forum_cloud/storage/database"
	"fmt"
)

func main() {
	fmt.Println(config.Configs)
	var res *[]string
	database.Client.Raw("show tables").Scan(&res)
	fmt.Println(res)
}
