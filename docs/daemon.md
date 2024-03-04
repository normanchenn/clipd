# scripts
restart_daemon.sh will restart the daemon, while setup_daemon.sh will create a new plist with your input configuration and restart the daemon
# starting daemon
```bash
launchctl load -w /path/to/clipd.plist
```

# stopping daemon
```bash
launchctl unload /path/to/clipd.plist
```
# recommended path
```bash
~/Library/LaunchAgents/
```
