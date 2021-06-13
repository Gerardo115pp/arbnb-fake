package main

import (
	"encoding/json"
)

type User struct {
	Id         uint   `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   uint64 `json:"password"`
	Is_manager bool   `json:"is_manager"`
}

func (self *User) String() string {
	return self.toString()
}

func (self *User) getId() uint {
	return self.Id
}

func (self *User) compareString(other string) bool {
	return self.toString() == other
}

func (self *User) toString() string {
	return self.Username
}

func (self *User) toJson() []byte {
	json_user, err := json.Marshal(self)
	if err != nil {
		logFatal(err)
	}
	return json_user
}
