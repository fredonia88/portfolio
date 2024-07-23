from django.db import models

class TicTacToeResult(models.Model):
    winner = models.CharField(max_length=1)
    created_at = models.DateTimeField(auto_now_add=True)
