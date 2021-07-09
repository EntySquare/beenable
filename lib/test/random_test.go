package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	GenerateRangeNum(10001, 20000)
	time.Sleep(time.Duration(2) * time.Second)
	GenerateRangeNum(10001, 20000)
	time.Sleep(time.Duration(2) * time.Second)
	GenerateRangeNum(10001, 20000)
}

// GenerateRangeNum 生成一个区间范围的随机数
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	fmt.Printf("rand is %v\n", randNum)
	fmt.Printf("rand is %v\n", int32(randNum))
	return randNum
}

func TestRange(t *testing.T) {
	a := []string{"Error", "Error", "Error", "Running"}
	var b []string
	for x, y := range a {
		if y == "Running" {
			//a = append(a[:x],a[x+1:]...)
			b = append(b, y)
		}
		fmt.Println(b, x)
	}
}
