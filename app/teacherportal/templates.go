package teacherportal

import (
	"html/template"
	"log"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	log.Println("templates rendered")
	rootTemplate, err = template.ParseFiles(
		"teacherportal/students.gohtml",
		"teacherportal/student.gohtml",
		"teacherportal/grades.gohtml",
	)
	if err != nil {
		return err
	}
	return nil
}
