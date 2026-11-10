#!/sbin/openrc-run

# THIS NEEDS TO GO TO /etc/init.d/putthenamehereplease
# name is extracted from there ,

name=$RC_SVCNAME
command="/usr/bin/agent"
pidfile="/run/$RC_SVCNAME.pid"
command_background="yes"

depend() {
    after net
    need localmount
}