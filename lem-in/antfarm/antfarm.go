package antfarm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name  string
	Coord []int
	Links []*Room
	Empty bool
}

type Ant struct {
	Name int
	Room *Room
	Path []*Room
}

type Antfarm struct {
	Ants  []*Ant
	Rooms map[string]*Room
	Start *Room
	End   *Room
}

// Make the AntFarm with the data from the file
func MakeFarm(file string) Antfarm {
	Antfarm := Antfarm{Rooms: make(map[string]*Room)}
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("ERROR: Problems reading the file", err)
	}
	text := strings.Split(string(data), "\n")
	for line := 1; line < len(text); line++ {
		if text[line] == "##start" {
			room := makeRoom(text[line+1])
			room.Empty = false
			Antfarm.Rooms[room.Name] = &room
			Antfarm.Start = &room
			line++
			continue
		}
		if text[line] == "##end" {
			room := makeRoom(text[line+1])
			Antfarm.Rooms[room.Name] = &room
			Antfarm.End = &room
			line++
			continue
		}
		if strings.Contains(text[line], "-") {
			Link(text[line], Antfarm)
			continue
		}
		room := makeRoom(text[line])
		Antfarm.Rooms[room.Name] = &room
	}
	Antfarm.Ants = Ants(text[0], Antfarm.Start)
	return Antfarm
}

// Make a Room with the data given in the string 's
func makeRoom(s string) Room {
	name := strings.Split(s, " ")[0]
	x, _ := strconv.Atoi(strings.Split(s, " ")[1])
	y, _ := strconv.Atoi(strings.Split(s, " ")[2])
	return Room{
		Name:  name,
		Coord: []int{x, y},
		Empty: true,
	}
}

// Link two Room between us
func Link(s string, f Antfarm) {
	room1 := f.Rooms[strings.Split(s, "-")[0]]
	room2 := f.Rooms[strings.Split(s, "-")[1]]
	room1.Links = append(room1.Links, room2)
	room2.Links = append(room2.Links, room1)
}

// Make a slice with all ants
func Ants(s string, Start *Room) []*Ant {
	nb_ants, _ := strconv.Atoi(s)
	tab_a := []*Ant{}
	for i := 1; i <= nb_ants; i++ {
		a := Ant{Name: i, Room: Start, Path: []*Room{Start}}
		tab_a = append(tab_a, &a)
	}
	return tab_a
}
