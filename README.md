# mukluk

An easy and modern OS deployment tool for clusters.

## Why / Goal
Today's network booting tools leave much to be desired. Mukluk is an attempt at creating a modern network-booting server for cluster deployment and management without legacy burdens.

## Features
* No external dependencies at runtime (except the config file / env vars)
* Boot any network-bootable OS, including Linux installers, thin-nodes, VMware, memtest, WinPE, etc
* Small enough project that it can be understood
* API for external inspection and node self-discovery

## Installation / Building
Pre-built binary versions can be found under the Releases tag at the top

If you would like to build it yourself:  
```
go get github.com/zgiles/mukluk
```

## Database
A database will need to be configured.  
If using MySQL, create the tables from docs/mysql/schema.sql  
If using Redis, no work needs to be done, mukluk will create keys on the fly and prefix them with "mukluk:"

## Usage / Config
mukluk currently pulls its config from the file config.toml, in the same folder as the binary. Fill in the config details. The sample config file is in the `docs/` folder  
Then run `./mukluk`  
Logging is to STDERR and pretty verbose.  
Optionally make some type of startup script. ( or wait until I do it.. see TODOs)  

## KNOWN ISSUES
* Redis plugin may be a bit behind.. mysql should be the only properly working one for now
* Checking of values in the URL is not fully implemented. It could be possible to get a bad field and database error (but it will probably be handled well). It is not recommended to make a mukluk server publicly accessible. It is not very secure.

## TODOs
* API docs
* 12 factor style config, etc, allow config in different location
* startup scripts
* tests
* update redis functionality
* detect schema on database server and update / change it
* flatfile database
* multiple databases / multiple readers writers
* benchmarking
* Debug GUI etc
* CLI tool
* nodediscovery to node conversion
* godocs
* variable verbosity / logging options
* better abstraction in the handlers.go file
* TLS
* dependency versioning

## Release History
* 0.1.1 - make action on failure an option through config
* 0.1.0 - Initial release

## License
Copyright (c) 2016 Zachary Giles  
Licensed under the MIT license.  
