package models

import (
	"database/sql"
	"time"
)

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
}

type TempPerson struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     string `json:"age"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
}

type Pizzeria struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PersonVisit struct {
	ID         int       `json:"id"`
	PersonID   int       `json:"person_id"`
	PizzeriaID int       `json:"pizzeria_id"`
	VisitDate  time.Time `json:"visit_date"`
}

type PersonVisitS struct {
	ID           int       `json:"id"`
	Person       string    `json:"name"`
	PizzeriaName string    `json:"pizzeria_name"`
	VisitDate    time.Time `json:"visit_date"`
}

type PersonOrder struct {
	ID        int       `json:"id"`
	PersonID  int       `json:"person_id"`
	MenuID    int       `json:"menu_id"`
	OrderDate time.Time `json:"order_date"`
}

type PersonOrderS struct {
	ID        int       `json:"id"`
	Person    string    `json:"name"`
	MenuID    int       `json:"menu_id"`
	OrderDate time.Time `json:"order_date"`
}

type Menu struct {
	ID         int           `json:"id"`
	PizzeriaID sql.NullInt64 `json:"pizzeria_id"`
	PizzaName  string        `json:"pizza_name"`
	Price      float64       `json:"price"`
}

type MenuS struct {
	ID        int     `json:"id"`
	Pizzeria  string  `json:"name"`
	PizzaName string  `json:"pizza_name"`
	Price     float64 `json:"price"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
