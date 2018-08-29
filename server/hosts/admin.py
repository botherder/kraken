# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.contrib import admin

from .models import Host

@admin.register(Host)
class HostAdmin(admin.ModelAdmin):
	pass
