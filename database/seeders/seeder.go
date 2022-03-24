package seeders

import "gohub/pkg/seed"

func Initialize() {
	seed.SetRunOrder([]string{
		"SeederUsersTable",
	})
}
