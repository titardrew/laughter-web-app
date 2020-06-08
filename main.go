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

// Article - Our struct for all articles
type Article struct {
    Id      string `json:"Id"`
    Title   string `json:"Title"`
    Desc    string `json:"desc"`
    Content string `json:"content"`
}

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

type Pictures struct {
    Base []Picture `json:"db"`
}

var Articles []Article
var Users []User
var PictureDB Pictures

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: homePage")
    json.NewEncoder(w).Encode(PictureDB)
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    for _, article := range Articles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
        }
    }
}

func returnMemeById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    for _, picture := range PictureDB.Base {
        if picture.Id == key {
            json.NewEncoder(w).Encode(picture)
        }
    }
}

func returnMemeRandom(w http.ResponseWriter, r *http.Request) {
    idx := rand.Intn(len(PictureDB.Base))
    json.NewEncoder(w).Encode(PictureDB.Base[idx])
}

func authorize(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var user User
    json.Unmarshal(reqBody, &user)

    if len(user.DeviceId) > 0 {
        known_user := false
        for _, b_user := range Users {
            if b_user.DeviceId == user.DeviceId {
                known_user = true
                user = b_user
                break
            }
        }
        if !known_user {
            user.Id = strconv.Itoa(len(Users))
            Users = append(Users, user)
            fmt.Println("New user!")
        } else {
            fmt.Println("Known user!")
        }
    }

    json.NewEncoder(w).Encode(Users)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    // get the body of our POST request
    // unmarshal this into a new Article struct
    // append this to our Articles array.    
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article
    json.Unmarshal(reqBody, &article)
    // update our global Articles array to include
    // our new Article
    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    for index, article := range Articles {
        if article.Id == id {
            Articles = append(Articles[:index], Articles[index+1:]...)
        }
    }

}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles)
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
    myRouter.HandleFunc("/auth", authorize).Methods("POST")
    myRouter.HandleFunc("/meme/{id}", returnMemeById).Methods("Get")
    myRouter.HandleFunc("/meme-random", returnMemeRandom)
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
    myRouter.HandleFunc("/article/{id}", returnSingleArticle)
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "5000"
    }

    log.Fatal(http.ListenAndServe(":" + port, myRouter))
}

func main() {
    Articles = []Article{
        Article{Id: "1", Title: "Hello",
                Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2",
                Desc: "Article Description", Content: "Article Content"},
    }
    Users = []User{
        User{Id: "0", DeviceId: "0x00013", History: []string{"138", "124"}},
        User{Id: "1", DeviceId: "0x00014"},
    }

    jsonFile, err := os.Open("db/pictures.json")
    if err != nil {
        fmt.Println(err)
    }
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &PictureDB)
    jsonFile.Close()

    handleRequests()
}
