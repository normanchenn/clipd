\launchctl unload ~/Library/LaunchAgents/clipd.plist
if [ -f /tmp/clipd.sock ]; then
	rm /tmp/clipd.sock
fi
\launchctl load -w ~/Library/LaunchAgents/clipd.plist
