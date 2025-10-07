from django.urls import path

from organizations.views.create_organization import CreateOrganization

urlpatterns = [
    path('create/', CreateOrganization.as_view(), name='create_organization'),
]
