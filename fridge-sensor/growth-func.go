package main

import (
	"math"
	"math/rand"
	)

//exponential decay
func growGrohe(stock *int) {
	oldStock := float64(*stock)
	*stock = int(math.Pow(oldStock, 0.95))
}

//barely ever changes
func growPfungstaedter(stock *int) {
	n := rand.Intn(50)
	switch {
	case n < 2: *stock--
	case n > 49: *stock++
	}
}

