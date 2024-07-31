[client]
port = 3306

[mysql]
no-auto-rehash
default-character-set = utf8mb4

[mysqld]
basedir = "{exePath}"
datadir = "{dataDir}"
port = {port}
server_id = 1
character-set-server = utf8mb4

default_authentication_plugin = mysql_native_password

explicit_defaults_for_timestamp = on
tls-version = ''
skip-mysqlx 

table_open_cache = 256

log_timestamps = SYSTEM

log-error = "{logDir}/mysql.log"

[mysqldump]
quick
max_allowed_packet = 512M
