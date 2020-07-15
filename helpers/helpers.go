package helpers

import (
	"backend/models"
)

func UserSliceToMap(elements *models.Users) map[int]map[string]string {
	elementMap := map[int]map[string]string{}
	for i, user := range *elements {
		elementMap[i] = map[string]string{}
		elementMap[i]["id"] = user.UUID.String()
		elementMap[i]["userName"] = user.UserName
		elementMap[i]["firstName"] = user.FirstName
		elementMap[i]["lastName"] = user.LastName
		elementMap[i]["email"] = user.Email
		elementMap[i]["role"] = user.Role.Role
		elementMap[i]["isDeleting"] = "false"
	}
	return elementMap
}
