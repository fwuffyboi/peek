# Peek
Take a peek at your server's usage statistics.


### What is Peek?
Peek is a simple utility tool designed for servers that allows quick monitoring of server usage statistics.

This application is designed to be used on a server, it may also be used on a personal computer, however this is not it's intended purpose.


### Auth Level
Peek has 3 levels of authentication, these are:
 - `0` - No authentication required for all pages OR API endpoints. (Highly not recommended)
 - `1` - Authentication required for all API endpoints that take action upon the server. (Reasonably secure, but not recommended as can leak sensitive information)
 - `2` - Authentication required, for ALL pages and API endpoints. (Most secure, best option for most use cases)


### Disclaimer
This application has not been tested thoroughly for exploits, thus, may not be a good idea to deploy on a server that is mission-critical or very sensitive. This application is secure for use on servers/equipment that is not of incredibly high value or highly targeted. 
If security vulnerabilities are discovered, please see our [security policy](https://github.com/fwuffyboi/peek/security/policy).


### Contributors
- [fwuffyboi](https://github.com/fwuffyboi) - Creator, Documentation, API


### Project TODO:
 - [ ] Add a screenshot to README

 - [ ] Add support for multiple servers to added to the web UI
 - [ ] Add a fallback if ipinfo.io doesnt work
 - [ ] Create a better way to get the server's IP
 - [ ] Stick to a standard for logging errors, etc.
 - [ ] Add support for logs to be seen through WebUI
 - [ ] Create a WebUI
 - [ ] Support at least RU and EN languages
 - [ ] Add authentication to WebUI by default (NO DEFAULT PASSWORDS.)


### API TODO:
 - [ ] Add authentication on API
 - [x] Add support for a yaml config
 - [x] Be able to get the server's country from the IP
 - [ ] Allow viewing disk storage use
 - [ ] Allow seeing RAID array information
 - [ ] Ram usage
 - [ ] CPU usage
 - [x] User friendly uptime
 - [x] User unfriendly uptime
 - [x] See hostname
 - [x] Be able to shut down server
 - [x] Custom rate limit to all API endpoints
 - [x] Be able to stop Peek - Note: Kinda works. It shuts down but does not tell the client. It just closes connection on them.
 - [ ] CPU temperature
 - [ ] Email if a certain stat gets to certain level
 - [ ] Be able to see if a certain port is responsive
 - [ ] Be able to see the logs of Peek
 - [ ] Be able to download the IP database if not present


 - [ ] Release V1.0.0 once api is done(doesn't need auth) and supports yaml config.
