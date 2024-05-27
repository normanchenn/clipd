# clipd
macOS clipboard history from the command line
## setup
```bash
./init.sh
# additional configuration can be done in plist/config.json

# to restart the service
./restart.sh
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

# help
clipd --help
```
