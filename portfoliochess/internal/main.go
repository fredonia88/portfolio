package main

import (
	"fmt"
)

func main() {

	// create a loop to play the game
	var err error
	cg := initGame()
	err = cg.newGame()
	if err != nil {
		handleError(err)
	}

	for {
		cg.displayBoard()
		fmt.Println("")
		cg.displayCapturedPieces()

		// get the move from and move to for the player
		fmt.Println("")
		fmt.Println("It is this player's turn: ", cg.player)
		fmt.Println("Enter the row and column of the piece you want to move")
		var tarRow, tarCol int
		fmt.Scanln(&tarRow, &tarCol)

		fmt.Println("Enter the row and column of the square you want to move to:")
		var desRow, desCol int
		fmt.Scanln(&desRow, &desCol)

		// make the move
		moveFrom, err := cg.getSquare(tarRow, tarCol)
		if err != nil {
			handleError(err)
		}

		if moveFrom.cp == nil {
			newChessError(errEmptySquare, "Square (%d %d) is empty. Choose a square with your piece", moveFrom.row, moveFrom.col)
		}

		moveTo, err := cg.getSquare(desRow, desCol)
		if err != nil {
			handleError(err)
		}

		moveToFullName := "  e-  "
		if moveTo.cp != nil {
			moveToFullName = moveTo.cp.fullName()
		}

		fmt.Println("")
		msg := fmt.Sprintf("Attempting to move %s at (%d, %d) to %s at (%d %d)...",
			moveFrom.cp.fullName(), moveFrom.row, moveFrom.col, moveToFullName, moveTo.row, moveTo.col)
		fmt.Println(msg)
		err = cg.makeMove(moveFrom, moveTo)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Move successful")

		// check for checkmate
		err = cg.checkMate()
		if err != nil {
			handleError(err)
		}

		// update the player, get the best computer move and make the computer move
		cg.player = "B"
		moveFromComp, moveToComp, err := cg.makeCompMove()
		if err != nil {
			handleError(err)
		}
		err = cg.makeMove(moveFromComp, moveToComp)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Computer move successful")

		// check for checkmate
		err = cg.checkMate()
		if err != nil {
			handleError(err)
		}

		score, err := cg.evalGame()
		if err != nil {
			handleError(err)
		}
		fmt.Println("Game score: ", score)

		cg.player = "W"
	}
}
