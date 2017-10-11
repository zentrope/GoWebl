#!/bin/sh
#
# PROVIDE: webl
# REQUIRE: DAEMON postgresql
# KEYWORD: shutdown
#

# https://www.freebsd.org/doc/en_US.ISO8859-1/articles/rc-scripting/index.html

. /etc/rc.subr

name=webl
rcvar=webl_enable

start_cmd="${name}_start"
stop_cmd="${name}_stop"

webl_dir="/usr/local/opt/webl"
webl_daemon="/usr/sbin/daemon"
webl_log="${webl_dir}/stdout.log"
webl_pid="${webl_dir}/webl.pid"
webl_conf="${webl_dir}/config.json"

command="${webl_dir}/webl"

load_rc_config $name
: ${webl_enable:=no}

webl_start() {
  echo "Hello" >> ${webl_log}
  (${webl_daemon} -p ${webl_pid} ${command} -c ${webl_conf}) >> ${webl_log} 2>&1
}

webl_stop() {
  /bin/kill `cat ${webl_pid}`
}

run_rc_command "$1"
