package grades

import "fmt"

const (
	GradeTest     = GradeType("Test")
	GradeHomework = GradeType("Homework")
	GradeQuiz     = GradeType("Quiz")
)

type GradeType string

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

type Students []Student

var students Students

func (s Students) GetById(ID int) (*Student, error) {
	for _, obj := range students {
		if obj.ID == ID {
			return &obj, nil
		}
	}
	return nil, fmt.Errorf("student with ID %v not found", ID)
}

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}
	result = result / float32(len(s.Grades))
	return result
}
