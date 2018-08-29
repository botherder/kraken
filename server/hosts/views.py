# -*- coding: utf-8 -*-
from __future__ import unicode_literals

import json

from django.shortcuts import render
from django.http import Http404

from .models import Host, Heartbeat
from detections.models import Detection
from autoruns.models import Autorun

# This is the profile for the host.
def profile(request, identifier):
	try:
		host = Host.objects.get(identifier=identifier)
	except Host.DoesNotExist:
		raise Http404("Host not found. Was it registered?")

	heartbeats = Heartbeat.objects.filter(host=host).order_by("-date")[:10]
	first_seen = Heartbeat.objects.filter(host=host).order_by("date").first()
	detections = Detection.objects.filter(host=host).order_by("-date")
	autoruns = Autorun.objects.filter(host=host).order_by("autorun_type", "-date")

	return render(request, "host.html", {"host": host, "heartbeats": heartbeats,
		"first_seen": first_seen, "detections": detections, "autoruns" : autoruns})

# This is to browse through all hosts.
def index(request):
	hosts = Host.objects.all().order_by("-last_seen", "computer_name")
	return render(request, "hosts.html", {"hosts": hosts})
