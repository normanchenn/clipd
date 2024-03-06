# clipd
get macOS clipboard history from the command line
## setup
```bash
# to create a plist following the prompts
./setup_plist.sh

# to setup the daemon with the plist
./setup_daemon.sh

# to download the cli
./setup_cli.sh
```
## usage
```bash
# get the last element from clipboard
clipd get

# get the nth last element from clipboard 
clipd get --at=n

# get the last n elements from clipboard
clipd get --last=5

# get elements from n to m from clipboard
clipd get --from=n --to=m
```
