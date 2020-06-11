package main

import (
    "encoding/json"
    "os"
    "fmt"
    "log"
    "strconv"
    "io/ioutil"
    "net/http"
    "math/rand"

    "github.com/gorilla/mux"
)

type User struct {
    Id        string   `json:"Id"`
    DeviceId  string   `json:"DeviceId"`
    History   []string `json:"history"`
}

type Thumbnail struct {
    Thumbnail string `json:"thumbnail"`
    Height    int    `json:"height"`
    Width     int    `json:"width"`
}

type Picture struct {
    Title       string    `json:"title"`
    Thumbnail   Thumbnail `json:"thumbnail"`
    Created_utc int       `json:"created_utc"`
    Author      string    `json:"author"`
    Id          string    `json:"id"`
    Ups         int       `json:"ups"`
    Downs       int       `json:"downs"`
    Media       string    `json:"media"`
}

type PicRequest struct {
    UserId string `json:"UserId"`
    PicId  string `json:"PicId"`
}

type Pictures struct {
    Base []Picture `json:"db"`
}

type Users struct {
    Base []User `json:"db"`
}

var UserDB Users
var PictureDB Pictures

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: homePage")
    json.NewEncoder(w).Encode(UserDB.Base)
}

func returnMemeById(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var request PicRequest
    var cur_user *User
    json.Unmarshal(reqBody, &request)
    idx := 0
    for pic_idx, picture := range PictureDB.Base {
        if picture.Id == request.PicId {
            //json.NewEncoder(w).Encode(picture)
            idx = pic_idx
        }
    }
    for i_user, user := range UserDB.Base {
        if user.Id == request.UserId {
            cur_user = &UserDB.Base[i_user]
        }
    }

    if cur_user != nil {
        fmt.Println("User established")
        fmt.Printf("Id: %s\n", cur_user.Id)
        fmt.Printf("History: %s\n", cur_user.History)
    } else {
        fmt.Println("User is not established")
    }

    if cur_user != nil {
        cur_user.History = append(cur_user.History, PictureDB.Base[idx].Id)
    }

    json.NewEncoder(w).Encode(PictureDB.Base[idx])
}

func returnMemeRandom(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var rec_request PicRequest
    var cur_user *User

    json.Unmarshal(reqBody, &rec_request)
    for i_user, user := range UserDB.Base {
        if user.Id == rec_request.UserId {
            cur_user = &UserDB.Base[i_user]
        }
    }

    if cur_user != nil {
        fmt.Println("User established")
        fmt.Printf("Id: %s\n", cur_user.Id)
        fmt.Printf("History: %s\n", cur_user.History)
    } else {
        fmt.Println("User is not established")
    }

    idx := rand.Intn(len(PictureDB.Base))
    if cur_user != nil {
        cur_user.History = append(cur_user.History, PictureDB.Base[idx].Id)
    }
    json.NewEncoder(w).Encode(PictureDB.Base[idx])
}

func returnMemeRecomended(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var rec_request PicRequest
    var cur_user *User
    json.Unmarshal(reqBody, &rec_request)
    for i_user, user := range UserDB.Base {
        if user.Id == rec_request.UserId {
            cur_user = &UserDB.Base[i_user]
        }
    }

    if cur_user != nil {
        fmt.Println("User established. Generating recomendation...")
        fmt.Printf("Id: %s\n", cur_user.Id)
        fmt.Printf("History: %s\n", cur_user.History)
    } else {
        fmt.Println("User is not established. Generating random pic...")
    }

    idx := rand.Intn(len(PictureDB.Base))
    if cur_user != nil {
        cur_user.History = append(cur_user.History, PictureDB.Base[idx].Id)
    }
    json.NewEncoder(w).Encode(PictureDB.Base[idx])
}

func authorize(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var user User
    json.Unmarshal(reqBody, &user)

    if len(user.DeviceId) > 0 {
        known_user := false
        for _, b_user := range UserDB.Base {
            if b_user.DeviceId == user.DeviceId {
                known_user = true
                user = b_user
                break
            }
        }
        if !known_user { // register
            user.Id = strconv.Itoa(len(UserDB.Base))
            UserDB.Base = append(UserDB.Base, user)
            fmt.Println("New user!")
        } else {
            fmt.Println("Known user!")
        }
    }

    json.NewEncoder(w).Encode(user)
}

func saveUserDB(w http.ResponseWriter, r *http.Request) {
    jsonUserDB, _ := json.Marshal(UserDB)
    ioutil.WriteFile("db/users_state.json", jsonUserDB, 0644)
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/auth", authorize).Methods("POST")
    myRouter.HandleFunc("/meme-id", returnMemeById).Methods("POST")
    myRouter.HandleFunc("/meme-random", returnMemeRandom).Methods("POST")
    myRouter.HandleFunc("/meme-rec", returnMemeRecomended).Methods("POST")
    myRouter.HandleFunc("/save-user-db", saveUserDB)
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "5000"
    }

    log.Fatal(http.ListenAndServe(":" + port, myRouter))
}

func main() {
    jsonFile, err := os.Open("db/pictures.json")
    if err != nil {
        fmt.Println(err)
    }
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &PictureDB)
    jsonFile.Close()

    jsonFile, err = os.Open("db/users_state.json")
    if err != nil {
        fmt.Println("users_state.json was not found. Booting new one..")
        jsonFile, err = os.Open("db/users.json")
    }
    byteValue, _ = ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &UserDB)
    jsonFile.Close()
    handleRequests()
}
