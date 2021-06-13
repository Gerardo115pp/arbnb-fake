package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	port     string
	host     string
	router   *Router
	state    *State
	sessions map[uint64]*User
	ok       []byte
}

func (self *Server) createSession(user *User, request *http.Request) uint64 {
	var session_key uint64 = shaAsInt64(fmt.Sprintf("%s:%s", user.Username, request.RemoteAddr))
	self.sessions[session_key] = user
	return session_key
}

func (self *Server) createUserFromRequest(request *http.Request) *User {
	var new_user *User = new(User)
	new_user.Id = self.state.getNewUserId()
	new_user.Username = request.FormValue("username")
	new_user.Name = request.FormValue("name")
	new_user.Phone = request.FormValue("phone")
	new_user.Email = request.FormValue("email")
	new_user.Password = shaAsInt64(request.FormValue("password"))
	new_user.Is_manager = false
	return new_user
}

func (self *Server) composeJson(key string, value string) []byte {
	return []byte(fmt.Sprintf("{\"%s\": %s}", key, value))
}

func (self *Server) composeResponse(response_value string) []byte {
	return self.composeJson("response", fmt.Sprintf("\"%s\"", response_value))
}

func (self *Server) composeError(error_value string) []byte {
	return self.composeJson("error", fmt.Sprintf("\"%s\"", error_value))
}

func (self *Server) composeHost() string {
	return fmt.Sprintf("%s:%s", self.host, self.port)
}

func (self *Server) enableCors(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Headers", "X-sk, Content-Type")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE, CONNECT")
		if request.Method == http.MethodOptions {
			response.WriteHeader(200)
			response.Write(self.ok)
			return
		}

		handler(response, request)
	}
}

func (self *Server) greet(response http.ResponseWriter, request *http.Request) {
	var data_to_hash string = request.URL.Query().Get("p")
	var greet_dialog = "hey"
	if data_to_hash != "" {
		greet_dialog = fmt.Sprintf("%d", shaAsInt64(data_to_hash))
	}
	response.WriteHeader(200)
	fmt.Fprintf(response, greet_dialog)
}

func (self *Server) getSessionKey(request *http.Request) (uint64, error) {
	var sk string = request.Header.Get("X-sk")
	if sk != "" {
		var session_key uint64 = stringToUint64(sk)
		return session_key, nil
	} else {
		return 0, fmt.Errorf("Missing X-sk header")
	}
}

func (self *Server) getUserFromRequest(request *http.Request) (*User, error) {
	var sk string = request.Header.Get("X-sk")
	if sk != "" {
		var session_key uint64 = stringToUint64(sk)
		if user, exists := self.sessions[session_key]; exists {
			return user, nil
		} else {
			return nil, fmt.Errorf("Users doesnt exists")
		}
	} else {
		return nil, fmt.Errorf("The user your looking for is not where you think it is...")
	}
}

func (self *Server) getRequestPostFormAsMap(request *http.Request) map[string]string {
	var form_items map[string]string = make(map[string]string)
	if len(request.Form) == 0 {
		request.ParseMultipartForm(32 << 20)
	}

	for key := range request.Form {
		form_items[key] = request.FormValue(key)
	}

	return form_items
}

func (self *Server) handleRegister(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		var new_user *User = self.createUserFromRequest(request)
		fmt.Printf("New user '%s'\n", new_user.Username)
		err := self.state.insertUser(new_user)
		if err == nil {
			response.WriteHeader(200)
			response.Write(self.ok)
		} else {
			fmt.Println("Username already exists")
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError(err.Error()))
		}
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError(fmt.Sprintf("Method '%s' is not allowed", request.Method)))
	}
}

func (self *Server) handleRooms(response http.ResponseWriter, request *http.Request) {
	var vendor *User
	vendor, _ = self.getUserFromRequest(request)

	switch request.Method {
	case http.MethodGet:
		var user_selector string = request.FormValue("guest")
		if user_selector == "*" || user_selector == "a" {
			response.WriteHeader(200)
			fmt.Fprint(response, self.state.getAllRooms(user_selector == "a"))
		} else if user_selector != "" {
			var vendor_products string = self.state.getUserReservations(vendor.Id)
			response.WriteHeader(200)
			fmt.Fprint(response, vendor_products)
		}
	case http.MethodPatch:
		var form_items map[string]string = self.getRequestPostFormAsMap(request)
		if len(form_items) > 0 {
			if err := self.state.makeReservation(vendor, &form_items); err == nil {
				response.WriteHeader(200)
				response.Write(self.ok)
			} else {
				response.WriteHeader(400)
				response.Write(self.composeError(err.Error()))
			}
		} else {
			fmt.Printf("error on form_items: %v\n", form_items)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError("missing parameters"))
		}
	case http.MethodDelete:
		var room_param string = request.FormValue("room")
		if room_param != "" {
			var room_id uint = stringToUint(room_param)
			var room *Room = self.state.getRoomByRoomNumber(room_id)
			if room != nil {
				if room.Used_by == int(vendor.Id) {
					room.Used_by = -1
					room.Reserved_from = ""
					room.Reserved_to = ""

					response.WriteHeader(200)
					response.Write(self.ok)
				} else {
					fmt.Printf("Room %d guest is not %d but %d\n", room.Number, vendor.Id, room.Used_by)
					response.WriteHeader(451)
					response.Write(self.composeError("fuck you dude, im not going to jail for you stupid shit"))
				}
			} else {
				response.WriteHeader(404)
				response.Write(self.composeError("Room not found"))
			}
		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("Missing room parameter"))
		}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("not allowed"))
	}
}

