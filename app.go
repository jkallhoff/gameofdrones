package main

import "fmt"
import "math"
import "sort"

//import "os"
//fmt.Fprintf(os.Stderr, "Distance: %d, CLOSEST: %d\n",distance, closestDistance)

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

	for {
		game.LoadRound()

		for _, drone := range game.Me().Drones {
			closestZones := drone.ClosestZones(game.Zones)

			if drone.TargetZone == nil {
				drone.TargetZone = closestZones[0].Zone
				drone.Move()
			}

			if drone.TargetZone != nil {
				if drone.AtZone != true {
					drone.Move()
				} else {
					drone.Move()

				}
			}
			//closestZones := d.ClosestZones(game.Zones)
			//d.MoveTo(closestZones[0].Zone.Center)
		}
	}
}

//POINT TYPES
type Point struct {
	X, Y int
}

func (start *Point) DistanceTo(end *Point) float64 {
	xpower := (end.X - start.X) * (end.X - start.X)
	ypower := (end.Y - start.Y) * (end.Y - start.Y)
	preSquare := xpower + ypower

	return math.Sqrt(float64(preSquare))
}

//ZONE TYPES
type Zone struct {
	Center        *Point
	ControllerId  int
	MaxEnemyShips int
}

type ZoneDistance struct {
	Zone     *Zone
	Distance float64
}

type Zones []*Zone
type ZoneDistances []ZoneDistance

func (a ZoneDistances) Len() int           { return len(a) }
func (a ZoneDistances) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZoneDistances) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

//DRONE TYPES
type Drone struct {
	Location   *Point
	TargetZone *Zone
	AtZone     bool
}

func (drone *Drone) Move() {
	fmt.Printf("%d %d\n", drone.TargetZone.Center.X, drone.TargetZone.Center.Y)
}

func (drone *Drone) MoveTo(point *Point) {
	fmt.Printf("%d %d\n", point.X, point.Y)
}

func (drone *Drone) ClosestZones(zones Zones) ZoneDistances {
	returnDistances := make(ZoneDistances, len(zones), len(zones))

	for i, z := range zones {
		distance := drone.Location.DistanceTo(z.Center)
		zoneDistance := ZoneDistance{Zone: z, Distance: distance}
		returnDistances[i] = zoneDistance
	}

	sort.Sort(returnDistances)

	if drone.TargetZone != nil {
		if returnDistances[0].Zone == drone.TargetZone && returnDistances[0].Distance <= 99 {
			drone.AtZone = true
		} else {
			drone.AtZone = false
		}
	}

	return returnDistances
}

type Drones []*Drone

//PLAYER TYPES
type Player struct {
	Drones Drones
}

func (player *Player) SendNextDroneTo(point *Point) {
	fmt.Printf("%d %d\n", point.X, point.Y)
}

type Players []*Player

//GAME TYPES
type Game struct {
	Players Players //The collection of players in the game.
	Zones   Zones
	MeId    int
	Init    bool
}

func (game *Game) Me() *Player {
	return game.Players[game.MeId]
}

//SYSTEM FUNCTIONS
func SetupGame(p, id, d, z int) *Game {
	game := new(Game)
	game.Players = make(Players, p, p)
	game.Zones = make(Zones, z, z)
	game.MeId = id
	game.Init = true

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
			fmt.Scan(&dv.Location.X, &dv.Location.Y)
		}
	}
}
