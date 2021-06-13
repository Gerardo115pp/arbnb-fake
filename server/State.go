package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const STORE_FILE = "./store.json"

type State struct {
	users        *List
	rooms        *List
	costs        map[string]int
	last_user_id uint
}

func (self *State) clearState() {
	self.users.clear()
}

func (self *State) getAllRooms(only_avaliable bool) string {
	var result string
	if only_avaliable {
		result = self.rooms.filter(func(c Content) bool { return c.(*Room).Used_by == -1 }).toJson()
	} else {
		result = self.rooms.toJson()
	}
	return result
}

func (self *State) getAllUsers() string {
	return self.users.toJson()
}

func (self *State) getClassCost(class_name string) (int, error) {
	if class_price, exists := self.costs[class_name]; exists {
		return class_price, nil
	} else {
		return -1, fmt.Errorf("class name doesnt exists")
	}
}

func (self *State) getUsersAndReservations(user_id uint) ([]byte, error) {
	var target_user Content = self.users.searchBy(func(c Content) bool { return c.getId() == user_id })
	if target_user != nil {

		// string with the form [...reservations]
		var user_reservations []Room
		err := json.Unmarshal([]byte(self.rooms.filter(func(c Content) bool { return c.(*Room).Used_by == int(user_id) }).toJson()), &user_reservations)
		if err != nil {
			logFatal(err)
		}
		user_details := &struct {
			User         *User  `json:"user"`
			Reservations []Room `json:"reservations"`
		}{}
		user_details.Reservations = user_reservations
		user_details.User = target_user.(*User)

		user_json, err := json.Marshal(user_details)
		if err != nil {
			logFatal(err)
		}
		return user_json, nil

	} else {
		return nil, fmt.Errorf("user doesnt exists")
	}
}

func (self *State) getNewUserId() uint {
	self.last_user_id++
	return self.last_user_id
}

func (self *State) getRoomByRoomNumber(room_number uint) *Room {
	var search_result Content = self.rooms.searchBy(func(c Content) bool { return c.(*Room).Number == room_number })
	if search_result != nil {
		return search_result.(*Room)
	} else {
		return nil
	}
}

func (self *State) getUserReservations(guest_id uint) string {
	return self.rooms.filter(func(c Content) bool { return c.getId() == guest_id }).toJson()
}

func (self *State) getUserByUsername(username string) *User {
	var search_result Content = self.users.searchBy(func(c Content) bool { return c.(*User).Username == username })
	if search_result != nil {
		return search_result.(*User)
	} else {
		return nil
	}
}

func (self *State) getUserById(target_id uint) *User {
	return self.users.searchBy(func(c Content) bool { return c.getId() == target_id }).(*User)
}

func (self *State) insertUser(new_user *User) error {
	if other := self.users.exists(new_user); other == nil {
		self.users.append(new_user)
		return self.save()
	} else {
		return fmt.Errorf("User '%s' already exists", new_user.toString())
	}
}

func (self *State) makeReservation(user *User, reservation_details *map[string]string) error {
	var requested_room *Room = self.getRoomByRoomNumber(stringToUint((*reservation_details)["room_id"]))
	if requested_room != nil {
		requested_room.Used_by = int(user.Id)
		requested_room.Reserved_from = (*reservation_details)["from"]
		requested_room.Reserved_to = (*reservation_details)["to"]
		return self.save()
	} else {
		return fmt.Errorf("room %s dosnt exists", (*reservation_details)["room_id"])
	}
}

func (self *State) loadState() error {
	var err error
	var state_json []byte
	state_json, err = ioutil.ReadFile(STORE_FILE)
	if err != nil {
		logFatal(err)
	}

	store_file_struct := &struct {
		Users      []User         `json:"users"`
		Rooms      []Room         `json:"rooms"`
		Costs      map[string]int `json:"costs"`
		LastUserID int            `json:"last_user_id"`
	}{}

	err = json.Unmarshal(state_json, store_file_struct)

	self.loadUsers(store_file_struct.Users)
	self.loadRooms(store_file_struct.Rooms)

	self.last_user_id = uint(store_file_struct.LastUserID)
	self.costs = store_file_struct.Costs
	return err
}

func (self *State) loadUsers(users []User) {
	for h := range users {
		self.users.append(&(users[h]))
	}
}

func (self *State) loadRooms(rooms []Room) {
	for h := range rooms {
		self.rooms.append(&(rooms[h]))
	}
}

func (self *State) save() error {
	var costs_data []byte
	costs_data, err := json.Marshal(self.costs)
	if err != nil {
		fmt.Println("Error while saving state:", err.Error())
		logFatal(err)
	}

	var json_state string = fmt.Sprintf("{\"rooms\": %s,\"users\": %s, \"costs\": %s, \"last_user_id\": %d}", self.rooms.toJson(), self.users.toJson(), string(costs_data), self.last_user_id)
	return ioutil.WriteFile(STORE_FILE, []byte(json_state), 0777)
}

// End of state functions

func createState() *State {
	var new_state *State = new(State)
	new_state.last_user_id = 0
	new_state.users = new(List)
	new_state.rooms = new(List)
	return new_state
}
