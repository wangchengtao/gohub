package seed

import "gorm.io/gorm"

// 存放所有 Seeder
var seeders []Seeder

var orderedSeederNames []string

type SeederFunc func(db *gorm.DB)

type Seeder struct {
	Func SeederFunc
	Name string
}

func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Func: fn,
		Name: name,
	})
}

func SetRunOrder(names []string) {
	orderedSeederNames = names
}
