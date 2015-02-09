package main

import "fmt"
import "math"
import "sort"
import "os"

//fmt.Fprintf(os.Stderr, "Distance: %d, CLOSEST: %d\n",distance, closestDistance)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

const radius = 100.0

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

		for di, drone := range game.Me().Drones {
			closestZones := drone.ClosestZones(game.Zones)
			currentZone := drone.CurrentZone
			var targetZone *Zone

			if currentZone != nil {
				fmt.Fprintf(os.Stderr, "Current Zone Isn't Empty: %d", di)
				currentEnemyShipCount := currentZone.MaxEnemyShips(game.Players)
				currentMyShipCount := currentZone.MyShips

				//WE CONTROL THE CURRENT ZONE
				if currentZone.ControllerId == game.MeId {
					fmt.Fprintf(os.Stderr, "And I control the current zone: %d", di)
					switch {
					case currentEnemyShipCount > 0 && (currentMyShipCount-currentEnemyShipCount) <= 1:
						targetZone = currentZone
						break
					case currentEnemyShipCount == 0 || (currentMyShipCount-currentEnemyShipCount) > 1:
						targetZone = closestZones.ClosestOpenEmpty(drone, game)
						if targetZone == nil {
							targetZone = closestZones.ClosestPush(drone, game)
							if targetZone == nil {
								targetZone = closestZones.ClosestTippingPoint(drone, game)
							}
						}
						break
					}
				} else {
					//ENEMY CONTROLS CURRENT NODE
					fmt.Fprintf(os.Stderr, "And I don't own the current zone: %d", di)
					if targetZone == nil {
						targetZone = closestZones.ClosestPush(drone, game)
						if targetZone == nil {
							targetZone = closestZones.ClosestTippingPoint(drone, game)
						}
					}
				}

			} else {
				fmt.Fprintf(os.Stderr, "Current Zone is empty: %d", di)
				targetZone = closestZones.ClosestOpenEmpty(drone, game)
				if targetZone == nil {
					targetZone = closestZones.ClosestPush(drone, game)
					if targetZone == nil {
						targetZone = closestZones.ClosestTippingPoint(drone, game)
					}
				}
			}

			if targetZone == nil {
				targetZone = closestZones.ClosestTippingPoint(drone, game)
			}

			drone.TargetZone = targetZone

			if drone.TargetZone != nil {
				drone.Move()
			}
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
	Center       *Point
	ControllerId int
	//MaxEnemyShips int
	MyShips int
}

type ZoneDistance struct {
	Zone     *Zone
	Distance float64
}

func (zone *Zone) MaxEnemyShips(players Players) int {
	var maxCount int
	maxCount = 0

	for _, p := range players {
		playerCount := 0
		for _, d := range p.Drones {
			if d.CurrentZone == zone {
				playerCount++
			}
		}
		if playerCount > maxCount {
			maxCount = playerCount
		}
	}

	return maxCount
}

type Zones []*Zone
type ZoneDistances []ZoneDistance

func (a ZoneDistances) Len() int           { return len(a) }
func (a ZoneDistances) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZoneDistances) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func (distances ZoneDistances) ClosestOpenEmpty(drone *Drone, game *Game) *Zone {
	for _, distance := range distances {
		if distance.Zone.ControllerId != game.MeId && distance.Zone.MaxEnemyShips(game.Players) == 0 {
			return distance.Zone
			break
		}
	}

	return nil
}

func (distances ZoneDistances) ClosestPush(drone *Drone, game *Game) *Zone {
	for _, distance := range distances {
		if distance.Zone.ControllerId != game.MeId && distance.Zone.MaxEnemyShips(game.Players)-distance.Zone.MyShips == 1 {
			return distance.Zone
			break
		}
	}
	return nil
}

func (distances ZoneDistances) ClosestTippingPoint(drone *Drone, game *Game) *Zone {
	var returnDifference int
	returnDifference = 10000
	var targetZone *Zone

	for _, distance := range distances {
		if distance.Zone.ControllerId != game.MeId {
			difference := distance.Zone.MaxEnemyShips(game.Players) - distance.Zone.MyShips

			if difference < returnDifference {
				returnDifference = difference
				targetZone = distance.Zone
			}
		}
	}

	return targetZone
}

//DRONE TYPES
type Drone struct {
	Location    *Point
	TargetZone  *Zone
	CurrentZone *Zone
}

func (drone *Drone) Move() {
	drone.MoveTo(drone.TargetZone.Center)
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

	return returnDistances
}

func (drone *Drone) SetCurrentZone(zones Zones) bool {
	for _, z := range zones {
		distance := drone.Location.DistanceTo(z.Center)

		if distance <= radius {
			drone.CurrentZone = z
			//fmt.Fprintf(os.Stderr, "SETTING CURRENT ZONE:",drone.CurrentZone)
			return true
		}
	}

	drone.CurrentZone = nil
	return false
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
		zv.MyShips = 0
		//zv.MaxEnemyShips = 0
	}

	for pi, pv := range game.Players {
		for _, dv := range pv.Drones {
			fmt.Scan(&dv.Location.X, &dv.Location.Y)
			if dv.SetCurrentZone(game.Zones) {
				if dv.CurrentZone.ControllerId == pi {
					if game.MeId == pi {
						dv.CurrentZone.MyShips++
					}
				}
			}
		}
	}
}
