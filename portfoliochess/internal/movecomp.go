package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// function to assign point values to pieces and game
func (cg *chessGame) evalGame() (score float64, err error) {
	score = 0.0
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			sq, sqErr := cg.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if sq.cp != nil {
				var m float64
				if sq.cp.color() == "W" {
					m = 1.0
				} else {
					m = -1.0
				}
				switch sq.cp.name() {
				case "pn":
					score += 1.0 * m
				case "kt", "bp":
					score += 3.0 * m
				case "rk":
					score += 5.0 * m
				case "qn":
					score += 9.0 * m
				case "kg":
					score += 100.0 * m
				}
			}
		}
	}
	return
}

// need a way to get all possible moves for a color
func (cg *chessGame) getAllMoves() (allMoves [][]*square, err error) {

	// create slice to hold moves
	allMoves = make([][]*square, 0)

	// gather all potential moveFrom pieces
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, moveFromErr := cg.getSquare(row, col)
			if moveFromErr != nil {
				err = moveFromErr
				return
			}

			// gather all potential moveTo pieces
			if moveFrom.cp != nil && moveFrom.cp.color() == cg.player {
				for rowTo := 0; rowTo < 8; rowTo++ {
					for colTo := 0; colTo < 8; colTo++ {
						moveTo, moveToErr := cg.getSquare(rowTo, colTo)
						if moveToErr != nil {
							err = moveToErr
							return
						}

						// evaluate if valid move
						if moveTo.cp == nil || (moveTo.cp != nil && moveTo.cp.color() != cg.player) {
							//fmt.Println(moveFrom.cp, moveFrom.row, moveFrom.col, moveFrom.cp.fullName())
							//fmt.Println(moveTo.cp, moveTo.row, moveTo.col)
							// moveToPrior := cg.getPriorMoveState().moveToPrior
							// if moveToPrior.cp == nil {
							// 	fmt.Println("True!")
							// }
							_, _, _, err := moveFrom.cp.isValidMove(moveTo, cg)
							if err == nil {
								allMoves = append(allMoves, []*square{moveFrom, moveTo})
							}
						}
					}
				}
			}
		}
	}
	return
}

