# generate plist
cd plist
python3 generate.py
cd ..

# install server
cd daemon
go install
cd ..

# install client
cd clipd
go install
cd ..

# move plist
mkdir -p ~/Library/LaunchAgents
mv plist/clipd.plist ~/Library/LaunchAgents/clipd.plist

\launchctl load -w ~/Library/LaunchAgents/clipd.plist
