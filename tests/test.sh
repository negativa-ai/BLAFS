#!/bin/bash

# Run this script under the root directory of the project

cd ./examples

python3 debloat.py profile ./redis /usr/bin/baffs
python3 debloat.py debloat ./redis /usr/bin/baffs
python3 debloat.py validate ./redis /usr/bin/baffs
