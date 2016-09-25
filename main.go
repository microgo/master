package main

import (
	"fmt"
	"master/resource"
	"master/resource/helper"
)

func main() {

	config := resource.ResourceConfig{
		PostgreSQLLogger: false,
		IsEnablePostgres: true,
		IsEnableRedis:    true,
		IsEnableRabbit:   true,
		IsEnableElastic:  true,
	}

	r, err := resource.Init(config)
	if err != nil {
		fmt.Println("[ERROR] Connect resource fail, app will be shutdown...", err)
		return
	}
	defer r.Close()
	fmt.Println("Master testing...")
	h := helper.NewResourceHelper(r)
	fmt.Println("Init helper", h)
}
