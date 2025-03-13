package main 

import (
	"fmt"
	"strings"
)


var moveTo *square
var moveFrom *square
var moveFromPrior *square
var moveToPrior *square

type chessGame struct {
	board [8][8]square
	player string
	wCaptured []chessPiece
	bCaptured []chessPiece
}

type square struct {
	row int // rank
	col int // file
	cp chessPiece
}

func initGame() *chessGame {
	var board [8][8]square
	var player string
	var wCaptured []chessPiece
	var bCaptured []chessPiece

	cg := &chessGame{board, player, wCaptured, bCaptured}

	return cg
}

func (cg *chessGame) newGame() (err error) {
	cg.player = "W"
	majorMinor := []string{"rk","kt","bp","qn","kg","bp","kt","rk"}
	for i := 0; i < 8; i++ {

		majorMinorConstructor, majorMinorErr := getPieceConstructor(majorMinor[i])
		if majorMinorErr != nil {
			err = majorMinorErr
			return
		}
		cg.board[0][i] = square{0, i, majorMinorConstructor(majorMinor[i] + "-B-0", 0, i)}
		cg.board[7][i] = square{7, i, majorMinorConstructor(majorMinor[i] + "-W-0", 7, i)}

		pawnConstructor, pawnErr := getPieceConstructor("pn")
		if pawnErr != nil {
			err = pawnErr
			return
		}
		cg.board[1][i] = square{1, i, pawnConstructor("pn-B-0", 1, i)}
		cg.board[6][i] = square{6, i, pawnConstructor("pn-W-0", 6, i)}

		cg.board[2][i] = square{2, i, nil}
		cg.board[3][i] = square{3, i, nil}
		cg.board[4][i] = square{4, i, nil}
		cg.board[5][i] = square{5, i, nil}
	}

	return
}

func (cg *chessGame) loadGame(mGame [8][8]string, mPlayer string, mWCaptured []string, mBCaptured []string) (err error) {
	cg.player = mPlayer
	for row := range mGame {
		for col := range mGame[row] {
			piece := strings.TrimSpace(strings.Split(mGame[row][col], "-")[0])

			if piece == "e" {
				cg.board[row][col] = square{row, col, nil}
				continue
			}
			
			constructor, conErr := getPieceConstructor(piece)
			if conErr != nil {
				err = conErr
				return
			}
			
			cg.board[row][col] = square{row, col, constructor(mGame[row][col], row, col)}
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
	
	return 
}

func (cg *chessGame) unloadGame() (board [8][8]string, player string, wCaptured []string, bCaptured []string) {
	for row := range cg.board {
		for col := range cg.board[row] {
			cp := cg.board[row][col].cp
			if cp == nil {
				board[row][col] = " e- "
			} else {
				board[row][col] = cp.fullName()
			}
		}
	}

	player = cg.player
	
	for _, cp := range cg.wCaptured {
		wCaptured = append(wCaptured, cp.fullName())
	}

	for _, cp := range cg.bCaptured {
		bCaptured = append(bCaptured, cp.fullName())
	}

	return
}

func (cg *chessGame) cloneGame(changePlayer bool) (cgSim *chessGame, err error) {
	cgSimBoard, cgSimPlayer, cgSimWCaptured, cgSimBCaptured := cg.unloadGame()
	cgSim = initGame()
	loadErr := cgSim.loadGame(cgSimBoard, cgSimPlayer, cgSimWCaptured, cgSimBCaptured)
	if loadErr != nil {
		err = loadErr
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
	for row := king.row - 1; row <= king.row + 1; row++ {
		for col := king.col - 1; col <= king.col + 1; col++ {
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
					cgSim.capturePiece(moveTo)
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
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, sqErr := cgSim.getSquare(row, col)
			if sqErr != nil {
				err = sqErr
				return
			}
			if moveFrom.cp != nil && moveFrom.cp.name() != "kg" && moveFrom.cp.color() == cgSim.player {
				for rowTo := 0; rowTo < 8; rowTo++ {
					for colTo := 0; colTo < 8; colTo++ {
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
									cgSim.capturePiece(moveTo)
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
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
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
	return &cg.board[row][col], nil
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

// function to check if piece positions align with board
func (cg *chessGame) displayBoardPositions() {
	for row := 0; row < 8; row++ {
		line := make([][2]int, 8)
		for col := 0; col < 8; col++ {
			line[col] = [2]int{cg.board[row][col].row, cg.board[row][col].col}
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