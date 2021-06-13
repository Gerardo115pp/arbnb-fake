package main

import (
	"encoding/json"
	"fmt"
)

type Room struct {
	Number        uint   `json:"number"`
	Used_by       int    `json:"used_by"`
	Class         string `json:"class"`
	Reserved_from string `json:"from"`
	Reserved_to   string `json:"to"`
}

func (self *Room) String() string {
	return self.toString()
}

func (self *Room) getId() uint {
	return self.Number
}

func (self *Room) compareString(other string) bool {
	return self.toString() == other
}

func (self *Room) toString() string {
	return fmt.Sprintf("room-%d", self.Number)
}

func (self *Room) toJson() []byte {
	json_user, err := json.Marshal(self)
	if err != nil {
		logFatal(err)
	}
	return json_user
}
