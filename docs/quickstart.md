# Quick Start

To get Inventory up and running, you'll need two things:

* A server running `inventory server`
* One or more clients running `inventory send` periodically.

## Server

Use the [installation](install) instructions to get the binary on your server.

```bash
sh -c "$(curl --location https://bketelsen.github.io/inventory/install.sh)" -- -d -b /usr/local/bin
```

Inventory doesn't need any escalated priveleges to run the server, so I recommend using a non-root user account to run it.

Grab a copy of the example [systemd unit](https://github.com/bketelsen/inventory/blob/main/contrib/inventory-server.service), modify it to your liking, and put it in `/etc/systemd/system`.

Create the configuration directory `/etc/inventory`.  Run `inventory config -c /etc/inventory/inventory.yml` to create a configuration file.

Modify the configuration file to adjust the `server:` section to fit your environment:

```yaml
server:
    http-port: 8000  // Port for dashboard
    listen: 0.0.0.0  // Listen IP for server
    rpc-port: 9999   // Port for client connections
```

Run `systemctl daemon-reload` to inform SystemD of the new unit. Then `systemctl start inventory-server.service` to start it.  You should also run `systemctl enable inventory-server` if you want it to start automatically.

Now open a browser and point to your server's IP address at port 8000. You should see the dashboard with no data, because we haven't started any client agents yet.

## Client / Reporter

On each system that you wish to track, follow the [installation](install) instructions to install the `inventory` binary.

!> If you want to track deployments on the server you installed earlier, you can skip the installation step, the binary is the same for the server and client.

Create (or modify) a configuration file in `/etc/inventory` using the steps above. For the client, we'll want to modify the `client:` section of the configuration:

```yaml
client:
    description: UM690               // a description of your server for the dashboard
    location: Office Rack, Shelf 2   // location of the server
    remote: 192.168.5.1:9999         // the IP:PORT of your server's inventory service
```

### Permissions

The `inventory send` command will try to connect to Docker and Incus to query running containers/instances. The user that runs the `inventory send` command needs to have permissions to talk to Docker and Incus.

I reccommend creating a new service account user and adding that user to the `docker` and `incus-admin` groups.

### Sending Your First Report

From a terminal/ssh session, run the `inventory send` command. If your configuration file is correctly configured you should see some log information.

Below is the output of me running `inventory send` on the same server where the Inventory service is deployed. Notice that there's an error in the output because I don't have Docker installed on that server.

```
root@beast:/home/bjk# inventory send
INFO  Config Used file=/etc/inventory/inventory.yml
INFO  Sending inventory
INFO  send remote=127.0.0.1:9999 location="Top of Rack" description="Dell 8950"
INFO  starting inventory client remote=127.0.0.1:9999
INFO  Getting host information
INFO  Host information hostname=beast IP=10.0.1.15
INFO  Getting network listeners
ERROR Docker not running, or can't connect error="Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
INFO  Found Incus container name=openwebui status=Running
INFO  Found Incus container name=stbeast status=Running
INFO  Inventory report sent successfully result=0
INFO  Inventory sent successfully
```

?> You'll notice in the output that there might be errors connecting to Incus or Docker.  If you don't have one of those running, that's to be expected. If you do have them running and you see an error, make sure your user belongs to the appropriate `docker` or `incus-admin` group.

Once you've sent your first report, refresh the page on your dashboard and check that the containers/instances you're running were reported.

### Scheduling Periodic Reports

Once you've validated the configuration of your client, you can use your favorite method to schedule reports to be sent to the server.  SystemD Timers and cron are two great options.

In the [contrib](https://github.com/bketelsen/inventory/blob/main/contrib/crontab) folder of the source code, there's an example crontab file.

```
*/2 * * * * /usr/local/bin/inventory send  >> /var/log/inventory.log 2>&1
```

This crontab entry runs `inventory send` every two minutes and sends the output to `/var/log/inventory.log`.

The reporting process is lightweight, so you can run it as frequently as you want without worrying about putting unnecessary load on your system.