example:
SETNX mssmhpc:mukluk:lock:12345 1
HMSET mssmhpc:mukluk:nodes:12345 hostname host1 ipv4address 1.2.3.4 macaddress a:b:c:d os_name redhat os_step 0 node_type physical oob_type ipmi heartbeat 0
SADD mssmhpc:mukluk:index:nodes:hostname:host1 12345
SADD mssmhpc:mukluk:index:nodes:ipv4address:1.2.3.4 12345
SADD mssmhpc:mukluk:index:nodes:macaddress:a:b:c:d 12345
SADD mssmhpc:mukluk:index:nodes:os_name:redhat 12345
SADD mssmhpc:mukluk:index:nodes:os_step:0 12345
SADD mssmhpc:mukluk:index:nodes:node_type:physical 12345
SADD mssmhpc:mukluk:index:nodes:oob_type:ipmi 12345
DEL mssmhpc:mukluk:lock:12345

SETNX mssmhpc:mukluk:lock:67890 1
HMSET mssmhpc:mukluk:nodes:67890 hostname host1 ipv4address 5.6.7.8 macaddress e:f:b:a os_name redhat os_step 0 node_type physical oob_type ipmi heartbeat 0
SADD mssmhpc:mukluk:index:nodes:hostname:host1 67890
SADD mssmhpc:mukluk:index:nodes:ipv4address:1.2.3.4 67890
SADD mssmhpc:mukluk:index:nodes:macaddress:a:b:c:d 67890
SADD mssmhpc:mukluk:index:nodes:os_name:redhat 67890
SADD mssmhpc:mukluk:index:nodes:os_step:0 67890
SADD mssmhpc:mukluk:index:nodes:node_type:physical 67890
SADD mssmhpc:mukluk:index:nodes:oob_type:ipmi 67890
DEL mssmhpc:mukluk:lock:67890


SETNX mssmhpc:mukluk:lock:35345 1
HMSET mssmhpc:mukluk:discoverednodes:35345 ipv4address 12.32.43.52 macaddress f:f:1:2 surpressed 0 enrolled 0 checkincount 1 heartbeat 0
SADD mssmhpc:mukluk:index:discoverednodes:ipv4address:12.32.43.52 35345
SADD mssmhpc:mukluk:index:discoverednodes:macaddress:f:f:1:2 35345
SADD mssmhpc:mukluk:index:discoverednodes:surpressed:0 35345
SADD mssmhpc:mukluk:index:discoverednodes:enrolled:0 35345
DEL mssmhpc:mukluk:lock:35345




SORT mssmhpc:mukluk:index:nodes:os_name:redhat BY nosort GET # GET mssmhpc:mukluk:nodes:*->hostname GET mssmhpc:mukluk:nodes:*->ipv4address GET mssmhpc:mukluk:nodes:*->macaddress GET mssmhpc:mukluk:nodes:*->os_name GET mssmhpc:mukluk:nodes:*->os_step GET mssmhpc:mukluk:nodes:*->node_type GET mssmhpc:mukluk:nodes:*->oob_type
