from django.conf.urls import url

from . import views

app_name = 'hosts'
urlpatterns = [
	url(r'^$', views.index, name='index'),
	url(r'view/(?P<identifier>\w+)/$', views.profile, name='profile'),
]
