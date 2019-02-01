# intented to be sourced.

function to_ssh_int32() {
    # Converts a int32 to ssh encoded value
    # Returns a format string for printf output the binary.
    l=${1}
    printf '\\x%02X\\x%02x\\x%02X\\x%02X' $((0xFF & l>>24)) $((0xFF &  l>>16)) $(( 0xFF &  l>>8)) $((0xFF & l>>0))
}

function to_ssh_int8() {
    # Converts a int8 to ssh encoded value
    # Returns a format string for printf output the binary.
    l=${1}
    printf '\\x%02X' $((0xFF & l ))
}

function to_ssh_string() {
    # Converts a string to ssh encoded string.  (32bit length prefix)
    # take arg 1 STRING, and returns as an escaped ssh message.
    # that can in turn be passed to printf to output binary.
    declare -i l
    l=${#1}
    printf '%s%q' "$(to_ssh_int32 l)" "$1"
}

function hist_send() {
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

    printf "$(\
        to_ssh_int32 totallen ; \
        to_ssh_int8 27 ; \
        to_ssh_string "${msgtype}"; \
        to_ssh_int32 payloadlen ; \
        to_ssh_string "${_cmd}" ; \
        to_ssh_string "${HOSTNAME}"; \
        to_ssh_string "${UID}" )"  \
        | nc -U  $SSH_AUTH_SOCK | xxd
}

function history_trap() {
    local _cmd
    _cmd=$(history 1)
    _cmd=${_cmd:0:7}
    if [[ $_cmd != $_lastCommand ]]
    then
        #echo CHANGEd _cmd=$_cmd _lastCommand=$_lastCommand BASH_COMMAND="$BASH_COMMAND"
        _lastCommand="$_cmd"

        # TODO send history here

    fi
}

trap history_trap DEBUG

