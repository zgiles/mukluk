CREATE TABLE `nodes` (
  `uuid` char(36) NOT NULL,
  `hostname` varchar(100) NOT NULL,
  `ipv4address` varchar(15) NOT NULL,
  `macaddress` varchar(12) NOT NULL,
  `os_name` varchar(100) NOT NULL,
  `os_step` tinyint(4) NOT NULL,
  `init_data` blob NOT NULL,
  `node_type` varchar(100) NOT NULL,
  `oob_type` varchar(100) NOT NULL,
  `heartbeat` int(11) NOT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `nodes_discovered` (
  `uuid` char(36) NOT NULL,
  `ipv4address` varchar(15) NOT NULL,
  `macaddress` varchar(12) NOT NULL,
  `surpressed` tinyint(4) NOT NULL DEFAULT '0',
  `enrolled` tinyint(4) NOT NULL DEFAULT '0',
  `checkincount` int(11) NOT NULL DEFAULT '0',
  `heartbeat` int(11) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `os` (
  `os_name` varchar(100) NOT NULL,
  `os_step` tinyint(4) NOT NULL,
  `boot_mode` varchar(100) NOT NULL,
  `boot_kernel` varchar(255) DEFAULT NULL,
  `boot_initrd` varchar(255) DEFAULT NULL,
  `boot_options` varchar(255) DEFAULT NULL,
  `next_step` tinyint(4) NOT NULL DEFAULT '0',
  `init_data` blob,
  PRIMARY KEY (`os_name`,`os_step`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
