## ftp-copy

## Introductioin

ftp-copy is used to copy files from/to ftp-server, while lftp does the actual file transfer.

## Environment variables

| parameter name   | Description  |
|  ----  | ----  |
| SOURCE_URL<br>TARGET_URL  | These two URL parameters represent the root directory on the remote ftp-server and the local root directory. The format of the remote URL is ftp://{user}:{pass}@{host}:{port}/{path}, and the format of the local URL is {local-path}. |
| ACTION  | The action of file copying, its value includes PUSH'/'PULL'/'PUSH_RECHECK'. 'PUSH' and 'PULL' do not need to be set, determined by SOURCE_URL and TARGET_URL('PUSH': local->ftp server; 'PULL': ftp server->local). 'PUSH_RECHECK' means that after pushing to the ftp-server, re-pull to verify the consistency of the transmission. |
| NUM_PGET_CONN  | Maximum number of connections to get the specified file using several connections, default value is 4 |
| ENABLE_RECHECK_SIZE  | For action equal to 'PUSH' or 'PULL', recheck the file size is consistent after the file transfer. The default value is 'yes' |
| ENABLE_LOCAL_RELAY  | Enable local machine as transfer relay if local files are stored on network storage. If ACTION is set to 'PUSH_RECHECK', ENABLE_LOCAL_RELAY is always yes. The default value is 'no' |
| RAM_DISK_GB  | The ramdisk size of the local machine transfer relay in GB. If set, this value should be greater than 2 times the maximum file bytes. The default value is 0, no ramdisk cache |

## Input Message

FILE_NAME: Relative file path to SOURCE_URL/TARGET_URL