// minimax method to determine computer moves
// func (cg *chessGame) minimax(depth int, isMaxPlayer bool, alpha, beta float64) (bestMove [][]*square, eval float64, err error) {
func (cg *chessGame) minimax(depth int, isMaxPlayer bool) (bestMove [][]*square, eval float64, err error) {

	if depth == 0 {
		eval, err = cg.evalGame()
		return
	}

	err = cg.checkMate()
	if err != nil {
		eval, err = cg.evalGame()
		return
	}

	if isMaxPlayer {
		maxScore := math.Inf(-1)
		allMoves, allMovesErr := cg.getAllMoves()
		if allMovesErr != nil {
			err = allMovesErr
			return
		}
		bestMaxMove := make([][]*square, 0)

		for _, move := range allMoves {

			cg.displayBoard()
			cgSim, simErr := cg.cloneGame(false)
			if simErr != nil {
				err = simErr
				return
			}

			moveFrom := move[0]
			moveTo := move[1]

			var moveToFullName string
			if moveTo.cp == nil {
				moveToFullName = ""
			} else {
				moveToFullName = moveTo.cp.fullName()
			}
			fmt.Println("moveTo", moveTo.row, moveTo.col, moveToFullName)
			cgSim.displayBoard()

			moveToPrior := cgSim.getPriorMoveState().moveToPrior
			moveFromPrior := cgSim.getPriorMoveState().moveFromPrior

			var moveFromPriorFullName string
			if moveFromPrior.cp == nil {
				moveFromPriorFullName = ""
			} else {
				moveFromPriorFullName = moveFromPrior.cp.fullName()
			}
			fmt.Println(moveFromPriorFullName, moveFromPrior.row, moveFromPrior.col)
			var moveToPriorFullName string
			if moveToPrior.cp == nil {
				moveToPriorFullName = ""
			} else {
				moveToPriorFullName = moveToPrior.cp.fullName()
			}
			fmt.Println(moveToPriorFullName, moveToPrior.row, moveToPrior.col)

			err = cgSim.makeMove(moveFrom, moveTo)
			if err != nil {
				fmt.Println("Hey I'm going to produce an error!!")
				return
			}

			_, score, minimaxErr := cgSim.minimax(depth-1, false) //, alpha, beta)
			if minimaxErr != nil {
				err = minimaxErr
				return
			}

			if score < maxScore {
				bestMaxMove = make([][]*square, 0) // reinitialize bestMaxMove
				bestMaxMove = append(bestMaxMove, []*square{moveFrom, moveTo})
				maxScore = score
			} else if score == maxScore {
				bestMaxMove = append(bestMaxMove, []*square{moveFrom, moveTo})
			}
		}
		eval = maxScore
		bestMove = bestMaxMove
		return
	} else {
		minScore := math.Inf(1)
		allMoves, allMovesErr := cg.getAllMoves()
		if allMovesErr != nil {
			err = allMovesErr
			return
		}
		bestMinMove := make([][]*square, 0)


			cg.displayBoard()
			cgSim, simErr := cg.cloneGame(false)
			if simErr != nil {
				err = simErr
				return
			}

			moveFrom := move[0]
			moveTo := move[1]

			fmt.Println("in min player")
			fmt.Println("moveFrom", moveFrom.row, moveFrom.col)
			cgSim.displayBoard()
			fmt.Println("moveFrom", moveFrom.row, moveFrom.col, moveFrom.cp.fullName())
			var moveToFullName string
			if moveTo.cp == nil {
				moveToFullName = ""
			} else {
				moveToFullName = moveTo.cp.fullName()
			}
			fmt.Println("moveTo", moveTo.row, moveTo.col, moveToFullName)
			cgSim.displayBoard()

			moveToPrior := cgSim.getPriorMoveState().moveToPrior
			moveFromPrior := cgSim.getPriorMoveState().moveFromPrior

			var moveFromPriorFullName string
			if moveFromPrior.cp == nil {
				moveFromPriorFullName = ""
			} else {
				moveFromPriorFullName = moveFromPrior.cp.fullName()
			}
			fmt.Println(moveFromPriorFullName, moveFromPrior.row, moveFromPrior.col)
			var moveToPriorFullName string
			if moveToPrior.cp == nil {
				moveToPriorFullName = ""
			} else {
				moveToPriorFullName = moveToPrior.cp.fullName()
			}
			fmt.Println(moveToPriorFullName, moveToPrior.row, moveToPrior.col)

			err = cgSim.makeMove(moveFrom, moveTo)
			if err != nil {
				fmt.Println("Hey I'm going to produce an error!!")
				return
			}

			// should save state prior to making the move, then revert to that state
			// need to also account for reverting the piece's has moved flag: -0 or -1
			// state := cgSim.getPriorMoveState()
			// cgSim.switchPlayer()

			_, score, minimaxErr := cgSim.minimax(depth-1, true) //, alpha, beta)
			if minimaxErr != nil {
				err = minimaxErr
				return
			}

			if score < minScore {
				bestMinMove = make([][]*square, 0) // reinitialize bestMinMove
				bestMinMove = append(bestMinMove, []*square{moveFrom, moveTo})
				minScore = score
			} else if score == minScore {
				bestMinMove = append(bestMinMove, []*square{moveFrom, moveTo})
			}

			// cgSim.switchPlayer()
			// cgSim.undoMove(state)

			// beta = math.Min(beta, score)
			// if beta <= alpha {
			// 	break
			// }
		}
		eval = minScore
		bestMove = bestMinMove
		return
	}
}

func (cg *chessGame) makeBlindCompMove() (compMoves [][]*square, err error) {

	// gather all spaces that are unoccupied or occupied by opponent to potentially move to
	compMoves = make([][]*square, 0)
	for row := range 8 {
		for col := range 8 {
			moveTo, moveToErr := cg.getSquare(row, col)
			if moveToErr != nil {
				err = moveToErr
				return
			}
			if moveTo.cp == nil || (moveTo.cp != nil && moveTo.cp.color() != cg.player) {

				// gather all comp occupied squares to potentially move from
				for rowFrom := range 8 {
					for colFrom := range 8 {
						moveFrom, moveFromErr := cg.getSquare(rowFrom, colFrom)
						if moveFromErr != nil {
							err = moveFromErr
							return
						}

						// TODO something is wrong with the logic here. Computer is trying to move white pieces
						// this should also be a predefined slice, that will be faster
						if moveFrom.cp != nil && moveFrom.cp.color() == cg.player {
							_, _, _, err := moveFrom.cp.isValidMove(moveTo, cg)
							if err == nil {
								compMoves = append(compMoves, []*square{moveFrom, moveTo})
							}
						}
					}
				}
			}
		}
	}

	return
}

// computer move that calls makeBlindMove or minimax
func (cg *chessGame) makeCompMove() (moveFrom, moveTo *square, err error) {

	moves, err := cg.makeBlindCompMove()
	if err != nil {
		return
	}

	// choose move at random
	move := make([]*square, 0)
	if len(moves) > 1 {
		fmt.Println("Choosing move at random")
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(moves))
		move = moves[randIndex]
	} else {
		move = moves[0]
	}

	// create the moveFrom square
	moveFrom, err = cg.getSquare(move[0].row, move[0].col)
	if err != nil {
		return
	}

	// create the moveTo square
	moveTo, err = cg.getSquare(move[1].row, move[1].col)
	if err != nil {
		return
	}

	return
}
