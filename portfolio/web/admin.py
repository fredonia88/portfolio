from django.contrib import admin
from .models import (
    Home, 
    Contact,
    Projects
)

# Register your models here.
admin.site.register(Home)
admin.site.register(Contact)
admin.site.register(Projects)