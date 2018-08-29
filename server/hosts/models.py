# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models

class Host(models.Model):
	identifier = models.CharField(max_length=255)
	user_name = models.CharField(max_length=255, null=True)
	computer_name = models.CharField(max_length=255, null=True)
	operating_system = models.CharField(max_length=255, null=True)
	version = models.CharField(max_length=10, null=True)
	last_seen = models.DateTimeField(auto_now_add=True, blank=True)

	def __str__(self):
		return "%s (%s)" % (self.computer_name, self.identifier)

class Heartbeat(models.Model):
	host = models.ForeignKey("Host", on_delete=models.CASCADE)
	date = models.DateTimeField(auto_now_add=True, blank=True)
	host_ip = models.GenericIPAddressField(null=True, default=None)
	country = models.CharField(max_length=2, null=True)
	city = models.CharField(max_length=50, null=True)

	def __str__(self):
		return "%s %s" % (self.date, self.host_ip)
