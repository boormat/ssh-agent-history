# ssh-agent-history
ssh agent with extra capabilities, history logging.

The goal is an extension to ssh-agent to fullfill a few things.

1. Centralised logging of all command history.  Add a small binary or script intended to 
be called on the CLIENT side bash command prompt to send a message over ssh-agent 
channel.
 The client can probably simply be a bash function to format the message.  (printf into the pipe?)

2. PKI Certificate managment.  Provide a PKCS#11 implementation that tunnels over ssh-agent.
   A shared library on target hosts, so that wget, curl, git, firefox etc can use your PKI 
   credentials.  Possibly pass through to a PKCS#11 thing such as GnomeKeyRing to manage all the 
   above.
## Build

Make em smaller.  4MB is OK.. but sheesh.
go build -ldflags="-s -w" 
Plus get UPX to pack them smaller.  (not in RHEL :-( )

Maybe the message code in go might be small fast enough to ignore bash?  (Especially since would still need
.so for PKSC#11 api)

 ## Bash
 
 Testing... with bash.
 Getting binary output needs a escape codes.  To do it dynamically, chain 2 printfs
   printf  '\x23\n' | xxd
   printf $(printf '\\x%d\\n' 23 ) | xxd

Talking to unix domain sockets seems to need netcat.
  printf  '\x0\x0\x0\x01\xB' | nc -U  $SSH_AUTH_SOCK
 
  printf  '\x0\x0\x0\x0A\x1B\x0\x0\x0\x05Hello' | nc -U  $SSH_AUTH_SOCK
  printf  '\x0\lx0\x0\x0E\x1B\x0\x0\x0\x05Hello\x0\x0\x0\x0' | nc -U  $SSH_AUTH_SOCK | xxd
  printf  '\x0\x0\x0\x0E\x1B\x0\x0\x0\x05Hello\x0\x0\x0\x0' | nc -U  $SSH_AUTH_SOCK | xxd

Message format is 4 bytes of message length, then 1 byte of message type, then per message.

E.g. list identities: Type  11 (0xB)
     agent extension: Type 27 (0x1B)  Then a string, then byte[]
     string is a 32bit length, then bytes. (ie same as message)
    byte[] in Extension message needs to be a string too (according to golang server)


Notes.   Ssh agent seems to close socket on response, so golang code setup to do the same.



## BASH Config

Hooking into bash history... 
There is trap DEBUG and PROMPT_COMMAND and PS0

PS0 does what we want, but is only on bash 4.4+.  :-(
PROMPT_COMMAND only fires when command finishes. Lame if it never does.
DEBUG trap fires for every command on the line, so you get duplicates.

BASH_COMMAND ... might be the same as history 1 maybe
The Debug trap thing would be bad if smashing a loop, so we will just detect a dupe to
prevent that.  (You can use A combo of DEBUG and PROMPT_COMMAND to supress the duplicates
but that is probably not required)


Just checking the history number We just use history to dedupe, 

    function PreCommand() {
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
    trap PreCommand DEBUG







### Notes


More complex way to simulate PS0, uses both DEBUG trap and PROMPT_COMMAND

        # This will run before any command is executed.
        function PreCommand() {
        if [ -z "$AT_PROMPT" ]; then
            return
        fi
        unset AT_PROMPT

        # Do stuff.
        echo "Running PreCommand"
        }
        trap "PreCommand" DEBUG

        # This will run after the execution of the previous full command line.  We don't
        # want it PostCommand to execute when first starting a bash session (i.e., at
        # the first prompt).
        FIRST_PROMPT=1
        function PostCommand() {
        AT_PROMPT=1

        if [ -n "$FIRST_PROMPT" ]; then
            unset FIRST_PROMPT
            return
        fi

        # Do stuff.
        echo "Running PostCommand"
        }
        PROMPT_COMMAND="PostCommand"

