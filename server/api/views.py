# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from django.contrib.gis.geoip2 import GeoIP2
from django.http import JsonResponse

from hosts.models import Host, Heartbeat
from detections.models import Detection
from autoruns.models import Autorun
from downloads.models import Download

import json

STATUS_SUCCESS = "OK"
STATUS_FAIL = "FAIL"

OK_DOWNLOAD = "OK_DOWNLOAD"
OK_AUTORUN_ALREADY_STORED = "OK_AUTORUN_ALREADY_STORED"
ERR_INVALID_METHOD = "ERR_INVALID_METHOD"
ERR_NOT_REGISTERED = "ERR_NOT_REGISTERED"

CODES = {
	OK_DOWNLOAD: "New version avilable for download",
	OK_AUTORUN_ALREADY_STORED: "This autorun record appears to have been stored already",
	ERR_INVALID_METHOD: "Invalid HTTP method",
	ERR_NOT_REGISTERED: "Host is not registered",
}

@csrf_exempt
def versioncheck(request):
	if request.method != "POST":
		return JsonResponse({"status": STATUS_FAIL,
			"code": ERR_INVALID_METHOD, "msg": CODES[ERR_INVALID_METHOD]})

	body_unicode = request.body.decode("utf-8")
	body = json.loads(body_unicode)

	current_sha1 = body["sha1"]

	latest_download = Download.objects.order_by("-release_date").first()

	if not latest_download or (current_sha1 == latest_download.sha1):
		return JsonResponse({
			"status": STATUS_SUCCESS,
			"code": "",
			"msg": "",
			"url": "",
		})
	else:
		return JsonResponse({
			"status": STATUS_SUCCESS,
			"code": OK_DOWNLOAD,
			"msg": CODES[OK_DOWNLOAD],
			"url": latest_download.url,
		})

@csrf_exempt
def register(request):
	if request.method != "POST":
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_INVALID_METHOD,
			"msg": CODES[ERR_INVALID_METHOD]})

	body_unicode = request.body.decode("utf-8")
	body = json.loads(body_unicode)

	host, created = Host.objects.update_or_create(
		identifier=body["identifier"],
		defaults={"user_name": body["user_name"], "computer_name": body["computer_name"],
		"operating_system": body["operating_system"], "version": body["version"]},
	)

	return JsonResponse({"status": STATUS_SUCCESS, "code": "", "msg": ""})

@csrf_exempt
def heartbeat(request):
	if request.method != "POST":
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_INVALID_METHOD,
			"msg": CODES[ERR_INVALID_METHOD]})

	# Because of the use of an Nginx proxy, we need to use this
	# trick to obtain the real IP address.
	host_ip = request.META.get("REMOTE_ADDR")
	if not host_ip:
		host_ip = request.META.get("HTTP_X_REAL_IP")
	if not host_ip:
		host_ip = request.META.get("HTTP_X_FORWARDED_FOR")

	# Initialize MaxMind database.
	geo = GeoIP2()
	# Look up the host IP address.
	geo_data = geo.city(host_ip)

	body_unicode = request.body.decode("utf-8")
	body = json.loads(body_unicode)

	try:
		host = Host.objects.get(identifier=body["identifier"])
	except Host.DoesNotExist:
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_NOT_REGISTERED,
			"msg": CODES[ERR_NOT_REGISTERED]})

	# Create new heartbeat record.
	heartbeat = Heartbeat(
		host=host,
		host_ip=host_ip,
		country=geo_data["country_code"],
		city=geo_data["city"],
	)
	heartbeat.save()

	# Save last seen datetime.
	host.last_seen = heartbeat.date
	host.save()

	# TODO: Add response in JSON providing URL to download update,
	# if new version.

	return JsonResponse({"status": STATUS_SUCCESS, "code": "", "msg": ""})

@csrf_exempt
def detection(request, identifier):
	if request.method != "POST":
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_INVALID_METHOD,
			"msg": CODES[ERR_INVALID_METHOD]})

	body_unicode = request.body.decode("utf-8")
	body = json.loads(body_unicode)

	try:
		host = Host.objects.get(identifier=identifier)
	except Host.DoesNotExist:
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_NOT_REGISTERED,
			"msg": CODES[ERR_NOT_REGISTERED]})

	detection = Detection(
		host=host,
		record_type=body["type"],
		signature=body["signature"],
		process_id=body["process_id"],
		image_name=body["image_name"],
		image_path=body["image_path"],
		sha1=body["sha1"],
		sha256=body["sha256"],
	)
	detection.save()

	return JsonResponse({"status": STATUS_SUCCESS, "code": "", "msg": ""})

@csrf_exempt
def autorun(request, identifier):
	if request.method != "POST":
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_INVALID_METHOD,
			"msg": CODES[ERR_INVALID_METHOD]})

	body_unicode = request.body.decode("utf-8")
	body = json.loads(body_unicode)

	try:
		host = Host.objects.get(identifier=identifier)
	except Host.DoesNotExist:
		return JsonResponse({"status": STATUS_FAIL, "code": ERR_NOT_REGISTERED,
			"msg": CODES[ERR_NOT_REGISTERED]})

	try:
		autorun = Autorun.objects.get(host=host, autorun_type=body["type"],
			image_path=body["image_path"], arguments=body["arguments"])
	except Autorun.DoesNotExist:
		autorun = Autorun(
			host=host,
			autorun_type=body["type"],
			image_name=body["image_name"],
			image_path=body["image_path"],
			arguments=body["arguments"],
			sha1=body["sha1"],
			sha256=body["sha256"],
		)
		autorun.save()

		return JsonResponse({"status": STATUS_SUCCESS, "code": "", "msg": ""})
	else:
		return JsonResponse({"status": STATUS_SUCCESS,
			"code": OK_AUTORUN_ALREADY_STORED,
			"msg": CODES[OK_AUTORUN_ALREADY_STORED]})
