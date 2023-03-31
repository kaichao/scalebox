#!/bin/bash

arr1=($(/app/bin/url_parser ${SOURCE_URL}))
arr2=($(/app/bin/url_parser ${TARGET_URL}))

MODE / REMOTE_HOST / REMOTE_PORT  / REMOTE_ROOT / REMOTE_USER

if [[ (${arr1[0]} == "LOCAL") && ((${arr2[0]} == "SSH") || (${arr2[0]} == "RSYNC")) ]]; then
    action="PUSH"
    local_root=${arr1[1]}
    remote_mode=${arr2[0]}
    remote_host=${arr2[1]}
    remote_port=${arr2[2]}
    remote_root=${arr2[3]}
    remote_user=${arr2[4]}
elif [[ (${arr2[0]} == "LOCAL") && ((${arr1[0]} == "SSH") || (${arr1[0]} == "RSYNC"))]]; then
    action="PULL"
    local_root=${arr2[1]}
    remote_mode=${arr1[0]}
    remote_host=${arr1[1]}
    remote_port=${arr1[2]}
    remote_root=${arr1[3]}
    remote_user=${arr1[4]}
else
    echo "Only one local and one remote allowed!" >&2
    exit
fi

cat > /env.sh << EOF
#!/bin/bash

export ACTION=$action
export LOCAL_ROOT=$local_root
export REMOTE_MODE=$remote_mode
export REMOTE_HOST=$remote_host
export REMOTE_PORT=$remote_port
export REMOTE_ROOT=$remote_root
export REMOTE_USER=$remote_user

EOF

env

chmod +x /env.sh

if [[ ${remote_mode} ==  "SSH" ]]; then
    version=$(ssh -p $remote_port ${remote_user}@${remote_host} rsync -V|grep version|awk '{print $3}')
    # version="3.2.3"
    if $(dpkg --compare-versions ${version} "ge" "3.2.3"); then 
        touch /rsync_ver_ge_323
    fi
    ssh -p ${remote_port} ${remote_user}@${remote_host} mkdir -p ${remote_root}
fi

exit $?
