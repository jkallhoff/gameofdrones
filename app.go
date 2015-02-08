package main

import "fmt"
import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	// P: number of players in the game (2 to 4 players)
	// ID: ID of your player (0, 1, 2, or 3)
	// D: number of drones in each team (3 to 11)
	// Z: number of zones on the map (4 to 8)
	var P, ID, D, Z int
	fmt.Scan(&P, &ID, &D, &Z)

	game := SetupGame(P, ID, D, Z)

	fmt.Fprintf(os.Stderr, "%s", game.Me())

	for {
		game.LoadRound()

		firstZone := game.Zones[0]

		for _, _ = range game.Me().Drones {
			fmt.Printf("%d %d\n", firstZone.Center.X, firstZone.Center.Y)
		}
	}
}

//Types
type Point struct {
	X, Y int
}

type Zone struct {
	Center        *Point
	ControllerId  int
	MaxEnemyShips int
}

type Zones []*Zone

type Drone struct {
	Location *Point
}

type Drones []*Drone

type Player struct {
	Drones Drones
}

type Players []*Player

type Game struct {
	Players Players //The collection of players in the game.
	Zones   Zones
	MeId    int
}

func (game *Game) Me() *Player {
	return game.Players[game.MeId]
}

func SetupGame(p, id, d, z int) *Game {
	game := new(Game)
	game.Players = make(Players, p, p)
	game.Zones = make(Zones, z, z)
	game.MeId = id

	for i := 0; i < z; i++ {
		// X: corresponds to the position of the center of a zone. A zone is a circle with a radius of 100 units.
		var zoneX, zoneY int
		fmt.Scan(&zoneX, &zoneY)

		zone := new(Zone)
		zone.Center = &Point{X: zoneX, Y: zoneY}
		game.Zones[i] = zone
	}

	for i := 0; i < p; i++ {
		player := new(Player)
		player.Drones = make(Drones, d, d)

		for ii := 0; ii < d; ii++ {
			drone := new(Drone)
			drone.Location = new(Point)
			player.Drones[ii] = drone
		}

		game.Players[i] = player
	}

	return game
}

func (game *Game) LoadRound() {
	for _, zv := range game.Zones {
		fmt.Scan(&zv.ControllerId)
	}

	for _, pv := range game.Players {
		for _, dv := range pv.Drones {
			fmt.Scan(dv.Location.X, dv.Location.Y)
		}
	}
}
