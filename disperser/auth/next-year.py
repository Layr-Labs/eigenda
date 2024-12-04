#!/usr/bin/env python3

# This script prints the date one year from now. This can't be done in a shell script since the date
# command differs substantially between Linux and macOS.

import datetime

now = datetime.datetime.now()
nextYear = now.replace(year=now.year + 1)
print(nextYear.strftime('%Y-%m-%d'))
