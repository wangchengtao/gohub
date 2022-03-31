package factories

import (
	"github.com/bxcodec/faker/v3"
	"gohub/app/models/category"
)

func MakeCategories(times int) []category.Category {
	var objs []category.Category

	faker.SetGenerateUniqueValues(true)

	for i := 0; i < times; i++ {
		model := category.Category{
			Name:        faker.Username(),
			Description: faker.Sentence(),
		}

		objs = append(objs, model)
	}

	return objs
}
