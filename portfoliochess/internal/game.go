package main

import (
	"fmt"
	"strings"
)

type square struct {
	row int // rank
	col int // file
	cp  chessPiece
}

type priorMoveState struct {
	moveFromPrior square
	moveToPrior   square
	promotedPawn  bool
	capturedPiece chessPiece
}

type chessGame struct {
	board     [8][8]*square
	player    string
	wCaptured []chessPiece
	bCaptured []chessPiece
	priorMove priorMoveState
}

func initGame() *chessGame {
	var board [8][8]*square
	var player string
	var wCaptured []chessPiece
	var bCaptured []chessPiece
	var priorMove priorMoveState

	cg := &chessGame{board, player, wCaptured, bCaptured, priorMove}

	return cg
}

func (cg *chessGame) newGame() (err error) {
	cg.player = "W"
	majorMinor := []string{"rk", "kt", "bp", "qn", "kg", "bp", "kt", "rk"}
	for i := 0; i < 8; i++ {

		majorMinorConstructor, majorMinorErr := getPieceConstructor(majorMinor[i])
		if majorMinorErr != nil {
			err = majorMinorErr
			return
		}
		cg.board[0][i] = &square{0, i, majorMinorConstructor(majorMinor[i]+"-B-0", 0, i)}
		cg.board[7][i] = &square{7, i, majorMinorConstructor(majorMinor[i]+"-W-0", 7, i)}

		pawnConstructor, pawnErr := getPieceConstructor("pn")
		if pawnErr != nil {
			err = pawnErr
			return
		}
		cg.board[1][i] = &square{1, i, pawnConstructor("pn-B-0", 1, i)}
		cg.board[6][i] = &square{6, i, pawnConstructor("pn-W-0", 6, i)}

		cg.board[2][i] = &square{2, i, nil}
		cg.board[3][i] = &square{3, i, nil}
		cg.board[4][i] = &square{4, i, nil}
		cg.board[5][i] = &square{5, i, nil}
	}

	return
}

func (cg *chessGame) loadGame(
	mGame [8][8]string,
	mPlayer string,
	mWCaptured []string,
	mBCaptured []string,
	mMoveFromPrior [2]int,
	mMoveToPrior [2]int,
	mPromotedPawn bool,
	mCapturedPiece string,
) (err error) {

	cg.player = mPlayer
	for row := range mGame {
		for col := range mGame[row] {
			piece := strings.TrimSpace(strings.Split(mGame[row][col], "-")[0])

			if piece == "e" {
				cg.board[row][col] = &square{row, col, nil}
				continue
			}

			constructor, conErr := getPieceConstructor(piece)
			if conErr != nil {
				err = conErr
				return
			}

			cg.board[row][col] = &square{row, col, constructor(mGame[row][col], row, col)}
		}
	}

	// load captured pieces
	for _, p := range mWCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, conErr := getPieceConstructor(piece)
		if conErr != nil {
			err = conErr
			return
		}
		cg.wCaptured = append(cg.wCaptured, constructor(p, -1, -1))
	}

	for _, p := range mBCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, conErr := getPieceConstructor(piece)
		if conErr != nil {
			err = conErr
			return
		}
		cg.bCaptured = append(cg.bCaptured, constructor(p, -1, -1))
	}

	// set moveFromPrior and moveToPrior
	moveFromPrior, err := cg.getSquare(mMoveFromPrior[0], mMoveFromPrior[1])
	if err != nil {
		return
	}

	moveToPrior, err := cg.getSquare(mMoveToPrior[0], mMoveToPrior[1])
	if err != nil {
		return
	}

	// create last captured piece
	var capturedPiece chessPiece
	if mCapturedPiece != "" {
		constructor := pieces[strings.Split(mCapturedPiece, "-")[0]]
		capturedPiece = constructor(mCapturedPiece, -1, -1)
	}

	// update prior move state
	cg.setPriorMoveState(*moveFromPrior, *moveToPrior, mPromotedPawn, capturedPiece)

	return
}

