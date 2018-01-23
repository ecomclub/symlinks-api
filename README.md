# symlinks-api
HTTP interface written in Go to create and delete symbolic links

# Technology stack
+ [Go](https://golang.org/) 1.9.x

# Setting up
For security, we recommend to download and install the app as root,
and let the files owned by `root:root` as default.

```bash
sudo git clone https://github.com/ecomclub/symlinks-api.git
cd symlinks-api
sudo go build main.go
```

Start application with CLI arguments:
+ Root directory to static files
+ HTTP/TCP port
+ Optional log file path

Example:

```bash
./main /var/www :3000 /var/log/app.log
```
