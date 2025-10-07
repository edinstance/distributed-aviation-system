from django.urls import path

from .jwks.jwks_view import jwks_view
from .views.login import Login
from .views.logout import Logout
from .views.refresh_token import Refresh
from .views.verify_token import VerifyToken

urlpatterns = [
    path("login/", Login.as_view(), name="login"),
    path("logout/", Logout.as_view(), name="logout"),
    path("refresh/", Refresh.as_view(), name="token_refresh"),
    path("verify-token/", VerifyToken.as_view(), name="verify_token"),
    path("jwks.json", jwks_view, name="jwks"),
]