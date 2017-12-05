[mysqld]    
server-id=1         ##设置server-id
log_bin = mysql-bin ##开启binlog    
binlog_format=row   ##指定格式为row 
datadir=/var/lib/mysql    
socket=/var/lib/mysql/mysql.sock    
\# Disabling symbolic-links is recommended to prevent assorted security risks
symbolic-links=0    
\# Settings user and group are ignored when systemd is used.    
\# If you need to run mysqld under a different user or group,    
\# customize your systemd unit file for mariadb according to the    
\# instructions in http://fedoraproject.org/wiki/Systemd    
    
[mysqld_safe]    
log-error=/var/log/mariadb/mariadb.log    
pid-file=/var/run/mariadb/mariadb.pid    
    
\#    
\# include all files from the config directory    
\#    
!includedir /etc/my.cnf.d        