package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	TODO:
		1. Finish TODOs
		2. Start writing game state to the database
		1. Play an actual game
		3. Build the mini max algo for computer moves
*/ 

var verticalDistance int
var horizontalDistance int
var isValidDiagonal bool 
var isVerticalOnly bool
var isHorizontalOnly bool

func setChessMove(moveFromRow, moveFromCol int, moveTo *square) {
	verticalDistance = moveTo.row - moveFromRow
	horizontalDistance = moveTo.col - moveFromCol
	isValidDiagonal = abs(verticalDistance) == abs(horizontalDistance)
	isVerticalOnly = moveTo.col == moveFromCol && moveTo.row != moveFromRow
	isHorizontalOnly = moveTo.col != moveFromCol && moveTo.row == moveFromRow
}

func abs(x int) int {
	if x < 0 {
		x *= -1
	}
	return x
}

func (cg *chessGame) evalCastle(moveFromRow, moveFromCol, moveToCol int, moveFromFullName string) (err error) {

	// collect all squares the king would start at, move through or end on while castling
	squaresToEval := make([]*square, 0, 4)

	endCol := 0
	colStep := -1
	if moveFromCol == 7 || moveToCol == 7 { 
		endCol = 7
		colStep = 1
	}

	for col := 4; col != endCol; col = col+colStep {
		sq, sqErr := cg.getSquare(moveFromRow, col)
		if sqErr != nil {
			err = sqErr
			return
		}
		squaresToEval = append(squaresToEval, sq)
	}

	evalCastleSquaresErr := cg.evalCastleSquares(squaresToEval)
	if evalCastleSquaresErr != nil {
		err = evalCastleSquaresErr
		return
	}

	return
}

func (cg *chessGame) evalCastleSquares(squaresToEval []*square) (err error) {

	// check that all squares are empty
	for i, sq := range squaresToEval {
		
		// skip the King's square
		if i == 0 {
			continue
		}
		if sq.cp != nil {
			err = newChessError(errCollision, "Castle will collide with %s at (%d %d)", sq.cp.fullName(), sq.row, sq.col)
			return
		}
	}

	// if castling left, remove the last square since the king doesn't move there
	if len(squaresToEval) == 4 {
		squaresToEval = squaresToEval[:3]
	}

	// create a simulated board
	cgSim, err := cg.cloneGame(false)
	if err != nil {
		return
	}

	// create kingRow and kingCol vars to move King back to original position
	king, kingErr := cgSim.getKing()
	if kingErr != nil {
		err = kingErr
		return
	}
	kingRow, kingCol := king.row, king.col

	for i, sq := range squaresToEval {
		if i == 0 {
			err = cgSim.inCheck()
			if err != nil {
				return
			}
		} else {
			kingMoveTo, kingMoveToErr := cgSim.getSquare(sq.row, sq.col) 
			if kingMoveToErr != nil {
				err = kingMoveToErr
				return
			}
			
			// simulate moving the King
			king.cp.updatePosition(sq.row, sq.col)
			kingMoveTo.cp = king.cp
			king.cp = nil

			err = cgSim.inCheck()
			if err != nil {
				return
			}

			// move the king back
			kingMoveTo.cp.updatePosition(kingRow, kingCol)
			king.cp = kingMoveTo.cp
			kingMoveTo.cp = nil
		}
	}
	return
}

func (cg *chessGame) castle(moveFrom, moveTo *square) {

	var king, rook chessPiece
	var kingCol, rookCol int 
	var row int 

	row = moveTo.row // can be moveTo or moveFrom since castle occurs on one row

	if moveFrom.col == 0 || moveTo.col == 0 {
		kingCol, rookCol = 2, 3
	} else {
		kingCol, rookCol = 6, 5
	}

	if moveFrom.cp.name() == "kg" {
		king, rook = moveFrom.cp, moveTo.cp
	} else {
		king, rook = moveTo.cp, moveFrom.cp
	}

	king.updatePosition(row, kingCol)
	rook.updatePosition(row, rookCol)

	// update the squares with new pieces
	kgsq, _ := cg.getSquare(row, kingCol)
	rksq, _ := cg.getSquare(row, rookCol)
	kgsq.cp, rksq.cp = king, rook

	kgsq.cp = moveTo.cp
	rksq.cp = moveFrom.cp
	kgsq.cp.setHasMoved()
	rksq.cp.setHasMoved()

	moveFrom.cp, moveTo.cp = nil, nil
}

