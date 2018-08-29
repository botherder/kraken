# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.contrib import admin

from .models import Download

@admin.register(Download)
class DownloadAdmin(admin.ModelAdmin):
	pass
