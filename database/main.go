package database

import (
	"fmt"

	"github.com/emidiaz3/event-driven-server/models"
)

func GetUserByCorrelativeDB(userID int) (*models.User, string, error) {
	user, err := GetUserByCorrelativeMysql(userID)
	fmt.Print(err)

	if err == nil {
		return user, "Main", nil
	}

	user, err = GetUserByCorrelative(userID)
	if err == nil {
		return user, "Replica", nil
	}

	return nil, "", err
}

func GetUsersDB(ids []string) ([]models.User, string, error) {

	users, err := GetUsers(ids)
	if err == nil {
		return users, "Replica", nil
	}
	fmt.Println("hola", err)

	users, err = GetUsersMysql(ids)
	fmt.Println("adios", err)

	if err == nil {
		return users, "Main", nil
	}

	return nil, "", err
}
