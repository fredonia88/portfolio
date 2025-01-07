import json
from .models import (
    TicTacToeResult,
    MedianIncomeByAgeConstantDollars,
    MedianIncomePercentChangeByAgeConstantDollars
)
from .forms import ContactForm
from .tic_tac_toe import TicTacToe, MoveIsTaken
from django.views import View
from django.shortcuts import render, redirect
from django.core.mail import send_mail
from django.conf import settings
from django.http import JsonResponse
from django.contrib import messages

class HomeView(View):
    template_name = 'home.html'

    def get(self, request):

        return render(request, self.template_name)

class ContactView(View):
    template_name = 'contact.html'
    form = ContactForm

    def get(self, request):
        context = {
            'title': 'Contact',
            'form': self.form
        }

        return render(request, self.template_name, context)
    
    def post(self, request):
        submission = self.form(request.POST)
        if submission.is_valid():
            name = submission.cleaned_data['name']
            email = submission.cleaned_data['email']
            message = submission.cleaned_data['message']

            send_mail(
                subject=f'New Portfolio Email From {name}!',
                message=f'Email: {email}\nMessage: {message}',
                from_email=settings.EMAIL_HOST_USER,
                recipient_list=[f'{settings.EMAIL_RECIPIENT}'],
                fail_silently=False,
            )
            messages.success(request, '- email sent!')

            return redirect('contact')
        
        else:
            if 'captcha' in submission.errors.as_data():
                del submission.errors['captcha']
                context = {'form': submission, 'captcha_error': True}
                
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
    model = MedianIncomePercentChangeByAgeConstantDollars

    def get(self, request):

        data = self.model.objects.values('demographic_age', 'starting_value_constant_dollars', 'ending_value_constant_dollars', 'percent_change_in_income')
        data = [entry for entry in data]

        context = {
            'title': 'Bls',
            'data': data
        }

        return render(request, self.template_name, context)


class BlsChartView(View):
    template_name = 'bls.html'
    model = MedianIncomeByAgeConstantDollars

    def get(self, request):

        datasets = {}
        data = self.model.objects.values('year', 'demographic_age', 'yearly_value_constant_dollars')
        for entry in data:
            year = str(entry['year'])
            age = entry['demographic_age'].replace(' years', '') # annoying bug -- chart.js can't handle length of label values
            income = entry['yearly_value_constant_dollars']

            if age not in datasets:
                datasets[age] = {'label': age, 'data': [], 'borderColor': '', 'fill': False, 'tension': 0.1}
            datasets[age]['data'].append({'x': year, 'y': income})

        return JsonResponse(list(datasets.values()), safe=False)

class BlsPercentChangeChartView(View):
    template_name = 'bls.html'
    model = MedianIncomePercentChangeByAgeConstantDollars

    def get(self, request):

        data = self.model.objects.values('demographic_age', 'starting_value_constant_dollars', 'ending_value_constant_dollars', 'percent_change_in_income')
        data = [entry for entry in data]

        return JsonResponse(data, safe=False)


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
        