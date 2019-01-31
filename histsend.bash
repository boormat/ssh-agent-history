# intented to be sourced.

function to_ssh_string() {
    # take arg 1...
    #return a 4 byte length then the string.
    # Hmmm longer than 256 :-(
    local a
    declare -i l
    l=${#1}
    echo $l
    a=${a}$(printf  '\x%h' $l)
    l=$(( l/256 ))
    a=${a}$(printf  '\x%h' $l)
    l=$((l/256))
    a=${a}$(printf  '\x%h' $l)
    l=$((l/256))
    a=${a}$(printf  '\x%h' $l)
    l=$((l/256))
    echo $a
    # $(printf  '\x0\x0\x0\x01\xB')
}
function hist_send() {
    # arg 1 is Comma
    # history -s
    local _cmd
    _cmd=$(history 1) l=${#1}
    _cmd=${_cmd:7}

    # in message to ssh-agent we want... 
    # command, hostname, User name
    $HOSTNAME $UID are builtin.
    # $USER is not builtin, lets avoid.
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

