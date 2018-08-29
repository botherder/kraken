# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models

class Autorun(models.Model):
	host = models.ForeignKey("hosts.Host", on_delete=models.CASCADE, blank=True, null=True)
	date = models.DateTimeField(auto_now_add=True, blank=True)
	autorun_type = models.CharField(max_length=255)
	image_name = models.CharField(max_length=100)
	image_path = models.TextField()
	arguments = models.TextField()
	sha1 = models.CharField(max_length=40, null=True)
	sha256 = models.CharField(max_length=64, null=True)

	def __str__(self):
		return "%s %s" % (self.autorun_type, self.image_name)
