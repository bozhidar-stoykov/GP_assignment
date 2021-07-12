package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

// User entity
type User struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Users []User

// Mock data container
var users Users

func homePage(wr http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(wr, "-------------------------------- Home Page --------------------------------")
	fmt.Fprintln(wr, "Endpoints:")
	fmt.Fprintln(wr, "/users  			*GET - returns all users / POST - creates a new user (use raw json for the body)*")
	fmt.Fprintln(wr, "/users/{email}  		*GET - returns user by email / DELETE - deletes a user by email*")
	fmt.Fprintln(wr, "/users/{partial email}  	*GET - returns all users that match the partial email*")
}

// GET all users
func getAllUsers(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(users)
}

// Filter users with specific email
func filterUsersByEmail(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var filtered Users

	var filter string = params["emailFilter"]
	var paterns [3]string = getRegexPaterns(filter)
	var regex string = paterns[0] + "[@]" + paterns[1] + "[.]" + paterns[2]

	for _, item := range users {
		res, _ := regexp.MatchString(regex, item.Email)
		if res {
			filtered = append(filtered, item)
		}
	}

	json.NewEncoder(wr).Encode(filtered)
}

// Split filter into three parts: 1-before "@", 2-between "@" and ".", 3-after "."
func getRegexPaterns(filter string) (paterns [3]string) {
	paterns[0] = `[\w+]`
	paterns[1] = `([a-z]+)`
	paterns[2] = `([a-z]+)`
	var secondSplit = filter

	if strings.Contains(filter, "@") {
		split1 := strings.Split(filter, "@")
		secondSplit = split1[1]

		// Return regex for multiple symbols if empty
		if len(strings.TrimSpace(split1[0])) != 0 {
			paterns[0] = split1[0]
		}
	}

	if strings.Contains(filter, ".") {
		split2 := strings.Split(secondSplit, ".")

		if len(strings.TrimSpace(split2[0])) != 0 {
			paterns[1] = split2[0]
		}

		if len(strings.TrimSpace(split2[1])) != 0 {
			paterns[2] = split2[1]
		}
	} else {
		paterns[1] = secondSplit
	}

	return paterns
}

// POST new user
func createUser(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)

	// Return created user
	json.NewEncoder(wr).Encode(user)
}

// GET user by Email
func getUser(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range users {
		if item.Email == params["email"] {
			json.NewEncoder(wr).Encode(item)
			return
		}
	}

	json.NewEncoder(wr).Encode(&User{})
}

// DELETE user by Email
func deleteUser(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range users {
		if item.Email == params["email"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}

	// Return remaning users
	json.NewEncoder(wr).Encode(users)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", getAllUsers).Methods("GET")
	myRouter.HandleFunc("/users", createUser).Methods("POST")

	myRouter.HandleFunc("/users/{emailFilter}", filterUsersByEmail).Methods("GET")

	myRouter.HandleFunc("/users/{email}", getUser).Methods("GET")
	myRouter.HandleFunc("/users/{email}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	// Mock data
	users = Users{
		User{Email: "john@hotmail.com", Phone: "+318457447", Password: "011235813"},
		User{Email: "jane@gmail.com", Phone: "+319677758", Password: "11010010"},
		User{Email: "richard81@gmail.com", Phone: "+3598983650", Password: "24688642"},
		User{Email: "alexander@abv.bg", Phone: "+3598874255", Password: "RK35_mS!"},
		User{Email: "oliver@kangaroo.au", Phone: "+6145770931", Password: "Oka936Rt"},
	}

	handleRequests()
}