func (cg *chessGame) unloadGame() (
	mBoard [8][8]string,
	mPlayer string,
	mWCaptured []string,
	mBCaptured []string,
	mMoveFromPrior [2]int,
	mMoveToPrior [2]int,
	mPromotedPawn bool,
	mCapturedPiece string,
) {

	// add pieces and empty space strings to board
	for row := range cg.board {
		for col := range cg.board[row] {
			cp := cg.board[row][col].cp
			if cp == nil {
				mBoard[row][col] = " e- "
			} else {
				mBoard[row][col] = cp.fullName()
			}
		}
	}

	// set player
	mPlayer = cg.player

	// add captured pieces' names to respective slice
	for _, cp := range cg.wCaptured {
		mWCaptured = append(mWCaptured, cp.fullName())
	}

	for _, cp := range cg.bCaptured {
		mBCaptured = append(mBCaptured, cp.fullName())
	}

	// set prior moves
	state := cg.getPriorMoveState()
	mMoveFromPrior = [2]int{state.moveFromPrior.row, state.moveFromPrior.col}
	mMoveToPrior = [2]int{state.moveToPrior.row, state.moveToPrior.col}

	// set promotedPawn and capturedPiece
	mPromotedPawn = state.promotedPawn
	if state.capturedPiece == nil {
		mCapturedPiece = ""
	} else {
		mCapturedPiece = state.capturedPiece.fullName()
	}

	return
}

func (cg *chessGame) cloneGame(changePlayer bool) (cgSim *chessGame, err error) {

	cgSimBoard, cgSimPlayer, cgSimWCaptured, cgSimBCaptured, cgSimMoveFromPrior, cgSimMoveToPrior, cgSimPromotedPawn, cgSimCapturedPiece := cg.unloadGame()
	cgSim = initGame()
	err = cgSim.loadGame(cgSimBoard, cgSimPlayer, cgSimWCaptured, cgSimBCaptured, cgSimMoveFromPrior, cgSimMoveToPrior, cgSimPromotedPawn, cgSimCapturedPiece)
	if err != nil {
		return
	}
	if changePlayer {
		if cgSim.player == "W" {
			cgSim.player = "B"
		} else {
			cgSim.player = "W"
		}
	}

	return
}

func (cg *chessGame) setPriorMoveState(moveFromPrior, moveToPrior square, promotedPawn bool, capturedPiece chessPiece) {
	cg.priorMove = priorMoveState{
		moveFromPrior: moveFromPrior,
		moveToPrior:   moveToPrior,
		promotedPawn:  promotedPawn,
		capturedPiece: capturedPiece,
	}
}

func (cg *chessGame) getPriorMoveState() priorMoveState {

	return cg.priorMove
}

var visitedSquare = make(map[*square]bool)

func (cg *chessGame) inCheck() (err error) {

	// get the king square
	king, kingErr := cg.getKing()
	if kingErr != nil {
		err = kingErr
		return
	}

	// denote if this square has been evaluated
	if visitedSquare[king] {
		return
	}
	visitedSquare[king] = true

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, sqErr := cg.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if moveFrom.cp != nil && moveFrom.cp.color() != cg.player {
				_, _, _, moveErr := moveFrom.cp.isValidMove(king, cg)
				if moveErr == nil {
					err = newChessError(errInCheck, "King will be in check at (%d %d) by %s at (%d %d)",
						king.row, king.col, moveFrom.cp.fullName(), moveFrom.cp.getRow(), moveFrom.cp.getCol())
					return
				}
			}
		}
	}
	return
}

