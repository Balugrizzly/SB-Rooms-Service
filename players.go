package main

import "sync"

type Players *[]Player

type Player struct {
	Mutex       *sync.Mutex
	Active      bool
	Name        string
	Matrix      Matrix
	LinSent     int
	LinReceived int
	LinBlocked  int
	HighestCmb  int
	LatestCmb   int
}
