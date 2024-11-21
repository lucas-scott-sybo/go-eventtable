package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/go-eventsource/tutorial"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var userRepo *tutorial.Queries
var db *pgx.Conn

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserOut struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserCreatedEvent struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var cur CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&cur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db error"))
		fmt.Println(err.Error())
		return
	}
	defer tx.Rollback(ctx)

	repoTx := userRepo.WithTx(tx)
	user, err := repoTx.CreateUser(ctx, tutorial.CreateUserParams{
		Name:     cur.Name,
		Password: cur.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating user"))
		fmt.Println(err.Error())
		return
	}

	event := UserCreatedEvent{
		Id:   uint64(user.ID),
		Name: user.Name,
	}

	data, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling event"))
		fmt.Println(err.Error())
		return
	}

	_, err = repoTx.CreateEvent(ctx, tutorial.CreateEventParams{
		AggregateID: int32(user.ID),
		Kind:        "UserCreated",
		Version:     "v1",
		Data:        data,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating event"))
		fmt.Println(err.Error())
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed saving user"))
		fmt.Println(err.Error())
		return
	}

	out, err := json.Marshal(UserOut{
		Id:        uint64(user.ID),
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error making response"))
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

type UpdateUserRequest struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserUpdatedEvent struct {
	Id              uint64 `json:"id"`
	Name            string `json:"name"`
	PasswordChanged bool   `json:"password"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var cur UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&cur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db error"))
		fmt.Println(err.Error())
		return
	}
	defer tx.Rollback(ctx)

	repoTx := userRepo.WithTx(tx)
	user, err := repoTx.UpdateUser(ctx, tutorial.UpdateUserParams{
		ID:       int64(cur.Id),
		Name:     cur.Name,
		Password: cur.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating user"))
		fmt.Println(err.Error())
		return
	}

	event := UserUpdatedEvent{
		Id:              uint64(user.ID),
		Name:            user.Name,
		PasswordChanged: !(cur.Password == user.Password),
	}

	data, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling event"))
		fmt.Println(err.Error())
		return
	}

	_, err = repoTx.CreateEvent(ctx, tutorial.CreateEventParams{
		AggregateID: int32(user.ID),
		Kind:        "UserUpdated",
		Version:     "v1",
		Data:        data,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating event"))
		fmt.Println(err.Error())
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed saving user"))
		fmt.Println(err.Error())
		return
	}

	out, err := json.Marshal(UserOut{
		Id:        uint64(user.ID),
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error making response"))
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching users"))
		fmt.Println(err.Error())
		return
	}

	data := make([]UserOut, len(users))
	for i, u := range users {
		data[i] = UserOut{
			Id:        uint64(u.ID),
			Name:      u.Name,
			CreatedAt: u.CreatedAt.Time.String(),
			UpdatedAt: u.UpdatedAt.Time.String(),
		}
	}

	out, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling users"))
		fmt.Println(err.Error())
		return
	}

	w.Write(out)
}

func User(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CreateUser(w, r)
	case "PUT":
		UpdateUser(w, r)
	case "GET":
		GetUsers(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

type EventOut struct {
	Id        uint64 `json:"id"`
	Kind      string `json:"kind"`
	Version   string `json:"version"`
	CreatedAt string `json:"createdAt"`
	Data      any    `json:"data"`
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	events, err := userRepo.GetEventsFrom(ctx, tutorial.GetEventsFromParams{
		CreatedAt: pgtype.Timestamptz{Time: time.Now().Add(time.Duration(-1) * time.Hour), Valid: true},
		Limit:     20,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching users"))
		return
	}

	data := make([]EventOut, len(events))
	fmt.Printf("found %d events\n", len(events))
	for i, ev := range events {
		var o map[string]any
		err := json.Unmarshal(ev.Data, &o)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error unmarshaling event data"))
			return
		}

		data[i] = EventOut{
			Id:        uint64(ev.ID),
			Kind:      ev.Kind,
			Version:   ev.Version,
			CreatedAt: ev.CreatedAt.Time.String(),
			Data:      o,
		}
	}

	out, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling events"))
		return
	}

	w.Write(out)
}

func Events(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetEvents(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func main() {
	ctx := context.Background()
	var err error
	db, err = pgx.Connect(ctx, "user=app dbname=app password=dev-pass host=127.0.0.1")
	if err != nil {
		fmt.Printf("Error connecting to db %v\n", err)
	}
	defer db.Close(ctx)
	userRepo = tutorial.New(db)

	http.HandleFunc("/users", User)
	http.HandleFunc("/events", Events)
	http.ListenAndServe("0.0.0.0:5000", nil)
}
