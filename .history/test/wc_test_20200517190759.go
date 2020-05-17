package main

import (
	"fmt"
	"math"
	"testing"
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

	fmt.Println(append(files,["333","4444"]))
}
