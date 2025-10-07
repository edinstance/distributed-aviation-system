from django.urls import path

from .views.create import create

urlpatterns = [
    path('create/', create, name='create'),
]
