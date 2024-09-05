package main

import (
	"fmt"
	"lem-in/antfarm"
	"lem-in/solver"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("[USAGE]: go run . $file")
		return
	}
	AntFarm := antfarm.MakeFarm(os.Args[1])
	for !solver.AllinEnd(AntFarm) {
		t := "false"
		for i := 0; i < len(AntFarm.Ants); i++ {
			solver.Solve(i, AntFarm, &t)
		}
		t = "false"
		fmt.Println()
	}
}
