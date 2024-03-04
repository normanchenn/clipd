\launchctl unload ~/Library/LaunchAgents/clipd.plist
rm /tmp/clipd.sock
\launchctl load -w ~/Library/LaunchAgents/clipd.plist
