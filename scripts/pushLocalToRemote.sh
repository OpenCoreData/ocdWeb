#!/bin/bash

# Push from local to remote...
# run from the scripts directory...
# dry run test:   --dry-run

LOCAL="/home/fils/src/go/src/opencoredata.org/ocdWeb/web/static/"
REMOTE="/mnt/dataVolumes/ocdweb/static/"

#scp -r $LOCAL  root@opencoredata.org:$REMOTE
rsync  -avzhe ssh $LOCAL root@opencoredata.org:$REMOTE
