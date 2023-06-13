package internal

import (
	"github.com/prometheus/client_golang/prometheus"
)

var CreateStudentCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "new_student",
	Help: "Creation of new student",
})

var GetStudentCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "getting_students",
	Help: "Getting info about student",
})

func init() {
	prometheus.MustRegister(CreateStudentCounter)
	prometheus.MustRegister(GetStudentCounter)
}
