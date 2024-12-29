from django.db import models

class TicTacToeResult(models.Model):
    difficulty_level = models.CharField(max_length=12)
    winner = models.CharField(max_length=3)
    created_at = models.DateTimeField(auto_now_add=True)

class MedianIncomeByAgeConstantDollars(models.Model):
    id = models.AutoField(primary_key=True)
    year = models.BigIntegerField(blank=True, null=True)
    demographic_age = models.TextField(blank=True, null=True)
    yearly_value_constant_dollars = models.FloatField(blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'median_income_by_age_constant_dollars'