func (cg *chessGame) capturePiece(moveTo *square) {

	moveTo.cp.updatePosition(-1, -1)
	if moveTo.cp.color() == "W" {
		cg.wCaptured = append(cg.wCaptured, moveTo.cp)
	} else {
		cg.bCaptured = append(cg.bCaptured, moveTo.cp)
	}
}

func (cg *chessGame) recoverLastCapturedPiece(moveTo *square) (err error) {

	var captured []chessPiece
	if cg.player == "W" {
		captured = cg.bCaptured
	} else {
		captured = cg.wCaptured
	}

	if len(captured) == 0 {
		err = newChessError(errNoCapturedPieces, "No captured pieces to recover")
		return
	}

	moveTo.cp, captured = captured[len(captured)-1], captured[:len(captured)-1]
	moveTo.cp.updatePosition(moveTo.row, moveTo.col)

	return

}

func (cg *chessGame) promotePawn(moveTo *square) {

	// will need a way to pause the game and have player select the new piece
	majorMinor := []string{"rk","kt","bp","qn"}
	n := rand.Intn(5)
	constructor := pieces[majorMinor[n]]
	cg.board[moveTo.row][moveTo.col] = square{moveTo.row, moveTo.col, constructor(majorMinor[n] + "-" + moveTo.cp.color() + "-0", moveTo.row, moveTo.col)}
}

func (cg *chessGame) checkCollision(moveFromRow, moveFromCol, moveToRow, moveToCol int, moveFromFullName string) (err error) {
	
	rowStep, colStep := 0, 0

	if moveFromRow != moveToRow {
		rowStep = (moveToRow - moveFromRow) / abs(moveToRow - moveFromRow)
	}
	if moveFromCol != moveToCol {
		colStep = (moveToCol - moveFromCol) / abs(moveToCol - moveFromCol)
	}

	for {
		if moveFromRow != moveToRow {
			moveFromRow += rowStep
		}
		if moveFromCol != moveToCol {
			moveFromCol += colStep
		}
		if moveFromRow == moveToRow && moveFromCol == moveToCol {
			break
		}

		if cg.board[moveFromRow][moveFromCol].cp != nil {
			msg := fmt.Sprintf("%s collides at (%d, %d)", moveFromFullName, moveFromRow, moveFromCol)
			err = fmt.Errorf(msg)
			return
		}
	}

	return
}

func (cg *chessGame) makeMove(moveFrom, moveTo *square) (err error) {

	// does the moveFrom square contain a chessPiece?
	if moveFrom.cp == nil {
		err = newChessError(errEmptySquare, "Square (%d %d) is empty. Choose a square with your piece", moveFrom.row, moveFrom.col)
		return
	}

	// does the basePiece belong to the player?
	if !(cg.player == moveFrom.cp.color()) {
		err = newChessError(errOccupiedSquare, "Square (%d %d) is occupied by %s. %s cannot move %s pieces. Choose a square with your piece", 
			moveFrom.row, moveFrom.col, moveFrom.cp.fullName(), cg.player, moveFrom.cp.color())
		return
	}

	// is the moveTo square occupied by the current player's piece?
	if moveTo.cp != nil && moveTo.cp.color() == moveFrom.cp.color() {
		moveToName := moveTo.cp.name()
		moveFromName := moveFrom.cp.name()

		// raise error if the move doesn't look like a castle
		if !((moveToName == "kg" && moveFromName == "rk") || (moveToName == "rk" && moveFromName == "kg")) {

			err = newChessError(errOccupiedSquare, "Square (%d %d) is occupied by your piece %s. Choose a valid square to move to.", 
				moveTo.row, moveTo.col, moveFrom.cp.fullName())
			return
		}
	}

	// get the implementation of chessPiece and call its specific isValidMove 
	enPassant, promotePawn, canCastle, err := moveFrom.cp.isValidMove(moveTo, cg)
	if err == nil {
		
		// if canCastle, then castle
		if canCastle {
			cg.castle(moveFrom, moveTo)
		} else {
			
			// set the from piece's new position
			moveFrom.cp.updatePosition(moveTo.row, moveTo.col)

			// if en passant, capture the prior piece, otherwise, if space is occupied by oppenent, capture its piece
			if enPassant {
				cg.capturePiece(moveToPrior)
				moveToPrior.cp = nil
			} else if moveTo.cp != nil {
				cg.capturePiece(moveTo)
			}
			
			// update the squares
			moveTo.cp = moveFrom.cp
			moveFrom.cp = nil
			moveTo.cp.setHasMoved()
			
			// if promotePawn, then promote it
			if promotePawn {
				cg.promotePawn(moveTo)
			}
		}
	}

	// update moveFromPrior and moveToPrior
	moveFromPrior, moveToPrior = moveFrom, moveTo

	return
}

