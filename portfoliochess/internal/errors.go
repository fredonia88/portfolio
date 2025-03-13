package main

import (
	"fmt"
	"errors"
	"runtime"
	"strings"
	"log"
)

var (
	errCollision = errors.New("COLLISION")
	errOutOfRange = errors.New("OUT_OF_RANGE")
	errEmptySquare = errors.New("EMPTY_SQUARE")
	errOccupiedSquare = errors.New("OCCUPIED_SQUARE")
	errInvalidMove = errors.New("INVALID_MOVE")
	errHasMovedConversion = errors.New("HAS_MOVED_CONVERSION")
	errHasMoved = errors.New("HAS_MOVED")
	errMissingPiece = errors.New("MISSING_PIECE")
	errInCheck = errors.New("IN_CHECK")
	errKingNotFound = errors.New("KING_NOT_FOUND")
	errNoCapturedPieces = errors.New("NO_CAPTURED_PIECES")
	errNoValidCompMoves = errors.New("NO_VALID_COMP_MOVES")
	errCheckMate = errors.New("CHECKMATE")
)

type chessError struct {
	details error
	message string
	file string
	line int
}

func handleError(err error) {
	if err != nil {
		log.Fatal("Error ", err)
	}
}

func (e *chessError) Error() string {
	return fmt.Sprintf("[%v]: %s (%s: %d)", e.details, e.message, e.file, e.line)
}

func newChessError(details error, message string, args ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	subdirs := strings.Split(file, "/")
	file = subdirs[len(subdirs)-1]

	return &chessError{
		details,
		fmt.Sprintf(message, args...),
		file,
		line,
	}
}
