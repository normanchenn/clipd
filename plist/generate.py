import os
import json
import pwd
from string import Template


class Config:
    def __init__(self, label="clipd.plist", executable_path="go/bin/daemon", output_path=".clipd/logs/output.log",
                 error_path=".clipd/logs/error.log", debug=True):
        user = pwd.getpwuid(os.getuid()).pw_name
        base_dir = os.path.expanduser(f"~{user}")
        self.label = label
        self.executable_path = os.path.join(base_dir, executable_path)
        self.output_path = os.path.join(base_dir, output_path)
        self.error_path = os.path.join(base_dir, error_path)
        self.debug = debug


def read_config(filename="config.json"):
    with open(filename, "r") as file:
        data = json.load(file)

    return Config(data.get("label"), data.get("executable_path"), data.get("output_path"), data.get("error_path"),
                  data.get("debug"))


def main():
    config = read_config()
    template = Template("""<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>${label}</string>
        <key>ProgramArguments</key>
        <array>
            <string>${executable_path}</string>
        </array>
        <key>KeepAlive</key>
        <true/>
        <key>StandardOutputPath</key>
        <string>${output_path}</string>
        <key>StandardErrorPath</key>
        <string>${error_path}</string>
        <key>Debug</key>
        <${debug}/>
    </dict>
</plist>""")
    plist = template.substitute(
        label=config.label,
        executable_path=config.executable_path,
        output_path=config.output_path,
        error_path=config.error_path,
        debug="true" if config.debug else "false"
    )
    with open(config.label, "w") as file:
        file.write(plist)


if __name__ == "__main__":
    main()