func (cg *chessGame) makeComputerMove() (move []*square, err error) {
	
	// gather all spaces that are unoccupied or occupied by opponent to potentially move to
	compMoves := make([][]*square, 0)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveTo, moveToErr := cg.getSquare(row, col)
			if moveToErr != nil {
				err = moveToErr
				return
			}
			if moveTo.cp == nil || (moveTo.cp != nil && moveTo.cp.color() != cg.player) {

				// gather all comp occupied squares to potentially move from
				for rowFrom := 0; rowFrom < 8; rowFrom++ {
					for colFrom := 0; colFrom < 8; colFrom++ {
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

	if len(compMoves) < 1 {
		err = newChessError(errNoValidCompMoves, "No valid comp moves found!")
		return
	} else {
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(compMoves))
		randMove := compMoves[randIndex]
		err = cg.makeMove(randMove[0], randMove[1])
		if err != nil {
			handleError(err)
		}
	}
	
	return
}

func (p *pawn) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(p.row, p.col, moveTo)

	// ensure pawn only moves forward
	if (p.color() == "W" && moveTo.row >= p.row) || (p.color() == "B" && moveTo.row <= p.row) {
		err = newChessError(errInvalidMove, "Pawns can only move forward.")
		return
	}

	// ensure pawn only moves one row at a time, but allow two rows off starting row
	advance := abs(verticalDistance)
	if advance > 1 {
		if advance == 2 && (p.row == 1 || p.row == 6) {
			// pawn is allowed to move two spaces forward off starting row
		} else {
			err = newChessError(errInvalidMove, "Pawns can only move one space at a time, unless moving two spaces from starting row.")
			return
		}
	}

	// determine if pawn should be promoted
	if (p.color() == "W" && moveTo.row == 0) || (p.color() == "B" && moveTo.row == 7) {
		promotePawn = true
	}

	// ensure pawn move is in the same column, unless capturing an opponent's piece or en passant
	lateral := abs(horizontalDistance)
	if lateral == 1 {

		// pawn is allowed to capture opponent's pieces
		if moveTo.cp != nil && moveTo.cp.color() != p.color() {
			return 

		// evaluate en passant
		} else if (moveToPrior.cp != nil &&
			moveToPrior.cp.name() == "pn" && 
			abs(moveToPrior.row - moveFromPrior.row) == 2 &&
			moveToPrior.row == p.row &&
			moveToPrior.col == moveTo.col) {
				enPassant = true
				return 
		} else {
			err = newChessError(errInvalidMove, "Pawns cannot move laterally unless capturing")
			return
		}
	} else if lateral > 1 {
		err = newChessError(errInvalidMove, "Pawns cannot move laterally more than one square")
		return
	}

	// is the moveTo space occupied?
	if moveTo.cp != nil {
		err = newChessError(errOccupiedSquare, "Square (%d %d is occupied by %s. Choose a valid square to move to.",
			moveTo.row, moveTo.col, moveTo.cp.fullName())
	}

	return
}

func (r *rook) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(r.row, r.col, moveTo)

	if !(isVerticalOnly != isHorizontalOnly) { // xor logic
		err = newChessError(errInvalidMove, "Rook can only move along a single row or column")
		return
	}
	
	// check for collision
	err = cg.checkCollision(r.row, r.col, moveTo.row, moveTo.col, r.fullName()) 
	if err != nil {
		return
	}

	// check if castling
	if isHorizontalOnly && moveTo.cp != nil && moveTo.cp.name() == "kg" && r.color() == moveTo.cp.color() {
		
		err = moveTo.cp.getHasMoved()
		if err != nil {
			return
		}

		err = r.getHasMoved()
		if err != nil {
			return
		}
		
		evalCastleErr := cg.evalCastle(r.row, r.col, moveTo.col, r.fullName())
		if evalCastleErr != nil {
			err = evalCastleErr
			return 
		}

		// the rook can castle
		canCastle = true
	}

	return
}

