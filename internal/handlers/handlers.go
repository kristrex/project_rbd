package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rbd_project/internal/auth"
	"rbd_project/internal/db"
	"rbd_project/internal/models"
	"strconv"

	_ "github.com/lib/pq"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		log.Println("Failed to insert user: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Failed to decode: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var storedUser models.User
	err = db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", user.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Invalid username or password: ", err)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			log.Println("Hz chto za oshibka: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !auth.CheckPasswordHash(user.Password, storedUser.Password) {
		log.Println("Invalid username or password in check hash", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		log.Println("Error generate jwt", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func GetPersons(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, age, gender, address FROM person")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var p models.Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.Address)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		persons = append(persons, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

func GetPizzerias(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM pizzeria")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pizzerias []models.Pizzeria
	for rows.Next() {
		var p models.Pizzeria
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pizzerias = append(pizzerias, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzerias)
}

func GetVisits(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT pv.id, p.name AS person_name, pi.name AS pizzeria_name, pv.visit_date FROM person_visits AS pv JOIN person AS p ON pv.person_id = p.id JOIN pizzeria AS pi ON pv.pizzeria_id = pi.id;")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var visits []models.PersonVisitS
	for rows.Next() {
		var v models.PersonVisitS
		err := rows.Scan(&v.ID, &v.Person, &v.PizzeriaName, &v.VisitDate)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		visits = append(visits, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visits)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT po.id, p.name AS person_name, menu_id, order_date FROM person_order AS po JOIN person AS p ON po.person_id = p.id")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.PersonOrderS
	for rows.Next() {
		var o models.PersonOrderS
		err := rows.Scan(&o.ID, &o.Person, &o.MenuID, &o.OrderDate)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetMenus(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, pizzeria_id, pizza_name, price FROM menu")
	if err != nil {
		log.Println("Failed to execute query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var m models.Menu
		err := rows.Scan(&m.ID, &m.PizzeriaID, &m.PizzaName, &m.Price)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		menus = append(menus, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menus)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var temp models.TempPerson
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		log.Println("Failed to decode: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование строки age в целое число
	age, err := strconv.Atoi(temp.Age)
	if err != nil {
		log.Println("Failed to convert age to int: ", err)
		http.Error(w, "Invalid age value", http.StatusBadRequest)
		return
	}

	// Создание структуры Person с преобразованным значением age
	person := models.Person{
		Name:    temp.Name,
		Age:     age,
		Gender:  temp.Gender,
		Address: temp.Address,
	}

	// Проверка на пустые поля
	if person.Name == "" || person.Gender == "" || person.Address == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO person (id, name, age, gender, address) SELECT COALESCE(MAX(id), 0) + 1, $1, $2, $3, $4 FROM person RETURNING id;`
	err = db.DB.QueryRow(query, person.Name, person.Age, person.Gender, person.Address).Scan(&person.ID)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// UpdatePerson обновляет существующего пользователя
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Invalid ID: ", err)
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	var temp models.TempPerson
	err = json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		log.Println("Failed to decode: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование строки age в целое число
	age, err := strconv.Atoi(temp.Age)
	if err != nil {
		log.Println("Failed to convert age to int: ", err)
		http.Error(w, "Invalid age value", http.StatusBadRequest)
		return
	}

	// Создание структуры Person с преобразованным значением age
	person := models.Person{
		Name:    temp.Name,
		Age:     age,
		Gender:  temp.Gender,
		Address: temp.Address,
	}

	// Проверка на пустые поля
	if person.Name == "" || person.Gender == "" || person.Address == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `UPDATE person SET name = $1, age = $2, gender = $3, address = $4 WHERE id = $5`
	_, err = db.DB.Exec(query, person.Name, person.Age, person.Gender, person.Address, id)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePerson удаляет пользователя
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Invalid ID: ", err)
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM person WHERE id = $1`
	_, err = db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SearchPersons(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	minAgeStr := r.URL.Query().Get("minAge")

	if city == "" || minAgeStr == "" {
		http.Error(w, "City and minAge are required", http.StatusBadRequest)
		return
	}

	minAge, err := strconv.Atoi(minAgeStr)
	if err != nil {
		http.Error(w, "Invalid minAge value", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, age, gender, address FROM person WHERE address LIKE $1 AND age >= $2`
	rows, err := db.DB.Query(query, "%"+city+"%", minAge)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.Name, &person.Age, &person.Gender, &person.Address)
		if err != nil {
			log.Println("Failed to scan row: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		persons = append(persons, person)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

func GetAverageAgeByCity(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	var avgAge sql.NullFloat64
	err := db.DB.QueryRow("SELECT get_average_age_by_city($1)", city).Scan(&avgAge)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if avgAge.Valid {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"average_age": %f}`, avgAge.Float64)
	} else {
		http.Error(w, "No data found", http.StatusNotFound)
	}
}

func DeletePersonsOlderThan(w http.ResponseWriter, r *http.Request) {
	maxAgeStr := r.URL.Query().Get("maxAge")
	if maxAgeStr == "" {
		http.Error(w, "maxAge is required", http.StatusBadRequest)
		return
	}

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		http.Error(w, "Invalid maxAge value", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("CALL delete_persons_older_than($1)", maxAge)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdatePersonWithTrigger(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	var temp models.TempPerson
	err = json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		log.Println("Failed to decode: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	age, err := strconv.Atoi(temp.Age)
	if err != nil {
		log.Println("Failed to convert age to int: ", err)
		http.Error(w, "Invalid age value", http.StatusBadRequest)
		return
	}

	person := models.Person{
		Name:    temp.Name,
		Age:     age,
		Gender:  temp.Gender,
		Address: temp.Address,
	}

	if person.Name == "" || person.Gender == "" || person.Address == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `UPDATE person SET name = $1, age = $2, gender = $3, address = $4 WHERE id = $5`
	_, err = db.DB.Exec(query, person.Name, person.Age, person.Gender, person.Address, id)
	if err != nil {
		log.Println("Failed query: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
