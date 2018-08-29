# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.shortcuts import render

from .models import Detection

# This is to browse through all detections.
def index(request):
	detections = Detection.objects.all().order_by("-date")
	return render(request, "detections.html", {"detections": detections})
