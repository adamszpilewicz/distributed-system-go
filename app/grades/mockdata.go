package grades

func init() {
	students = []Student{
		Student{
			ID:        1,
			FirstName: "Adam",
			LastName:  "Szpilewicz",
			Grades:    []Grade{Grade{Title: "Quiz 1", Type: GradeQuiz, Score: 85}},
		},

		Student{
			ID:        2,
			FirstName: "Renata",
			LastName:  "Szpilewicz",
			Grades:    []Grade{Grade{Title: "Quiz 1", Type: GradeQuiz, Score: 82}},
		},
	}
}
