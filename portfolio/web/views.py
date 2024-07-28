import json
from .models import TicTacToeResult
from .tic_tac_toe import TicTacToe, MoveIsTaken
from django.views import View
from django.shortcuts import render
from django.http import JsonResponse

class HomeView(View):
    template_name = 'home.html'

    def get(self, request):

        return render(request, self.template_name)

class ContactView(View):
    template_name = 'contact.html'

    def get(self, request):
        context = {
            'title': 'Contact'
        }

        return render(request, self.template_name, context)

class AboutView(View):
    template_name = 'about.html'

    def get(self, request):
        context = {
            'title': 'About'
        }

        return render(request, self.template_name, context)

class ErcotView(View):
    template_name = 'ercot.html'

    def get(self, request):
        context = {
            'title': 'Ercot'
        }

        return render(request, self.template_name, context)

class BlsView(View):
    template_name = 'bls.html'

    def get(self, request):
        context = {
            'title': 'Bls'
        }

        return render(request, self.template_name, context)

class TicTacToeView(View):
    template_name = 'tictactoe.html'
    model = TicTacToeResult

    def get(self, request):
        if request.headers.get('x-requested-with') == 'XMLHttpRequest':
            return self.get_stats(request)

        context = {
            'title': 'Tic Tac Toe'
        }

        return render(request, self.template_name, context)

    def get_stats(self, request):

        difficulty_level = request.GET.get('difficulty', 'All')
        if difficulty_level == 'All':
            result = self.model.objects.all()
        else:
            result = self.model.objects.filter(difficulty_level=difficulty_level)

        games_played = result.count()
        user_wins = result.filter(winner='X').count()
        comp_wins = result.filter(winner='O').count()
        ties = result.filter(winner='Tie').count()

        if games_played > 0:
            user_win_rate = round((user_wins / games_played) * 100, 1)
            comp_win_rate = round((comp_wins / games_played) * 100, 1)
            tie_rate = round((ties / games_played) * 100, 1)
        else:
            user_win_rate = comp_win_rate = tie_rate = 0.0

        stats = {
            'difficulty_level': difficulty_level,
            'games_played': games_played,
            'user_wins': user_wins,
            'comp_wins': comp_wins,
            'ties': ties,
            'user_win_rate': user_win_rate,
            'comp_win_rate': comp_win_rate,
            'tie_rate': tie_rate
        }

        return JsonResponse(stats)

class TicTacToeBoardView(View):
    template_name = 'tictactoe.html'

    def get(self, request):
        level = request.headers['Difficulty-Level']
        game = TicTacToe(difficulty_level=level)
        request.session['board'] = game.board
        request.session['difficulty-level'] = game.difficulty_level

        return JsonResponse({'board': game.board})
    
    def post(self, request):
        try:
            board = request.session.get('board')
            level = request.session.get('difficulty-level')
            game = TicTacToe(board=board, difficulty_level=level)

            body = json.loads(request.body)
            if body == 'compMove':
                game.comp_move()
            else:
                row = int(body['row'])
                col = int(body['col'])
                game.user_move(row, col)

            winner = game.victory_for(game.board)
            if winner: 
                TicTacToeResult.objects.create(
                    winner = winner,
                    difficulty_level = level
                )

            request.session['board'] = game.board
            return JsonResponse({'board': game.board, 'winner': winner, 'difficulty_level': game.difficulty_level})
        except MoveIsTaken as e:
            response = {'error': str(e)}
            return JsonResponse(response, status=400)
        