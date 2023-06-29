#!/bin/bash

arr1=($(/usr/local/bin/url_parser ${SOURCE_URL}))
arr2=($(/usr/local/bin/url_parser ${TARGET_URL}))

MODE / REMOTE_HOST / REMOTE_PORT  / REMOTE_ROOT / REMOTE_USER

if [[ (${arr1[0]} == "LOCAL") && (${arr2[0]} == "FTP") ]]; then
    if [[ $ACTION != '' ]]; then 
        action="$ACTION"
    else
        action="PUSH"
    fi
    local_root=${SOURCE_URL}
    remote_url=${TARGET_URL}
elif [[ (${arr2[0]} == "LOCAL") && (${arr1[0]} == "FTP")]]; then
    action="PULL"
    local_root=${TARGET_URL}
    remote_url=${SOURCE_URL}
else
    echo "Only one local and one remote allowed!" >&2
    exit
fi

cat > /env.sh << EOF
#!/bin/bash

export ACTION=$action
export LOCAL_ROOT=$local_root
export REMOTE_URL=$remote_url

EOF

env

chmod +x /env.sh


if [[ $RAM_DISK_GB != '' ]]; then 
    mount -t tmpfs -o size=${RAM_DISK_GB}G,mode=0777 tmpfs /work
    code=$?
    echo "create_tmpfs,ret_code="$code
fi
