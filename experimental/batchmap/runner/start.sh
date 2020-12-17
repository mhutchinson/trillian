#!/bin/bash
NETIP=$(ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
socat TCP4-LISTEN:8099,bind=${NETIP},reuseaddr,fork TCP4:localhost:8099 &

python -m apache_beam.runners.portability.local_job_service_main --port=8099
