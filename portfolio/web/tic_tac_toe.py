import random

class NumberMustBeInRange(Exception):
    '''used in enter_move method -> input must be an integer'''
    pass

class MoveIsTaken(Exception):
    '''used in enter_move method -> move must be available'''
    pass

class TicTacToe:
    
    def __init__(self, board=None):
        if board:
            self.board = board
        else:
            self.board = [['' for c in range(3)] for r in range(3)]
            self.comp_move
        
    @property
    def comp_move(self):
        free_spaces = self.make_list_of_free_fields()
        space = random.randint(a=0, b=len(free_spaces) - 1)
        row, col = free_spaces[space][0], free_spaces[space][1]
        self.board[row][col] = 'X'
    
    def user_move(self, row, col):
        if not 0 <= row < 3 or not 0 <= col < 3:
            raise NumberMustBeInRange('Row or Column is out of range! Select a square in range!')
        elif (row, col) not in self.make_list_of_free_fields():
            raise MoveIsTaken('Move is taken! Select a different square!')
        else:
            self.board[row][col] = 'O'
    
    def victory_for(self):
        if not self.make_list_of_free_fields():
            return 'Tie'

        for sign in ['X', 'O']:
            win_list = [[sign] * 3]
            
            for x in range(3):
                if self.board[x] in win_list:
                    return sign
                    
            for y in range(3):
                vertical_list = []
                for x in range(3):
                    vertical_list.append(self.board[x][y])
                if vertical_list in win_list:
                    return sign
                    
            if [self.board[0][0], self.board[1][1], self.board[2][2]] in win_list or [self.board[0][2], self.board[1][1], self.board[2][0]] in win_list:
                return sign
                
        return None

    def make_list_of_free_fields(self):
        free_spaces = []
        for x in range(len(self.board)):
            for y in range(len(self.board[x])):
                if self.board[x][y] not in ('X', 'O'):
                    free_spaces.append((x, y))

        return free_spaces