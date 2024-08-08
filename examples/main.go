package main

import (
	"context"
	"fmt"

	"github.com/lazar-push/gorm-gen/examples/biz"
	"github.com/lazar-push/gorm-gen/examples/conf"
	"github.com/lazar-push/gorm-gen/examples/dal"
	"github.com/lazar-push/gorm-gen/examples/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()
}

func main() {
	// start your project here
	fmt.Println("hello world")
	defer fmt.Println("bye~")

	query.SetDefault(dal.DB)
	biz.Query(context.Background())
}