func (k *knight) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(k.row, k.col, moveTo)

	verticalMove := abs(verticalDistance)
	horizontalMove := abs(horizontalDistance)

	// knight can only move 3 squares; 2 vertical and 1 lateral, or 1 vertical and 2 lateral
	if !((verticalMove == 2 && horizontalMove == 1) || (verticalMove == 1 && horizontalMove == 2)) {
		err = newChessError(errInvalidMove, "Knights must move both vertically and horizontally, and a total of 3 squares")
		return
	}

	// makeMove checks moveTo is nil or occuppied by opponent -- no need to call checkCollision here

	return
}

func (b *bishop) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(b.row, b.col, moveTo)

	// bishop must move the same number of squares both vertically and horizontally
	if !isValidDiagonal {
		err = newChessError(errInvalidMove, "Bishops must move both vertically and horizontally the same number of squares")
		return
	}

	// check for collision
	err = cg.checkCollision(b.row, b.col, moveTo.row, moveTo.col, b.fullName())
	if err != nil {
		return
	}

	return
}

func (q *queen) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(q.row, q.col, moveTo)

	if !(isVerticalOnly != isHorizontalOnly) {

		// check diagonal movement
		if !isValidDiagonal {
			err = newChessError(errInvalidMove, "When moving vertically and horizontally, queens must move the same number of squares")
			return
		}
	}

	// check for collision
	err = cg.checkCollision(q.row, q.col, moveTo.row, moveTo.col, q.fullName())
	if err != nil { 
		return
	}

	return 
}

func (k *king) isValidMove(moveTo *square, cg *chessGame) (enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	enPassant, promotePawn, canCastle = false, false, false
	setChessMove(k.row, k.col, moveTo)

	// check if the king can castle
	if isHorizontalOnly && moveTo.cp != nil && moveTo.cp.name() == "rk" && k.color() == moveTo.cp.color() {

		err = moveTo.cp.getHasMoved()
		if err != nil {
			return
		}

		err = k.getHasMoved()
		if err != nil {
			return
		}
		
		evalCastleErr := cg.evalCastle(k.row, k.col, moveTo.col, k.fullName())
		if evalCastleErr != nil {
			err = evalCastleErr
			return 
		}

		// the king can castle
		canCastle = true
	} else {

		// check for collision
		err = cg.checkCollision(k.row, k.col, moveTo.row, moveTo.col, k.fullName())
		if err != nil {
			return
		}

		// clone the game
		cgSim, cgSimErr := cg.cloneGame(false)
		if cgSimErr != nil {
			err = cgSimErr
			return
		}

		// get the king's square
		king, kingErr := cgSim.getKing()
		if kingErr != nil {
			err = kingErr
			return
		}

		// get moveTo square from sim game
		kingMoveTo, kingMoveToErr := cgSim.getSquare(moveTo.row, moveTo.col)
		if kingMoveToErr != nil {
			err = kingMoveToErr
			return
		}

		// simulate moving the king
		king.cp.updatePosition(moveTo.row, moveTo.col)
		kingMoveTo.cp = king.cp
		king.cp = nil

		// if the king is in check, return the error
		err = cgSim.inCheck()
		if err != nil {
			return
		}

		if (abs(verticalDistance) > 1 || abs(horizontalDistance) > 1) {
			err = newChessError(errInvalidMove, "Kings can only move one square at a time, unless castling")
			return
		}
	}

	return
}