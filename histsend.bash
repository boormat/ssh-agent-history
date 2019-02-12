# intended to be sourced.

function __ha_to_ssh_int32() {
    # Converts a int32 to ssh encoded value
    # Returns a format string for printf output the binary.
    l=${1}
    printf '\\x%02X\\x%02x\\x%02X\\x%02X' $((0xFF & l>>24)) $((0xFF &  l>>16)) $(( 0xFF &  l>>8)) $((0xFF & l>>0))
}

function __ha_to_ssh_int8() {
    # Converts a int8 to ssh encoded value
    # Returns a format string for printf output the binary.
    l=${1}
    printf '\\x%02X' $((0xFF & l ))
}

function __ha_to_ssh_string() {
    # Converts a string to ssh encoded string.  (32bit length prefix)
    # take arg 1 STRING, and returns as an escaped ssh message.
    # that can in turn be passed to printf to output binary.
    declare -i l
    l=${#1}
    printf '%s%b' "$(to_ssh_int32 l)" "$1"
}

function __ha_hist_get() {
    # arg 1 is Comma
    # history -s
    local _cmd
    _cmd="$(history 1)"
    _cmd=${_cmd:7}

    # in message to ssh-agent we want...
    # command, hostname, User name
    # $HOSTNAME $UID are builtin.
    # $USER is not builtin, so avoid
    local msgtype="HISTORY"
    local payloadlen="$(( 4 + ${#_cmd} + 4 + ${#HOSTNAME} + 4 + ${#UID} ))" # 3 strings
    local totallen="$(( 1 + 4 + ${#msgtype} + 4 + payloadlen ))"

    printf %s%s%s%s%s%s%s \
        "$(__ha_to_ssh_int32 totallen )" \
        "$(__ha_to_ssh_int8 27 )" \
        "$(__ha_to_ssh_string "${msgtype}" )" \
        "$(__ha_to_ssh_int32 payloadlen )" \
        "$(__ha_to_ssh_string "${_cmd}" )" \
        "$(__ha_to_ssh_string "${HOSTNAME}" )" \
        "$(__ha_to_ssh_string "${UID}" )"
}


function history_trap() {
    local _cmd=$(history 1)
    local _cmdid=${_cmd:0:7}
    if [[ $_cmdid != $_lastCommand ]]
    then
        _lastCommand="$_cmdid"
        printf "$(hist_get)"    | nc -U  $SSH_AUTH_SOCK > /dev/null
    fi
}

trap history_trap DEBUG

