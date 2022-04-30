# macaddrobfus
A tool to obfuscate MAC addresses in log files in a consistent manner.
This tool is in its early development stage, don't use it for important things.

- This tool aims to obfuscate a MAC address consistently throughout the same file: if `12:34:56:78:9a:bc` becomes `12:34:56:de:f0:12` in one place, it should be consistently so in other places
- This tool touches only the NIC part (the last 24 bits) and preserves the first OUI (vendor) part (the first 24 bits)

## Install
```
go install github.com/goura/macaddrobfus/cmd/macaddrobfus@latest
```

## Usage
macaddrobfus reads from stdin and writes to stdout, i.e.
```
cat log.json | macaddrobfus -secret "this is a secret" > output.json
```

If you don't supply a secret as an option, the secret will be randamly chosen when it runs. There is no way to know what random secret was chosen.
This means when you run it again, MAC addresses will be obfuscated in a different way.

If you want to keep your MAC address obfuscation consistent between multiple executions, or if you are working on multiple files that need to be consistently obfuscated, you should specify a secret.
