import sys
import base64
import struct

# get the second field from the public key file.
keydata = base64.b64decode(
  open('key.pub').read().split(None)[1])

parts = []
while keydata:
    # read the length of the data
    dlen = struct.unpack('>I', keydata[:4])[0]

    # read in <length> bytes
    data, keydata = keydata[4:dlen+4], keydata[4+dlen:]

    parts.append(data)
This leaves us with an array that, for an RSA key, will look like:

['ssh-rsa', '...bytes in exponent...', '...bytes in modulus...']
We need to convert the character buffers currently holding e (the exponent) and n (the modulus) into numeric types. There may be better ways to do this, but this works:

e_val = eval('0x' + ''.join(['%02X' % struct.unpack('B', x)[0] for x in
    parts[1]]))
n_val = eval('0x' + ''.join(['%02X' % struct.unpack('B', x)[0] for x in
    parts[2]]))
