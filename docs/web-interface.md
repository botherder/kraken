# Installing the Web Interface

The web interface is built using Django. You can run it using Python 3, which will require the following dependencies:

    $ sudo apt install python3 python3-dev python3-pip python3-mysqldb
    $ sudo pip3 install Django python-decouple django-geoip2-extras

To configure your Krakan Django app, instead of modifying `server/settings.py` you can create a file named `.env` inside the `server/` folder with the following content:

```shell
SECRET_KEY=your_secret_key
DEBUG=True
DB_NAME=kraken
DB_USER=user
DB_PASSWORD=pass
STATIC_ROOT=/home/user/kraken/server/static/
GEOIP_PATH=/home/user/geoip/
```

Change those values appropriately. The `GEOIP_PATH` variable should point to a folder containing your [MaxMind GeoLite2 City](https://dev.maxmind.com/geoip/geoip2/geolite2/) database.

After having configured the settings in the `.env` file, you will need to initialize the database with:

    $ python3 manage.py makemigrations autoruns detections downloads hosts
    $ python3 manage.py migrate

If you want to run the server using Gunicorn, you can install it with:

    $ sudo pip3 install gunicorn

You can create a Gunicorn systemd service by creating a `kraken.service` file in `/etc/systemd/system` like the following:

    Description=Gunicorn Application Server handling Kraken Servers
    After=network.target

    [Service]
    User=user
    Group=www-data
    WorkingDirectory=/home/user/kraken/server/
    ExecStart=/usr/local/bin/gunicorn --workers 3 --bind unix:/home/user/kraken-server.sock server.wsgi:application
    Restart=always

    [Install]
    WantedBy=multi-user.target

You can then configure your webserver to proxy requests to the unix socket at `/home/user/kraken-server.sock`.
