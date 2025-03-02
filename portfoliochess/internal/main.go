package main

import (
	"fmt"
)

func main() {
	// initialize err var and game
	var err error
	cg := initGame()

	// parse vars from front end and load the game
	mCreateNewGame := false
	mBoard := [8][8]string{
		{"rk-B-0", "  e-  ", "  e-  ", "  e-  ", "kg-B-0", "bp-B-1", "kt-B-1", "rk-B-0"},
		{"pn-B-0", "  e-  ", "pn-B-0", "qn-B-0", "pn-B-0", "pn-B-0", "  e-  ", "pn-B-0"},
		{"kt-B-0", "pn-W-0", "  e-  ", "pn-B-0", "bp-B-0", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "pn-B-1", "  e-  "},
		{"pn-W-0", "pn-B-1", "pn-W-1", "pn-W-1", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "bp-W-1", "bp-W-0", "kt-W-0", "  e-  "},
		{"  e-  ", "  e-  ", "qn-W-1", "  e-  ", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "kt-W-0", "  e-  ", "  e-  ", "kg-W-0", "  e-  ", "  e-  ", "rk-W-0"},
	}
	/*
		{"rk-B-0", "kt-B-0", "bp-B-0", "qn-B-0", "kg-B-0", "bp-B-0", "kt-B-0", "rk-B-0"},
		{"pn-B-0", "pn-B-0", "pn-B-0", "qn-B-0", "pn-B-0", "pn-B-0", "pn-B-0", "pn-B-0"},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "kt-W-0", "bp-W-0", "qn-W-0", "kg-W-0", "bp-W-0", "kt-W-0", "rk-W-0"},
	*/
	mPlayer := "B" // User = W or Computer = B
	mMoveFrom := []int{0, 4}
	mMoveTo := []int{0, 0}
	mMoveFromPrior := [2]int{6, 0}
	mMoveToPrior := [2]int{4, 0}
	mWCaptured := []string{}
	mBCaptured := []string{"pn-B-0"}
	
	if mCreateNewGame {
		err = cg.newGame()
		if err != nil {
			handleError(err)
		}
		fmt.Println()
		fmt.Println("New game:")
		cg.displayBoard()
	} else {
		err = cg.loadGame(mBoard, mPlayer, mWCaptured, mBCaptured)
		if err != nil {
			handleError(err)
		}
		fmt.Println()
		fmt.Println("Game loaded:")
		cg.displayBoard()

		moveFromPrior, err = cg.getSquare(mMoveFromPrior[0], mMoveFromPrior[1])
		if err != nil {
			handleError(err)
		}

		moveToPrior, err = cg.getSquare(mMoveToPrior[0], mMoveToPrior[1])
		if err != nil {
			handleError(err)
		}

		moveFrom, err := cg.getSquare(mMoveFrom[0], mMoveFrom[1])
		if err != nil {
			handleError(err)
		}

		moveTo, err := cg.getSquare(mMoveTo[0], mMoveTo[1])
		if err != nil {
			handleError(err)
		}

		fmt.Println()
		msg := fmt.Sprintf("Attempting move from (%d, %d) to (%d, %d)", moveFrom.row, moveFrom.col, moveTo.row, moveTo.col)
		fmt.Println(msg)
		err = cg.makeMove(moveFrom, moveTo)
		if err != nil {
			handleError(err)
		}
		fmt.Println()
		fmt.Println("Board after move:")
		cg.displayBoard()
		fmt.Println()
		cg.displayCapturedPieces()
		cg.displayBoardPositions()
	}
}