package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math"
	"math/rand"
)

type Direction string

const (
	Up    = "up"
	Down  = "down"
	Left  = "left"
	Right = "right"
)

var Directions = []Direction{Up, Down, Left, Right}

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "ChrisMcKenzie/snakey-mcksnakeface", // TODO: Your Battlesnake username
		Color:      "#BADA55",                           // TODO: Choose color
		Head:       "default",                           // TODO: Choose head
		Tail:       "default",                           // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
	PrintGameState(state)
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
	PrintGameState(state)
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	coord := moveSafely(state, 1)

	nextMove := CoordToDirection(state.You.Head, coord)

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove, Shout: randomShout()}
}

func CoordToDirection(you Coord, dest Coord) Direction {
	if you.X > dest.X {
		return Left
	} else if you.X < dest.X {
		return Right
	} else if you.Y > dest.Y {
		return Down
	} else if you.Y < dest.Y {
		return Up
	}

	return ""
}

func main() {
	RunServer()
}

func moveSafely(state GameState, depth int) Coord {
	// coord := randomDirection(state.You.Head)

	// pick closer point, between food and smaller snake
	coord := closestFood(state.You.Head, state)
	snake := closestVulnerableSnakeHead(state.You.Head, state)
	// target a vulnerable snake if we are closer than some food
	if snake != nil && distanceToPoint(state.You.Head, *snake) < distanceToPoint(state.You.Head, coord) {
		coord = *snake
	}

	coord = nextCoordinate(state.You.Head, coord)

	if isSafeMove(coord, state.Board, state.You, state.Board.Snakes) {
		return coord
	}

	coord = randomDirection(state.You.Head)
	if isSafeMove(coord, state.Board, state.You, state.Board.Snakes) {
		return coord
	}

	if depth > 50 {
		return coord
	}
	return moveSafely(state, depth+1)
}

func isSafeMove(coord Coord, board Board, mySnake Battlesnake, snakes []Battlesnake) bool {
	// log.Println(coord)
	// log.Println(board)
	// Check if the next move is within the boundaries of the board
	if coord.X < 0 || coord.X == board.Width || coord.Y < 0 || coord.Y == board.Height {
		return false
	}

	// Check if the next move would collide with any snake's body
	for _, snake := range snakes {
		for _, segment := range snake.Body {
			if coord.X == segment.X && coord.Y == segment.Y {
				return false
			}
		}
	}

	// neck := mySnake.Body[1]
	// if neck.X == coord.X { // Neck is left of head, don't move left
	// 	return false
	// } else if neck.Y == coord.Y { // Neck is below head, don't move down
	// 	return false
	// }

	// // Check if the next move would collide with your own snake's body
	for _, segment := range append(mySnake.Body, mySnake.Head) {
		if coord.X == segment.X && coord.Y == segment.Y {
			return false
		}
	}

	return true
}

func randomDirection(snake Coord) Coord {
	dir := rand.Intn(len(Directions))

	var coord Coord
	switch intToDir(dir) {
	case Up:
		coord = Coord{
			X: snake.X,
			Y: snake.Y + 1,
		}
	case Down:
		coord = Coord{
			X: snake.X,
			Y: snake.Y - 1,
		}
	case Left:
		coord = Coord{
			X: snake.X - 1,
			Y: snake.Y,
		}
	case Right:
		coord = Coord{
			X: snake.X + 1,
			Y: snake.Y,
		}
	}

	return coord
}

func closestFood(snake Coord, state GameState) Coord {
	var closest *Coord
	for _, food := range state.Board.Food {
		if closest == nil {
			closest = &food
		} else if distanceToPoint(snake, food) < distanceToPoint(snake, *closest) {
			closest = &food
		}
	}

	return *closest
}

// from where we are now
// to where we want to go
// returns the next coordinate
func nextCoordinate(from, to Coord) Coord {
	switch {
	case from.X < to.X:
		return Coord{X: from.X + 1, Y: from.Y}
	case from.X > to.X:
		return Coord{X: from.X - 1, Y: from.Y}
	case from.Y < to.Y:
		return Coord{X: from.X, Y: from.Y + 1}
	case from.Y > to.Y:
		return Coord{X: from.X, Y: from.Y - 1}
	}

	return to
}

// manhattan distance |X1 – X2| + |Y1 – Y2|
func distanceToPoint(snake Coord, point Coord) int {
	return int(
		math.Abs(float64(snake.X)-float64(point.X)) +
			math.Abs(float64(snake.Y)-float64(point.X)))
}

func randomShout() string {
	shouts := []string{
		"gophers > crabs",
		"🪦🦀",
		"mess with the best die like the rest",
		"ditto will write golang",
	}

	return shouts[rand.Intn(len(shouts))]
}

func intToDir(i int) Direction {
	switch i {
	case 0:
		return Up
	case 1:
		return Down
	case 2:
		return Left
	case 3:
		return Right
	}

	panic("unknown direction")
}

// If there is a bigger snake, find its head, and pick a coordinate away from it so
// we run away from it.
func closestVulnerableSnakeHead(you Coord, state GameState) *Coord {
	var closest *Battlesnake
	for _, other := range state.Board.Snakes {
		if other.Length+1 < state.You.Length {
			if closest == nil {
				closest = &other
			} else if distanceToPoint(you, other.Head) < distanceToPoint(you, closest.Head) {
				closest = &other
			}
		}
	}
	if closest != nil {
		return &closest.Head
	}
	return nil
}

// // If we are
// func runFromBigSnake(me Battlesnake, others []Battlesnake) *Coord {
// 	for _, others := range others {
// 		if me.Length > others.Length {
// 			// its a smaller snake, ignore
// 			continue
// 		}

// 	}

// 	return nil
// }

// type Target int

// const (
// 	TargetFood Target = iota
// 	TargetSnake
// 	AvoidSnake
// )

// func targetPriority(you Battlesnake) Target {
// 	if you.Length < 7 || you.Health < 26 {
// 		return TargetFood
// 	}

// 	return TargetSnake
// }
