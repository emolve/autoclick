package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestAddRandomTime(t *testing.T) {
	for i := 0; i < 10; i++ {

		_, randomTime := AddRandomTime(time.Second * 360)
		fmt.Printf("%d:%s\n", i, randomTime)
		time.Sleep(time.Second * 2)
	}
	for i := 0; i < 10; i++ {

		_, randomTime := AddRandomTime(time.Second * 240)
		fmt.Printf("%d:%s\n", i, randomTime)
		time.Sleep(time.Second * 2)
	}

}
