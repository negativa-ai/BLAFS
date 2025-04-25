# Automating BAFFS

This directory provide examples of how to automate BAFFS using scripts.

To debloat an example, using the `debloat.py` script.
For instance, to debloat the `redis` container, you can run the following commands:

```bash
python3 debloat.py profile ./redis /usr/bin/baffs # run some profiling workloads
python3 debloat.py debloat ./redis /usr/bin/baffs # debloat based on the profiling workloads
python3 debloat.py validate ./redis /usr/bin/baffs # validate the debloated container
```