func (cg *chessGame) checkMate() (err error) {

	// clone the game
	cgSim, simErr := cg.cloneGame(true)
	if simErr != nil {
		err = simErr
		return
	}

	// is the king in check? if not, no need to check for check mate
	err = cgSim.inCheck()
	if err == nil {
		return
	}

	// set new king square for simulated game
	king, kingErr := cgSim.getKing()
	if kingErr != nil {
		err = kingErr
		return
	}

	// check if king can move to a square that is not in check
	for row := king.row - 1; row <= king.row+1; row++ {
		for col := king.col - 1; col <= king.col+1; col++ {
			if row < 0 || row > 7 || col < 0 || col > 7 {
				continue
			}
			moveTo, sqErr := cgSim.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if moveTo.cp == nil || (moveTo.cp != nil && moveTo.cp.color() != cgSim.player) {

				// simulate moving the king to that position
				king.cp.updatePosition(moveTo.row, moveTo.col)
				var pieceCaptured bool
				if moveTo.cp != nil {
					cgSim.capturePiece(*moveTo)
					pieceCaptured = true
				}
				moveTo.cp = king.cp
				king.cp = nil

				// exit if king is not in check
				err = cgSim.inCheck()
				if err == nil {

					// undo the move
					king.cp = moveTo.cp
					king.cp.updatePosition(king.row, king.col)
					moveTo.cp = nil
					if pieceCaptured {
						cgSim.recoverLastCapturedPiece(moveTo)
					}
					return
				}

				// undo the move
				king.cp = moveTo.cp
				king.cp.updatePosition(king.row, king.col)
				moveTo.cp = nil
				if pieceCaptured {
					cgSim.recoverLastCapturedPiece(moveTo)

				}
			}
		}
	}

	// loop through each piece and check if it can move to a square that will block the check
	for row := range 8 {
		for col := range 8 {
			moveFrom, sqErr := cgSim.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if moveFrom.cp != nil && moveFrom.cp.name() != "kg" && moveFrom.cp.color() == cgSim.player {
				for rowTo := range 8 {
					for colTo := range 8 {
						moveTo, sqErr := cgSim.getSquare(rowTo, colTo)
						if sqErr != nil {
							err = sqErr
							return
						}
						if moveTo.cp == nil || (moveTo.cp != nil && moveTo.cp.color() != cgSim.player) {
							_, _, _, err = moveFrom.cp.isValidMove(moveTo, cgSim)
							if err == nil {

								// simulate moving the piece
								moveFrom.cp.updatePosition(rowTo, colTo)
								var pieceCaptured bool
								if moveTo.cp != nil {
									cgSim.capturePiece(*moveTo)
									pieceCaptured = true
								}
								moveTo.cp = moveFrom.cp
								moveFrom.cp = nil

								// check if king is in check, and clear the visited variable
								visitedSquare = make(map[*square]bool)
								err = cgSim.inCheck()
								if err == nil {
									return
								}

								// undo the move
								moveFrom.cp = moveTo.cp
								moveFrom.cp.updatePosition(row, col)
								moveTo.cp = nil
								if pieceCaptured {
									cgSim.recoverLastCapturedPiece(moveTo)
								}
							}
						}
					}
				}
			}
		}
	}
	err = newChessError(errCheckMate, "Checkmate! %s wins!", cg.player)
	return
}

func (cg *chessGame) getKing() (king *square, err error) {
	for row := range 8 {
		for col := range 8 {
			sq, sqErr := cg.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if sq.cp != nil && sq.cp.name() == "kg" && sq.cp.color() == cg.player {
				king = sq
				return
			}
		}
	}
	err = newChessError(errKingNotFound, "King not found")
	return
}

func (cg *chessGame) getSquare(row, col int) (*square, error) {
	if row > 7 || row < 0 || col > 7 || col < 0 {
		return &square{}, newChessError(errOutOfRange, "Selected square (%d %d) is out of range", row, col)
	}
	return cg.board[row][col], nil
}

func (cg *chessGame) displayBoard() {
	for row := 0; row < 8; row++ {
		line := make([]string, 8)
		for col := 0; col < 8; col++ {
			if cg.board[row][col].cp == nil {
				line[col] = "|  e-  |"
			} else {
				line[col] = "|" + cg.board[row][col].cp.fullName() + "|"
			}
		}
		fmt.Println(line)
	}
}

func (cg *chessGame) displayCapturedPieces() {
	var w []string
	var b []string
	for _, piece := range cg.wCaptured {
		w = append(w, piece.fullName())
	}
	for _, piece := range cg.bCaptured {
		b = append(b, piece.fullName())
	}
	fmt.Println("White pieces captured:", w)
	fmt.Println("Black pieces captured:", b)
}
