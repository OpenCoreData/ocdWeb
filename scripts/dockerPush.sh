#!/bin/bash

DOCKERVER="0.9.10"

docker save opencoredata/ocdweb:$DOCKERVER | bzip2 | pv |  ssh -i /home/fils/.ssh/id_rsa root@opencoredata.org 'bunzip2 | docker load'

