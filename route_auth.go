package main

import (
	"cardmanager/data"
	"encoding/json"
	"fmt"
	"net/http"
)

func setCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func preflight(w http.ResponseWriter, r *http.Request) (isOption bool) {
	isOption = r.Method == "OPTIONS"
	if isOption {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Length", "0")
		w.Header().Set("Content-Type", "text/plain")
	}
	return
}

func signupAccount(w http.ResponseWriter, r *http.Request) {

	isOption := preflight(w, r)
	fmt.Println(isOption)
	if isOption {
		return
	}

	setCors(w, r)

	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user data.User
	json.Unmarshal(body, &user)

	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	output, err := json.MarshalIndent(&user, "", "\t\t")
	if err != nil {
		return
	}
	w.Write(output)
	return

}

func authenticate(w http.ResponseWriter, r *http.Request) {
	isOption := preflight(w, r)
	if isOption {
		return
	}
	err := r.ParseForm()

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var requestUser data.User
	json.Unmarshal(body, &requestUser)

	user, err := data.UserByEmail(requestUser.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}

	if user.Password == data.Encrypt(requestUser.Password) {
		err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		session, err := user.Session()
		if err != nil {
			danger(err, "Cannot find session")
		}
		cookie := http.Cookie{
			Name:     "user_session",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}

	output, err := json.MarshalIndent(&user, "", "\t\t")
	if err != nil {
		return
	}
	setCors(w, r)
	w.Write(output)

}

func sessionCheck(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	valid := false
	if err != nil {
		fmt.Println("session ng")
	} else {
		fmt.Println("session ok!!")
		valid = true
	}
	m1 := map[string]interface{}{
		"valid": valid,
	}

	output, err := json.MarshalIndent(&m1, "", "\t\t")
	if err != nil {
		return
	}
	setCors(w, r)

	w.Write(output)

}
