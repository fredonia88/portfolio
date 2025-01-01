from django.urls import path
from .views import (
    HomeView, 
    ContactView,
    AboutView,
    TicTacToeView,
    TicTacToeBoardView,
    ErcotView,
    BlsView,
    BlsChartView,
    BlsPercentChangeChartView
)

urlpatterns = [
    path('', HomeView.as_view(), name='home'),
    path('contact', ContactView.as_view(), name='contact'),
    path('about', AboutView.as_view(), name='about'),
    path('tictactoe/', TicTacToeView.as_view(), name='tictactoe'),
    path('tictactoe/board', TicTacToeBoardView.as_view(), name='tictactoe_board'),
    #path('ercot', ErcotView.as_view(), name='ercot'),
    path('bls/', BlsView.as_view(), name='bls'),
    path('bls/chart-data', BlsChartView.as_view(), name='bls_chart_data'),
    path('bls/median-income-percent-change', BlsPercentChangeChartView.as_view(), name='bls_median_income_percent_change'),
]