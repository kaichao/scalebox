# dir-list

## Introduction

dir-list is a common module in scalebox. Its function is to traverse the file list of the local directory or the rsync/ftp remote directory, generate messages and send it to the subsequent module.

In this container image, curlftpfs is used to mount ftp-server as a local file system to manipulate metadata.

dir-list supports four types of directories:
- local: Local Directory
- rsync-over-ssh: Server directory backed by the ssh-based rsync protocol
- native rsync: Server directory that supports the standard rsync protocol
- ftp: Directory on ftp service that supports ssl encryption

## Environment variables

| variable name   | Description  |
|  ----  | ----  |
| SOURCE_URL  | See table below. |
| REGEX_FILTER | File filtering rules represented by regular expressions |
| RSYNC_PASSWORD | Non-anonymous rsync user password |
| REGEX_2D_DATASET | treat current dir as 2d dataset, and output its metadata; regex for extracting x/y in 2d dataset |
| INDEX_2D_DATASET | index for 2d dataset |


### SOURCE_URL

| type | description |
| --- | ---- |
| local | represented by an absolute path ```</absolute-path> ```|
| rsync | anonymous access: ```rsync://<rsync-host><rsync-base-dir>```<br/> non-anonymous access: ```rsync://<rsync-user>@<rsync-host><rsync-base-dir>```|
| rsync-over-ssh | The ssh public key is stored in the ssh-server account to support password-free access <br/> ``` [<ssh-user>@]<ssh-host>[:<ssh-port>][#<jump-servers>#]<ssh-base-dir>```, default ssh-user is root, default ssh-port is 22. The format of jump-servers is ```[<user1>@]<host1>[:<port1>],[<user>@]<host2>[:<port2>] ``` |
| ftp | anonymous access: ```ftp://<ftp-host>/<ftp-base-dir>```<br/> non-anonymous access: ```ftp://<ftp-user>:<ftp-pass>@<ftp-host>[:<ftp-port>]/<ftp-base-dir>``` |

## Input Message

DIR_NAME: 
- relative_path: Relative to the subdirectory of SOURCE_URL, use "." to represent the current directory of SOURCE_URL
- SOURCE_URL + relative_path: #-seperated 

## Output Message

## App Exit Code
