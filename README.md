# golangmutualssl

A simple test app showing how to do mutually authenticated SSL in go.

All the certificates are already in the appropriate directories.

The CA cert is `MyCA.crt' and is available to both the client and the server. 
The client and server certs are issued by that CA.

The client verifies that the server has a cert issued by the CA and that the
server common name is "myserver". Notice that it overrides the normal DNS-based 
checks.

The service verifies that the client has a cert issued by the  and that the
client common name is "myclient".
