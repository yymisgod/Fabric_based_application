package controller

import "github.com/hyperledger/fabric-samples/modelworker/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	IsAdmin   string
}

var users []User

func init() {

	admin1 := User{LoginName: "LiFeng", Password: "123456", IsAdmin: "T"}
	admin2 := User{LoginName: "LiFeng2", Password: "123456", IsAdmin: "T"}
	admin3 := User{LoginName: "lifeng", Password: "1", IsAdmin: "T"}
	admin4 := User{LoginName: "1", Password: "1", IsAdmin: "T"}
	alice := User{LoginName: "alice", Password: "123456", IsAdmin: "F"}
	bob := User{LoginName: "bob", Password: "123456", IsAdmin: "F"}

	users = append(users, admin1)
	users = append(users, admin2)
	users = append(users, admin3)
	users = append(users, admin4)
	users = append(users, alice)
	users = append(users, bob)

}

func isAdmin(cuser User) bool {
	if cuser.IsAdmin == "T" {
		return true
	}
	return false
}
