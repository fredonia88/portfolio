package main

import (
	"fmt"
	"log"
	"math/rand"
)

/*
	TODO:
		1. Need a more elegant way to handle errors
		2. Check for checkmate
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

func (cg *chessGame) castleThruSquares(moveFromRow, moveFromCol, moveToCol int, moveFromFullName string) (squaresToEval []*square, err error) {

	// collect all squares the king would start at, move through or end on while castling
	squaresToEval = make([]*square, 0, 3)

	endCol := 1
	colStep := -1
	if moveFromCol == 7 || moveToCol == 7 { 
		endCol = 7
		colStep = 1
	}

	for col := 4; col != endCol; col = col+colStep {
		sq, sqErr := cg.getSquare(moveFromRow, col)
		if sqErr != nil {
			msg := fmt.Sprintf("%s collides at (%d, %d)", moveFromFullName, moveFromRow, col)
			err = fmt.Errorf(msg)
			return
		}
		squaresToEval = append(squaresToEval, sq)
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

func (cg *chessGame) promotePawn(moveTo *square) {

	// will need a way to pause the game and have player select the new piece
	majorMinor := []string{"rk","kt","bp","qn"}
	n := rand.Intn(5)
	constructor := pieces[majorMinor[n]]
	cg.board[moveTo.row][moveTo.col] = square{moveTo.row, moveTo.col, constructor(majorMinor[n] + "-" + moveTo.cp.color() + "-0", moveTo.row, moveTo.col)}
}

func (cg *chessGame) checkCollision(moveFromRow, moveFromCol, moveToRow, moveToCol int, moveFromFullName string) (doesCollide bool, err error) {
	
	doesCollide = true

	rowStep, colStep := 0, 0

	if moveFromRow != moveToRow {
		rowStep = (moveToRow - moveFromRow) / abs(moveToRow - moveFromRow)
	}
	if moveFromCol != moveToCol {
		colStep = (moveToCol - moveFromCol) / abs(moveToCol - moveFromCol)
	}

	for row, col := moveFromRow+rowStep, moveFromCol+colStep; row != moveToRow || col != moveToCol; row, col = row+rowStep, col+colStep {
		if cg.board[row][col].cp != nil {
			msg := fmt.Sprintf("%s collides at (%d, %d)", moveFromFullName, row, col)
			err = fmt.Errorf(msg)
			return
		}
	}

	doesCollide = false

	return
}

func (cg *chessGame) makeMove(moveFrom, moveTo *square) {

	// does the moveFrom square contain a chessPiece?
	if moveFrom.cp == nil {
		err := fmt.Errorf("This square is empty. Choose a square with your piece")
		log.Fatal("Error:", err)
	}

	// does the basePiece belong to the player?
	if !(cg.player == moveFrom.cp.color()) {
		err := fmt.Errorf("%s cannot move %s color pieces", cg.player, moveFrom.cp.color())
		log.Fatal("Error:", err)
	}

	// is the moveTo square occupied by the current player's piece?
	if moveTo.cp != nil && moveTo.cp.color() == moveFrom.cp.color() {
		moveToName := moveTo.cp.name()
		moveFromName := moveFrom.cp.name()

		// raise error if the move doesn't look like a castle
		if !((moveToName == "kg" && moveFromName == "rk") || (moveToName == "rk" && moveFromName == "kg")) {
			err := fmt.Errorf("moveTo square is occupied with: %s", moveTo.cp.fullName())
			log.Fatal("Error:", err)
		}
	}

	// get the implementation of chessPiece and call its specific isValidMove 
	if canMove, enPassant, promotePawn, canCastle, err := moveFrom.cp.isValidMove(moveTo, cg); canMove {
		
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
				moveTo.cp.updatePosition(-1, -1)
				cg.capturePiece(moveTo)
			}
			
			// update the squares
			moveTo.cp = moveFrom.cp
			moveFrom.cp = nil
			
			// if promotePawn, then promote it
			if promotePawn {
				cg.promotePawn(moveTo)
			}
		}
	} else if err != nil {
		log.Fatal("Error:", err)
	}
}

func (p *pawn) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(p.row, p.col, moveTo)

	// ensure pawn only moves forward
	if (p.color() == "W" && moveTo.row >= p.row) || (p.color() == "B" && moveTo.row <= p.row) {
		err = fmt.Errorf("%s can only move forward", p.fullName())
		return
	}

	// ensure pawn only moves one row at a time, but allow two rows off starting row
	advance := abs(verticalDistance)
	if advance > 1 {
		if advance == 2 && (p.row == 1 || p.row == 6) {
			// pawn is allowed to move two spaces forward off starting row
		} else {
			err = fmt.Errorf("%s can only move one space at a time, or two from starting row", p.fullName())
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
			canMove = true
			return 

		// evaluate en passant
		} else if (moveToPrior.cp != nil &&
			moveToPrior.cp.name() == "pn" && 
			abs(moveToPrior.row - moveFromPrior.row) == 2 &&
			moveToPrior.row == p.row &&
			moveToPrior.col == moveTo.col) {
				canMove = true
				enPassant = true
				return 
		} else {
			err = fmt.Errorf("%s cannot move laterally unless capturing or using en passant", p.fullName())
			return
		}
	} else if lateral > 1 {
		err = fmt.Errorf("%s cannot move laterally more than one square", p.fullName())
		return
	}

	canMove = true

	return
}

func (r *rook) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(r.row, r.col, moveTo)

	if !(isVerticalOnly != isHorizontalOnly) { // xor logic
		err = fmt.Errorf("%s can only move along a single row or column", r.fullName())
		return
	}
	
	// check for collision
	if doesCollide, collideErr := cg.checkCollision(r.row, r.col, moveTo.row, moveTo.col, r.fullName()); doesCollide {
		err = collideErr 
		return
	}

	// check if castling
	if isHorizontalOnly && moveTo.cp != nil && moveTo.cp.name() == "kg" && r.color() == moveTo.cp.color() && 
		!moveTo.cp.hasMoved() && !r.hasMoved() {
		
		squaresToEval, _ := cg.castleThruSquares(r.row, r.col, moveTo.col, r.fullName())
		
		// check if any squaresToEval would be in check
		for s := range squaresToEval {
			if canCheck, inCheckErr, _ := cg.inCheck(squaresToEval[s]); canCheck {
				err = inCheckErr
				return
			}
		}

		// the rook can castle
		canCastle = true
	}

	canMove = true 

	return
}

func (k *knight) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(k.row, k.col, moveTo)

	verticalMove := abs(verticalDistance)
	horizontalMove := abs(horizontalDistance)

	// knight can only move 3 squares; 2 vertical and 1 lateral, or 1 vertical and 2 lateral
	if !((verticalMove == 2 && horizontalDistance == 1) || (verticalMove == 1 && horizontalMove == 2)) {
		err = fmt.Errorf("%s must move vertically and horizontally, and can only move a total of 3 squares", k.fullName())
		return
	}

	// makeMove checks moveTo is nil or occuppied by opponent -- no need to call checkCollision here

	canMove = true

	return
}

func (b *bishop) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(b.row, b.col, moveTo)

	// bishop must move the same number of squares both vertically and horizontally
	if !isValidDiagonal {
		err = fmt.Errorf("%s must move vertically and horizontally the same number of squares", b.fullName())
		return
	}

	// check for collision
	if doesCollide, collideErr := cg.checkCollision(b.row, b.col, moveTo.row, moveTo.col, b.fullName()); doesCollide {
		err = collideErr 
		return
	}
	
	canMove = true

	return
}

func (q *queen) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(q.row, q.col, moveTo)

	if !(isVerticalOnly != isHorizontalOnly) {

		// check diagonal movement
		if !isValidDiagonal {
			err = fmt.Errorf("If %s moves diagonally, it must move vertically and horizontally the same number of squares", q.fullName())
			return
		}
		err = fmt.Errorf("%s can only move along a single row or column", q.fullName())
		return
	}

	// check for collision
	if doesCollide, collideErr := cg.checkCollision(q.row, q.col, moveTo.row, moveTo.col, q.fullName()); doesCollide {
		err = collideErr 
		return
	}

	canMove = true

	return 
}

func (k *king) isValidMove(moveTo *square, cg *chessGame) (canMove, enPassant, promotePawn, canCastle bool, err error) {

	// set named return values and chessMove vars
	canMove, enPassant, promotePawn, canCastle = false, false, false, false
	setChessMove(k.row, k.col, moveTo)

	if isHorizontalOnly && moveTo.cp != nil && moveTo.cp.name() == "rk" && k.color() == moveTo.cp.color() && 
		!moveTo.cp.hasMoved() && !k.hasMoved() {
		
		squaresToEval, _ := cg.castleThruSquares(k.row, k.col, moveTo.col, k.fullName())
		
		// check if any squaresToEval would be in check
		for s := range squaresToEval {
			if canCheck, inCheckErr, _ := cg.inCheck(squaresToEval[s]); canCheck {
				err = inCheckErr
				return
			}
		}

		// the king can castle
		canCastle = true
	} else {

		if canCheck, inCheckErr, _ := cg.inCheck(moveTo); canCheck {
			err = inCheckErr
			return
		}

		if (abs(verticalDistance) > 1 || abs(horizontalDistance) > 1) {
			err = fmt.Errorf("%s can only move one square at a time", k.fullName())
			return
		}
	}

	canMove = true

	return
}