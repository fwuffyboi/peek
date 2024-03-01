# Peek
Take a peek at your server's usage statistics.


### What is Peek?
Peek is a simple and insanely fast utility tool designed for Linux servers that allows quick monitoring of server usage statistics.

This application is designed to be used on a server, however it may also be used on a personal computer.


### Features
- [x] Insanely fast
- [x] Simple to use
- [x] Easy to configure (It's just a yaml file)


### Installation
To install Peek, you will need to have Go installed on your system. If you do not have Go installed, you can download it from [here](https://golang.org/dl/) (https://golang.org/dl).

Once you have Go installed, you can install Peek by running the following commands in your directory of choosing:
```bash
git clone https://github.com/fwuffyboi/peek.git # Clone the repo
cd peek/src        # Go into the necessary directory
go build -o peek . # Build the file
sudo chmod +x peek # Make the file executable
sudo mv peek /usr/local/bin/peek # Move the file to /usr/local/bin
cd ../.. # Get out of the directory
sudo rm -rf peek # Delete the unnecessary repo
```

Then just run the command `peek` in your terminal to start. You can now access the server at its default port of `http://0.0.0.0:42649`.

⚠️⚠️⚠️⚠️⚠️

**WARNING: Currently, Peek does not have any authentication. This means that anyone on your server's local network can access the server's API, and they _WILL_ be able to access _ALL_ statistics and endpoints that are enabled in the configuration file. However, the default configuration is considered to be a "Safe default", allowing anyone on the local network to _ONLY_ view the logs of Peek (This is not sensitive information) or be able to see system information EXCEPT the server's public IP address. No actions (such as shutting the server down, or stopping peek) can be taken from the API on these defaults.** 

⚠️⚠️⚠️⚠️⚠️


### How to configure

To configure Peek, you must first run the application after moving it to /usr/local/bin/peek. This will create a default configuration in /home/{YOUR_USERNAME}/.config/peek called peek.config.yaml. It is recommended to stop the Peek application before editing this, as that can lead to unsaved changes. Once edited, start Peek again and it will load your new configuration. If there is an issue with it not working, please feel free to create a GitHub issue.

### Log file

The name of the log file in the Peek configuration file is what the log file will be called.
The log file's location will remain in `/home/YOUR_USERNAME/.config/peek/`. This only changes the file's name. Nothing else. The default log file value is "peek.log".

### Logging level

Peek allows you to choose what level of logging you would prefer.
The default is INFO, this shows most information that users would care about, and is very helpful if something goes wrong. This is the recommended option.

There are also the other options: WARN, ERR and FATA.

WARN only shows warnings in the program, and isn't very helpful. ERR only shows errors, and FATA only shows what caused a program to stop running. It is highly recommended to stick to the default. As this shows information that is critically helpful during debugging.

# How to uninstall

Run these commands after stopping the Peek application:
```bash
sudo rm -rf /home/{YOUR_USERNAME}/.config/peek
sudo rm -rf /usr/local/bin/peek
```


### Screenshots
![Screenshot](/src/assets/readme/ss-api-full.png)

The above screenshot shows the full API response from the server located at `/api/stats/all`.
This is the most detailed response you can get from the API. Below is shown an example of what peek would log for this request.
I request the /api/stats/all endpoint from my Pixel 6a device on IP 192.168.0.57, to my server at 192.168.0.80. The server's hostname is fedorable.

`{"clientIP":"192.168.0.57","dataLength":716,"hostname":"fedorable","latency":4,"level":"info","method":"GET","msg":"192.168.0.57 - fedorable [17/Feb/2024:10:22:08 +0000] \"GET /api/stats/all/\" 200 716 \"\" \"Mozilla/5.0 (Android 14; Mobile; rv:122.0) Gecko/122.0 Firefox/122.0\" (4ms)","path":"/api/full/","referer":"","statusCode":200,"time":"2024-02-17T10:22:08Z","userAgent":"Mozilla/5.0 (Android 14; Mobile; rv:122.0) Gecko/122.0 Firefox/122.0"}`

This request took 4ms(precisely 4.029865ms) total to complete. This includes getting all the data, reverse geolocating the IP, etc., and sending the response back to the client.

![Screenshot](/src/assets/readme/ss-api-endpoints.png)

The above screenshot shows the API endpoints available to the client. This is the response from the `/api/` endpoint.

![Screenshot](/src/assets/readme/ss-api-index.png)

The above screenshot shows the index page of the API. This is the response from the `/` endpoint.


### Credits/Contributors
- [fwuffyboi](https://github.com/fwuffyboi) - Creator, Documentation, API
- [db-ip.com](https://db-ip.com) - IP-to-country geolocation database


### Project TODO:
 - [ ] Add support for multiple servers to added to the web UI
 - [ ] Stick to a standard for logging errors, etc.
 - [ ] Create a WebUI
 - [ ] Support at least RU and EN languages
 - [ ] Add authentication to WebUI by default (NO DEFAULT PASSWORDS)


### TODO (Not in order of importance):
 - [x] Add screenshots to README
 - [ ] Add authentication on API - Note: This will be done after V1.0.0
 - [x] Add support for a yaml config
 - [x] Be able to get the server's country from the IP
 - [x] Allow viewing disk storage
 - [ ] Allow seeing RAID array information
 - [x] Ram usage
 - [x] CPU usage
 - [x] User friendly uptime
 - [x] User unfriendly uptime
 - [x] See hostname
 - [x] Be able to shut down server
 - [x] Custom rate limit to all API endpoints
 - [x] Be able to stop Peek - Note: Kinda works. It shuts down but does not tell the client. It just closes connection on them.
 - [x] CPU temperature
 - [ ] GPU temperature
 - [ ] GPU Usage
 - [ ] Auto-updating option
 - [ ] Allow selecting specific flag type (twitter, equal height, equal width, etc.)
 - [x] Make the API easier to parse
 - [ ] System time and timezone
 - [ ] Change where the config is stored
 - [ ] Add a fallback if ipinfo.io doesn't work
 - [ ] Create a better way to get the server's IP
 - [ ] See open sessions (ssh, etc.) and who/where they are from
 - [ ] Email if a certain stat gets to certain level
 - [ ] Be able to see if a certain port on the local network is responsive (ping)
 - [x] Be able to see the logs of Peek
 - [ ] Improve viewing the logs of Peek
 - [ ] Be able to see the config from api
 - [ ] Be able to see the logs of systemd processes
 - [x] Only ever return JSON from the whole API, no HTML
 - [ ] Be able to download the IP database if not present (from GitHub)
 - [ ] Be able to update the IP database if it is outdated
 - [ ] Be able to manually override the server's country


 - [ ] Release V1.0.0 once api is done (Without auth)


### Security
If security vulnerabilities are discovered, please see the Peek [security policy](https://github.com/fwuffyboi/peek/security/policy).
