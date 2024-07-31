[client]
port = {port}

[mysql]
no-auto-rehash
default-character-set = utf8

[mysqld]
basedir = "{exePath}"
datadir = "{dataDir}"
port = {port}
server_id = 1
character-set-server = utf8

explicit_defaults_for_timestamp = on
skip-ssl

table_open_cache = 256

log_timestamps = SYSTEM

log-error = "{logDir}/mysql.log"
log_syslog = 0
[mysqldump]
quick
max_allowed_packet = 512M
