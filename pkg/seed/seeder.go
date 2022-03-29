package seed

import (
	"gohub/pkg/console"
	"gohub/pkg/database"
	"gorm.io/gorm"
)

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

func GetSeeder(name string) Seeder {
	for _, seeder := range seeders {
		if name == seeder.Name {
			return seeder
		}
	}

	return Seeder{}
}

func RunAll() {
	// 先运行 ordered 的
	executed := make(map[string]string)
	for _, name := range orderedSeederNames {
		seeder := GetSeeder(name)

		if len(seeder.Name) > 0 {
			console.Warning("Running Ordered Seeder: " + seeder.Name)
			seeder.Func(database.DB)
			executed[name] = name
		}
	}

	// 再运行剩下的
	for _, seeder := range seeders {
		// 过滤已运行的
		if _, ok := executed[seeder.Name]; !ok {
			console.Warning("Running Seeder: " + seeder.Name)
			seeder.Func(database.DB)
		}
	}
}

func RunSeeder(name string) {
	seeder := GetSeeder(name)

	if len(seeder.Name) > 0 {
		seeder.Func(database.DB)
	} else {
		console.Error("Seeder not found: " + name)
	}
}
