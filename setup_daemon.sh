cd plist
go install
plist

mkdir ~/Library/LaunchAgents
cp clipd.plist ~/Library/LaunchAgents/clipd.plist

\launchctl unload ~/Library/LaunchAgents/clipd.plist
rm /tmp/clipd.sock
\launchctl load -w ~/Library/LaunchAgents/clipd.plist
