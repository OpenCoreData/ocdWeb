# Scripts

#  deployment scripts

* script for syncing to from local to remote data volume for static web

```
rsync --dry-run -avzhe ssh ~/Data/OCDDataVolumes/ocdweb root@opencoredata.org:/mnt/dataVolumes/ocdweb/
```

```
docker save earthcube/p418webui:latest | bzip2 | pv |  ssh -i /Users/dfils/.ssh/id_rsa root@149.165.157.173 'bunzip2 | docker load'

```
