#!/bin/bash

function parse_remote_path() {
    local s=$1
    if [[ $s =~ ^(/.*)$ ]]; then
        # /my-root
        echo "LOCAL $1"
    elif [[ $s =~ ^([^@]+@[^:/]+)(:[0-9]+)?(/.*)$ ]]; then
        # user@myhost:22/my-root
        # user@myhost/my-root
        ssh_host=${BASH_REMATCH[1]}
        ssh_port=${BASH_REMATCH[2]}
        data_root=${BASH_REMATCH[3]}
        # ssh_url="${ssh_host}:${data_root}"
        if [ "$ssh_port" == "" ]; then
            ssh_port="22"
        else
            ssh_port=${ssh_port:1}
        fi
        echo "SSH $data_root $ssh_host $ssh_port"
    else
        echo "wrong message format, message: $1" >&2
        return 26
    fi
}

# dir_path="/my-dir/my-subdir"
# parse_remote_path $dir_path
# dir_path="user@myhost/my-root"
# parse_remote_path $dir_path

function get_ssh_option() {
    local json="$1"
    local url_name="$2"
    local jump_servers_name="$3"

    local url=$(get_parameter "$json" "$url_name")
    read -r mode data_root ssh_host ssh_port < <(parse_remote_path $url)

    echo "url:$url" >&2
    echo "mode:$mode" >&2
    echo "ssh_host:$ssh_host" >&2
    echo "ssh_port:$ssh_port" >&2
    echo "data_root:$data_root" >&2
    
    jump_servers=$(get_parameter "$json" "$jump_servers_name") 
    jump_servers_option=""
    if [ $jump_servers ]; then
        jump_servers_option="-J '${jump_servers}' "
    fi
    ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x ${jump_servers_option}"
    option="-p ${ssh_port} ${jump_servers_option}"
    echo $option
}

# source "/usr/local/bin/functions.sh"
# headers='{
#   "source_jump_servers": "jump-servers:22",
#   "source_url": "user@myhost:10022/my-root"
# }'
# get_ssh_option "$headers" "source_url" "source_jump_servers"

function get_data_root() {
    local url="$1"

    data_root=$(parse_remote_path "$url"| cut -d' ' -f2)
    echo $data_root
}

# get_data_root "$headers" "source_url"

function get_mode() {
    local url="$1"
    mode=$(parse_remote_path "$url"| cut -d' ' -f1)
    echo $mode
}

function get_ssh_host() {
    local url="$1"
    ssh_host=$(parse_remote_path "$url"| cut -d' ' -f3)
    echo $ssh_host
}

function get_ssh_cmd(){
    local json="$1"
    local url_name="$2"
    local jump_servers_name="$3"

    ssh_option=$(get_ssh_option "$json" "$url_name" "$jump_servers_name")
    url=$(get_parameter "$json" "$url_name")
    ssh_host=$(get_ssh_host $url)
    echo "ssh $ssh_option $ssh_host"
}
