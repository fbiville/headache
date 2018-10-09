#!/bin/bash

dump_exclude() {
	exclusion=$2
	dump_file=$1

	sed -i -e "s&file:$2&xfile:$2&" "$1"
}

dump_exclude $1 $2
