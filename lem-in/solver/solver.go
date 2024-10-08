package solver

import (
	"fmt"
	"lem-in/antfarm"
)

// Make the ant 'a' move and print the movement if its possible
func Solve(n int, f antfarm.Antfarm, t *string) {
	a := f.Ants[n]
	if a.Room == f.End {
	} else {
		if LinkToEnd(a, f) && *t == "false" {
			if a.Room == f.Start {
				*t = "true"
			}
			Move(a, f.End, f)
		} else {
			for _, room := range a.Room.Links {
				if room.Empty && !isVisited(a, room) {
					if a.Room == f.Start && room == f.End {
						continue
					}
					Move(a, room, f)
					break
				}
			}

		}
	}
}

func Move(a *antfarm.Ant, room *antfarm.Room, f antfarm.Antfarm) {
	prev := a.Room
	a.Room = room
	if prev != f.Start {
		prev.Empty = true
	}
	if room != f.End {
		room.Empty = false
	}
	a.Path = append(a.Path, a.Room)
	fmt.Printf("L%d-%s ", a.Name, room.Name)
}

// Verify if the actual Room of the ant 'a' is linked to the end
func LinkToEnd(a *antfarm.Ant, f antfarm.Antfarm) bool {
	for _, room := range a.Room.Links {
		if room == f.End {
			return true
		}
	}
	return false
}

func AllinEnd(f antfarm.Antfarm) bool {
	for _, ants := range f.Ants {
		if ants.Room != f.End {
			return false
		}
	}
	return true
}

func isVisited(a *antfarm.Ant, r *antfarm.Room) bool {
	for _, room := range a.Path {
		if room == r {
			return true
		}
	}
	return false
}

func AllPath(f antfarm.Antfarm) [][]*antfarm.Room {
	tab := [][]*antfarm.Room{}
	str := []*antfarm.Room{f.Start}
	a := &antfarm.Ant{Name: 0,Room: f.Start}
	for _, room := range f.Start.Links {
		str = append(str, room)
		Move(a,room,f)
		if room == f.End {
			tab = append(tab, str)
			str = []*antfarm.Room{f.Start}
		}
		for a.Room != f.End {
			if LinkToEnd(a, f){
				Move(a, f.End, f)
				str = append(str,f.End)
				tab = append(tab, str)
				str = []*antfarm.Room{f.Start}
			} else {
				for _, room := range a.Room.Links {
					if room.Empty && !isVisited(a, room) {
						Move(a, room, f)
						str = append(str, room)
						break
					}
				}
			}
		}
	}
	return tab
}

func FollowPath(ant *antfarm.Ant ,path []*antfarm.Room,f antfarm.Antfarm) {
	for _, room := range path {
		if ant.Room != room {
			Move(ant,room,f)
		}
	}
}