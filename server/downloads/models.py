# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models

class Download(models.Model):
	release_date = models.DateTimeField(auto_now_add=True, blank=True)
	version = models.CharField(max_length=255, null=True)
	url = models.TextField()
	sha1 = models.CharField(max_length=40, default="")
	downloads_count = models.IntegerField(null=True)

	def __str__(self):
		return "%s %s" % (self.release_date, self.sha1)
