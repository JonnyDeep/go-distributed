package grades

func init(){
	students = []Student{
		{
			ID: 1,
			FirstName: "Nick",
			LastName: "Carter",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type: GradeQuiz,
					Score: 85,
				},
				{
					Title: "Quiz 2",
					Type: GradeQuiz,
					Score: 90,
				},
				{
					Title: "Final Exam",
					Type: GradeExam,
					Score: 94,
				},
			},
		},
		{
			ID: 2,
			FirstName: "Roberto",
			LastName: "Baggio",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type: GradeQuiz,
					Score: 100,
				},
				{
					Title: "Quiz 2",
					Type: GradeQuiz,
					Score: 90,
				},
				{
					Title: "Final Exam",
					Type: GradeExam,
					Score: 100,
				},
			},
		}
	}
}