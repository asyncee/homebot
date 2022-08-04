# Homebot

That's my home automation telegram bot that can download various content from the 
popular torrent tracker.

*Architecture*: clean architecture.

*Deploy*: via systemd because the bot is designed to be run single instance on a single 
machine and does not need to scale.

# Installation

As this is my personal bot so there are some hardcoded defaults, like pre-existing 
`smb` user.

```bash
wget https://github.com/asyncee/homebot/releases/download/latest/homebot-linux-amd64 -O /usr/bin/homebot

mkdir /etc/homebot
wget https://raw.githubusercontent.com/asyncee/homebot/main/env.example -O /etc/homebot/env
vim /etc/homebot/env
chmod 644 /etc/homebot/env

wget https://raw.githubusercontent.com/asyncee/homebot/main/homebot.service -O /etc/systemd/system/homebot.service
systemctl daemon-reload
systemctl enable homebot
systemctl start homebot
```
