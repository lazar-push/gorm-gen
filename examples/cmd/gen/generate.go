package main

import (
	gen "github.com/lazar-push/gorm-gen"
	"github.com/lazar-push/gorm-gen/examples/conf"
	"github.com/lazar-push/gorm-gen/examples/dal"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()

	prepare(dal.DB) // prepare table for generate
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
	})

	g.UseDB(dal.DB)

	// generate all table from database
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
