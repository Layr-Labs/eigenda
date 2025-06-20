#!/bin/bash
# poll the proxy endpoint until we get a 0 return code or 2mins have passed, in that case exit 1
timeout_time=$(($(date +%s) + 120))

while (( $(date +%s) <= timeout_time )); do
  if curl -X GET 'http://localhost:6666/health'; then
    exit 0
  else
    sleep 5
  fi
done

exit 1