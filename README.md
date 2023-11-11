# signal-cmd-executor

This is a small helper utility that can e.g. run alongside kubernetes pods to execute certain actions when a program issues a signal.

Possible use cases include invoking complex curl commands when a SIGUSR1 is received.

## Configuration

The program expects a config file in `/etc/config.json`, which may be overridden by setting the `CONFIG_FILE` env variable.
A sample config can be found in `config.json` in this repository.

All commands are passed to `/bin/bash -c "${command}"`, so all advanced bash features are available here.

The prebuilt image contains curl and wget and is alpine based, however, if you need extra software in your image, feel free to build your own image 
based off the one provided here.
