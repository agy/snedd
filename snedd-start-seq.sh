#!/bin/bash

set -u
set -e

SNEDD_EXPIRER_FN_NAME=${1:-"snedd-expirer"}

function show_motd() {
	# motd is required to start with a newline
	echo
	figlet "self destruct sequence initiated"
	cat /run/snedd/triggered

	exit 0
}

function get_instance_id() {
	local instance_id=$(curl --fail --silent http://169.254.169.254/latest/meta-data/instance-id)
	echo ${instance_id}
}

function get_region() {
	local region=$(curl --fail --silent http://169.254.169.254/latest/meta-data/placement/availability-zone)
	# trim last character
	echo ${region:0:-1}
}

function get_doc() {
	local doc=$(curl --fail --silent http://169.254.169.254/latest/dynamic/instance-identity/pkcs7)
	echo ${doc}
}

function exec_lambda() {
	local region=${1}
	local instance_id=${2}
	local doc=${3}
	local fn_name=${4}

	local result=$(aws lambda invoke \
			--invocation-type RequestResponse \
			--function-name ${fn_name} \
			--region ${region} \
			--payload '{"instance-id":"'${instance_id}'", "pkcs7":"'${doc}'"}')
	echo ${result}
}

function main() {
	install -d -m 0700 -o root -g root /run/snedd
	[ -f /run/snedd/triggered ] && show_motd # already run

	local fn_name=${SNEDD_EXPIRER_FN_NAME}

	local instance_id=$(get_instance_id)
	local region=$(get_region)
	local doc=$(get_doc)
	local result=$(exec_lambda ${region} ${instance_id} ${doc} ${fn_name})

	echo ${result} > /run/snedd/triggered
	show_motd
}

main
