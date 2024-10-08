module github.com/lazar-push/gorm-gen/tests

go 1.16

require (
	github.com/lazar-push/gorm-gen v0.3.19
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	gorm.io/driver/mysql v1.5.6
	gorm.io/driver/sqlite v1.4.4
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.9
	gorm.io/hints v1.1.1 // indirect
	gorm.io/plugin/dbresolver v1.5.0
)

replace github.com/lazar-push/gorm-gen => ../
