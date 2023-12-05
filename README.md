# Peek
Take a peek at your server's usage statistics.


### What is Peek?
Peek is a simple utility tool designed for servers that allows quick monitoring of server usage statistics.

It can also be used if you're curious about how much RAM or CPU usage a (linux or darwin derivative)device is using.

This application is designed to be used on a server, it may also be used on a personal computer, however this is not it's intended purpose.


### Auth Level
Peek has 3 levels of authentication, these are:
 - `0` - No authentication required for all pages OR API endpoints. (Highly not recommended)
 - `1` - Authentication required for all API endpoints that take action upon the server. (Reasonably secure, but not recommended as can leak sensitive information)
 - `2` - Authentication required, for ALL pages and API endpoints. (Most secure, best option for most use cases)


### Disclaimer
This application has not been tested thoroughly for exploits, thus, may not be a good idea to deploy on a server that is mission-critical or very sensitive. This application is secure for use on servers/equipment that is not of incredibly high value or highly targeted. 
If security vulnerabilities are discovered, please see our [security policy](https://github.com/fwuffyboi/peek/security/policy).


### Contributors:
- [fwuffyboi](https://github.com/fwuffyboi)


### Project TODO:
 - [ ] Add a screenshot to README

 - [ ] Add support for multiple servers to be linked (how??)
 - [ ] Add fallback if ipinfo.io doesnt work in getIP()
 - [ ] Create a better way to get the server's IP
 - [ ] Stick to a standard for logging errors, etc.
 - [ ] Add support for logs to be seen through WebUI
 - [ ] Add a functioning fucking WebUI
 - [ ] Support at least RU and EN
 - [ ] Add authentication to WebUI by default (NO DEFAULT PASSWORDS.)


### API TODO:
 - [ ] Add authentication on API
 - [ ] Add support for a yaml config
 - [x] Be able to get the server's country from the IP
 - [ ] Allow viewing disk storage use
 - [ ] Allow seeing RAID array information
 - [ ] Ram usage
 - [ ] CPU usage
 - [x] User friendly uptime
 - [x] User unfriendly uptime
 - [x] See hostname
 - [x] Be able to shut down server
 - [ ] Be able to stop Peek
 - [ ] CPU temperature
 - [ ] Email if a certain stat gets to certain level


 - [ ] Release V1.0.0 once api is done to basic level and supports yaml config.
