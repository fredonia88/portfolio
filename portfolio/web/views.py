import json
from .models import (
    Home, 
    Contact,
    Projects
)
from .tic_tac_toe import TicTacToe, MoveIsTaken
from django.views.generic import ListView
from django.views import View
from django.shortcuts import render
from django.http import JsonResponse

class HomeView(ListView):
    model = Home
    template_name = 'home.html'

class ContactView(ListView):
    model = Contact
    template_name = 'contact.html'

class About(ListView):
    model = Contact
    template_name = 'about.html'

class ErcotView(ListView):
    model = Projects
    template_name = 'ercot.html'

class BlsView(ListView):
    model = Projects
    template_name = 'bls.html'

class TicTacToeView(View):
    template_name = 'tictactoe.html'

    def get(self, request):
        return render(request, self.template_name)

class TicTacToeBoard(View):
    template_name = 'tictactoe.html'

    def get(self, request):
        game = TicTacToe()
        request.session['board'] = game.board
        return JsonResponse({'board': game.board})
    
    def post(self, request):
        try:
            board = request.session.get('board')
            if board is None:
                raise ValueError('Game not found')
            
            game = TicTacToe(board=board)
            body = json.loads(request.body)
            if body == 'compMove':
                response = game.comp_move
            else:
                row = int(body['row'])
                col = int(body['col'])
                response = game.user_move(row, col)

            winner = game.victory_for()
            request.session['board'] = game.board
            return JsonResponse({'board': game.board, 'winner': winner})
        except MoveIsTaken as e:
            response = {'error': str(e)}
            return JsonResponse(response, status=400)
        