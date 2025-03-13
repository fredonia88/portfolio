package main

import (
	"fmt"
)

func main() {
	// initialize err var and game
	// var err error
	// cg := initGame()

	// parse vars from front end and load the game
	// CreateNewGame := true
	/*
	mBoard := [8][8]string{
		{"rk-B-0", "  e-  ", "kg-B-0", "qn-B-0", "  e-  ", "bp-B-1", "kt-B-1", "rk-B-0"},
		{"pn-B-0", "  e-  ", "pn-B-0", "pn-B-0", "pn-B-0", "pn-B-0", "  e-  ", "pn-B-0"},
		{"kt-W-1", "pn-W-0", "  e-  ", "  e-  ", "bp-B-0", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "pn-B-1", "  e-  "},
		{"pn-W-0", "pn-B-1", "pn-W-1", "  e-  ", "bp-W-1", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "pn-W-1", "  e-  ", "bp-W-0", "kt-W-0", "  e-  "},
		{"  e-  ", "  e-  ", "qn-W-1", "  e-  ", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "  e-  ", "  e-  ", "  e-  ", "kg-W-0", "  e-  ", "  e-  ", "rk-W-0"},
	}
		{"rk-B-0", "kt-B-0", "bp-B-0", "qn-B-0", "kg-B-0", "bp-B-0", "kt-B-0", "rk-B-0"},
		{"pn-B-0", "pn-B-0", "pn-B-0", "qn-B-0", "pn-B-0", "pn-B-0", "pn-B-0", "pn-B-0"},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "kt-W-0", "bp-W-0", "qn-W-0", "kg-W-0", "bp-W-0", "kt-W-0", "rk-W-0"},
	*/
	// mPlayer := "W" // User = W or Computer = B
	// mMoveFrom := []int{2, 1}
	// mMoveTo := []int{1, 1}
	// mMoveFromPrior := [2]int{6, 0}
	// mMoveToPrior := [2]int{4, 0}
	// mWCaptured := []string{}
	// mBCaptured := []string{"pn-B-0"}
	
	// if mCreateNewGame {
	// 	err = cg.newGame()
	// 	if err != nil {
	// 		handleError(err)
	// 	}
	// } else {
	// 	err = cg.loadGame(mBoard, mPlayer, mWCaptured, mBCaptured)
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	moveFromPrior, err = cg.getSquare(mMoveFromPrior[0], mMoveFromPrior[1])
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	moveToPrior, err = cg.getSquare(mMoveToPrior[0], mMoveToPrior[1])
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	moveFrom, err := cg.getSquare(mMoveFrom[0], mMoveFrom[1])
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	moveTo, err := cg.getSquare(mMoveTo[0], mMoveTo[1])
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	err = cg.makeMove(moveFrom, moveTo)
	// 	if err != nil {
	// 		handleError(err)
	// 	}

	// 	cm, err := cg.checkMate()
	// 	if err != nil {
	// 		handleError(err)
	// 	}
	// }

	// create a loop to play the game
	var err error
	createNewGame := false
	cg := initGame()
	if createNewGame {
		err = cg.newGame()
		if err != nil {
			handleError(err)
		}
	} else {
		mBoard := [8][8]string{
			{"rk-B-0", "kt-B-0", "bp-B-0", "qn-B-0", "kg-B-0", "bp-B-0", "kt-B-0", "rk-B-0"},
			{"pn-B-0", "  e-  ", "  e-  ", "pn-B-0", "pn-B-0", "pn-B-0", "  e-  ", "  e-  "},
			{"  e-  ", "pn-B-1", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "pn-B-1", "pn-B-1"},
			{"  e-  ", "  e-  ", "pn-B-1", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
			{"  e-  ", "pn-W-1", "  e-  ", "  e-  ", "  e-  ", "pn-W-1", "  e-  ", "  e-  "},
			{"  e-  ", "  e-  ", "pn-W-1", "  e-  ", "  e-  ", "  e-  ", "pn-W-1", "  e-  "},
			{"pn-W-0", "  e-  ", "  e-  ", "pn-W-0", "pn-W-0", "  e-  ", "  e-  ", "pn-W-0"},
			{"rk-W-0", "kt-W-0", "bp-W-0", "qn-W-0", "kg-W-0", "bp-W-0", "kt-W-0", "rk-W-0"},
		}
		mPlayer := "W"
		mWCaptured := []string{}
		mBCaptured := []string{}
		// mMoveFromPrior := [2]int{1, 1}
		// mMoveToPrior := [2]int{2, 1}

		err = cg.loadGame(mBoard, mPlayer, mWCaptured, mBCaptured)
		if err != nil {
			handleError(err)
		}
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

		// update the player and make the computer move
		cg.player = "B"
		cg.makeComputerMove()

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