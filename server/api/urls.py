from django.conf.urls import url

from . import views

urlpatterns = [
	url(r'^versioncheck/$', views.versioncheck),
	url(r'^register/$', views.register),
	url(r'^heartbeat/$', views.heartbeat),
	url(r'^detection/(?P<identifier>\w{40})/$', views.detection),
	url(r'^autorun/(?P<identifier>\w{40})/$', views.autorun),
]
