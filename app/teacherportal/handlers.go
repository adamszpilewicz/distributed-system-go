package teacherportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/grades"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	log.Println("len of pathSegments: ", len(pathSegments))
	switch len(pathSegments) {
	case 2: // /students
		sh.renderStudents(w, r)
	case 3: // /students/{:id}
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.renderStudent(w, r, id)
	case 4: // /students/{:id}/grades
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.ToLower(pathSegments[3]) != "grades" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.renderGrades(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func (sh studentsHandler) renderStudents(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error retrieving students: ", err)
		}
	}()

	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		return
	}
	res, err := http.Get(serviceURL + "/students")
	if err != nil {
		return
	}
	var s grades.Students
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("students.gohtml").Execute(w, s)
	if err != nil {
		return
	}

}

func (sh studentsHandler) renderStudent(w http.ResponseWriter, r *http.Request, id int) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error retrieving students: ", err)
		}
	}()
	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		return
	}
	res, err := http.Get(fmt.Sprintf("%v/students/%v", serviceURL, id))
	if err != nil {
		return
	}
	var s grades.Student
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("student.gohtml").Execute(w, s)
	if err != nil {
		return
	}

}

func (studentsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", fmt.Sprintf("/students/%v", id))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
	title := r.FormValue("Title")
	gradeType := r.FormValue("Type")
	score, err := strconv.ParseFloat(r.FormValue("Score"), 32)
	if err != nil {
		log.Println("Failed to parse score: ", err)
		return
	}
	g := grades.Grade{
		Title: title,
		Type:  grades.GradeType(gradeType),
		Score: float32(score),
	}
	data, err := json.Marshal(g)
	if err != nil {
		log.Println("Failed to convert grade to JSON: ", g, err)
	}

	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		log.Println("Failed to retrieve instance of Grading Service", err)
		return
	}
	res, err := http.Post(fmt.Sprintf("%v/students/%v/grades", serviceURL, id), "application/json", bytes.NewBuffer(data))
	log.Println("response from post: ", res)
	log.Println(bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to save grade to Grading Service", err)
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.Println("Failed to save grade to Grading Service. Status: ", res.StatusCode)
		return
	}
}

//
//func (sh studentsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {
//	if r.Method != http.MethodPost {
//		w.WriteHeader(http.StatusMethodNotAllowed)
//		return
//	}
//	defer func() {
//		w.Header().Add("location", fmt.Sprintf("/students/%v", id))
//		w.WriteHeader(http.StatusTemporaryRedirect)
//	}()
//	title := r.FormValue("Title")
//	gradeType := r.FormValue("Type")
//	score, err := strconv.ParseFloat(r.FormValue("Score"), 32)
//	if err != nil {
//		log.Println("failed to parse score: ", err)
//	}
//	g := grades.Grade{
//		Title: title,
//		Type:  grades.GradeType(gradeType),
//		Score: float32(score),
//	}
//	log.Println("marshalling data")
//	data, err := json.Marshal(g)
//	if err != nil {
//		log.Println("failed to convert grade to JSON", g, err)
//		return
//	}
//
//	log.Println("printing marshalled data: ", data)
//	serviceURL, err := registry.GetProvider(registry.GradingService)
//	if err != nil {
//		log.Println("failed to receive instance og Grading Service", err)
//		return
//	}
//	res, err := http.Post(fmt.Sprintf("%v/students/%v/grades", serviceURL, id),
//		"application/json", bytes.NewBuffer(data))
//	if err != nil {
//		log.Println("Failed to save grade to Garding System", err)
//		return
//	}
//	if res.StatusCode != http.StatusCreated {
//		log.Println("failed to save grade to Grading System, Status: ",
//			res.StatusCode)
//		return
//	}
//}

func RegisterHandlers() {
	log.Println("registering handlers")
	http.Handle("/", http.RedirectHandler("/students",
		http.StatusPermanentRedirect))
	h := new(studentsHandler)

	http.Handle("/students", h)
	http.Handle("/students/", h)
}
