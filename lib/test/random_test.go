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
	return randNum
}
