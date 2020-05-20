package main

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

type Salary struct {
	base int
	num  int
}

type Employee struct {
	base int
	Salary
}

func (s *Salary) getSalary() int {
	//base不会重载啊
	return s.base * s.num
}

// func (e *Employee) getSalary() int {
// 	return e.base * e.num
// }

const (
	//UNASSIGN task not assign yet
	UNASSIGN = iota
	// ASSIGN have beean assign
	ASSIGN
	// SUCCESS success
	SUCCESS
	// FAILED failed
	FAILED
	// TIMEOUT timeout
	TIMEOUT
)

func TestMap(t *testing.T) {
	fmt.Println("efg")
	files := []string{"fakename1", "fakename2", "fakename3"}
	fmt.Println(files[1:3])
	fmt.Println(math.Min(1, 2))
	task := map[string]string{}
	task["b"] = "b"
	fmt.Println(task)
	for k, v := range task {
		fmt.Println(k, v)
	}

	bucket := make([][]string, 10)
	bucket[0] = []string{"a", "b"}
	fmt.Println(bucket)
	ftemp := []string{"333", "4444"}
	fmt.Println(append(files, ftemp...))
	fmt.Println("abc" + "def")
	fmt.Printf("abc %s", "def")

	salary := Salary{1, 100}
	employee := Employee{2, salary}
	fmt.Println(employee)
	fmt.Println(employee.getSalary())
	fmt.Println(employee.base)

	fmt.Println(UNASSIGN, ASSIGN, SUCCESS)
	twoprint()
	time.Sleep(time.Duration(5) * time.Second)
}

var a string
var once sync.Once

func setup() {
	println("init")
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
	go doprint()
	go doprint()
}