func (self *Server) handleBudgets(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var class_param, days_param string
		var budget_json []byte

		class_param = request.URL.Query().Get("class")
		days_param = request.URL.Query().Get("days")

		if class_param != "" && days_param != "" {
			days := stringToInt(days_param)

			budget := &struct {
				Cost_per_day int `json:"costs_per_day"`
				Total        int `json:"total"`
			}{}

			price_per_day, err := self.state.getClassCost(class_param)

			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				response.Write(self.composeError("wrong class name"))
			}

			budget.Cost_per_day = price_per_day
			budget.Total = price_per_day * days

			budget_json, err = json.Marshal(budget)
			if err != nil {
				response.WriteHeader(500)
				response.Write(self.composeError("we fucked up... sorry"))
				logFatal(err)
			}

			response.WriteHeader(200)
			response.Write(budget_json)
		} else {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError("Missing parameters"))
		}
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) handleUser(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		var guest_id_param string = request.FormValue("guest")
		var guest_data []byte
		var err error
		if guest_id_param == "*" {
			response.WriteHeader(200)
			fmt.Fprint(response, self.state.getAllUsers())
			return
		} else if guest_id_param != "" {
			guest_data, err = self.state.getUsersAndReservations(stringToUint(guest_id_param))
			if err != nil {
				fmt.Println("Error:", err.Error())
				response.WriteHeader(500)
				response.Write(self.composeError("server error, sorry for the inconvinience"))
				return
			}
		} else {
			var user *User
			user, _ = self.getUserFromRequest(request)
			guest_data, err = self.state.getUsersAndReservations(user.Id)
			if err != nil {
				fmt.Println("Error:", err.Error())
				response.WriteHeader(500)
				response.Write(self.composeError("server error, sorry for the inconvinience"))
				return
			}

		}
		response.WriteHeader(200)
		response.Write(guest_data)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method '%s' not allowed", request.Method)
	}
}

func (self *Server) login(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var username string = request.URL.Query().Get("username")
		var passwd string = request.URL.Query().Get("password")
		if username != "" && passwd != "" {
			var target *User = self.state.getUserByUsername(username)
			if target != nil {
				if target.Password == shaAsInt64(passwd) {
					// correct login
					var session_key uint64 = self.createSession(target, request)
					fmt.Printf("Session created for user '%s' on '%s': %d\n", username, request.RemoteAddr, session_key)

					response.WriteHeader(200)
					response.Write([]byte(fmt.Sprintf("{\"response\": \"%d\", \"is_manager\": %t}", session_key, target.Is_manager)))
				} else {
					fmt.Println("wrong password", target.Password, "!==", shaAsInt64(passwd))
					//wrong password
					response.WriteHeader(http.StatusBadRequest)
					response.Write(self.composeError("\"wrong password\""))
				}
			} else {
				//wrong username
				response.WriteHeader(http.StatusNotFound)
				response.Write(self.composeError("\"user doesnt exists\""))
			}
		} else {
			// missing information
			fmt.Printf("Incomplete login credentials username='%s' password='%s'\n", username, passwd)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError("\"missing information\""))
		}
	}
}

func (self *Server) logout(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPatch {
		var session_key uint64
		session_key, _ = self.getSessionKey(request)
		delete(self.sessions, session_key)
		fmt.Printf("Session %d was finished\n", session_key)
		response.WriteHeader(200)
		response.Write(self.ok)
	} else {
		fmt.Println("Wrong method, user was not logged out")
		response.WriteHeader(http.StatusNotModified)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) sessionExists(sk string) bool {
	var session_key uint64 = stringToUint64(sk)
	_, exists := self.sessions[session_key]
	return exists
}

func (self *Server) validateSession(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		var session_key string = request.Header.Get("X-sk")
		fmt.Println("Session:", session_key)
		if session_key != "" && self.sessionExists(session_key) {
			handler(response, request)
		} else {
			response.WriteHeader(http.StatusForbidden)
			response.Write(self.composeError("Invalid or missing session key"))
		}
	}
}

func (self *Server) run() {
	if err := self.state.loadState(); err != nil {
		logFatal(err)
	}

	fmt.Printf("\nState loaded!\n\nUsers: %d\nRooms: %d\nCosts: %d\nLast_user_id: %d\n\n", self.state.users.length, self.state.rooms.length, len(self.state.costs), self.state.last_user_id)

	self.router.registerRoute(NewRoute("/", true), self.enableCors(self.greet))
	self.router.registerRoute(NewRoute("/register", true), self.enableCors(self.handleRegister))
	self.router.registerRoute(NewRoute("/user", true), self.enableCors(self.validateSession(self.handleUser)))
	self.router.registerRoute(NewRoute("/reservations", true), self.enableCors(self.validateSession(self.handleRooms)))
	self.router.registerRoute(NewRoute("/budget", true), self.enableCors(self.handleBudgets))
	self.router.registerRoute(NewRoute(`/login`, true), self.enableCors(self.login))
	self.router.registerRoute(NewRoute(`/logout`, true), self.enableCors(self.validateSession(self.logout)))

	fmt.Println("Lisiting on '", self.composeHost(), "'")
	http.ListenAndServe(self.composeHost(), self.router)
}

func createServer(port int) *Server {
	var new_server *Server = new(Server)
	var server_port string = os.Getenv("GSERVER_PORT")
	if server_port == "" {
		server_port = "5006"
	}
	var server_host string = os.Getenv("GSERVER_HOST")
	if server_host == "" {
		server_host = "127.0.0.1"
	}

	new_server.router = createRouter()
	new_server.state = createState()
	new_server.sessions = make(map[uint64]*User)
	new_server.host = server_host
	new_server.port = server_port
	new_server.ok = new_server.composeResponse("ok")

	return new_server
}

func main() {
	var server *Server = createServer(5006)
	server.run()
}
