package main

import (
	"github.com/ivalue2333/scheduler/cmd/method"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 本路径下执行 go run main.go method.go
func main() {
	conf := gen.Config{
		OutPath:      "../dal/query",
		ModelPkgPath: "../dal/rdm",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	}
	// 这个必须单独放在一个包里，否则生成不了
	g := gen.NewGenerator(conf)

	dsn := "root:123456@tcp(localhost:3306)/mydb?charset=utf8&parseTime=True&loc=Local"
	// 使用单数作为表名
	gormDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	g.UseDB(gormDB) // reuse your gorm db

	g.ApplyInterface(func(method method.InstanceMethod) {}, g.GenerateModel("scheduler_instance"))
	g.ApplyInterface(func(method method.TaskMethod) {}, g.GenerateModel("scheduler_task"))

	// Generate the code.
	g.Execute()
}
