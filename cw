#!/bin/bash

command="$@"
creds | fzf -m | while read cred; do
	user="$(echo $cred | cut -d: -f1)"
	pass="${cred#*:}"
	mod_command="$(echo $command | sed "s/USER/$user/g;s/PASS/$pass/g")"
	echo $mod_command
done
