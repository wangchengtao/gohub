package factories

import (
	"github.com/bxcodec/faker/v3"
	"gohub/app/models/topic"
)

func MakeTopics(times int) []topic.Topic {
	var objs []topic.Topic

	faker.SetGenerateUniqueValues(true)

	for i := 0; i < times; i++ {
		model := topic.Topic{
			Title:      faker.Sentence(),
			Body:       faker.Paragraph(),
			CategoryID: "3",
			UserID:     "1",
		}

		objs = append(objs, model)
	}

	return objs
}
