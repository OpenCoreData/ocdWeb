#!/bin/bash
DOCKERVER=$(<../VERSION)

docker save opencoredata/ocdweb:$DOCKERVER | bzip2 | pv |  ssh -i /home/fils/.ssh/id_rsa root@opencoredata.org 'bunzip2 | docker load'

