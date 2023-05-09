# list-dir

## Introduction

list-dir is a common module in scalebox. Its function is to traverse the file list of the local directory or the rsync/ftp remote directory, generate messages and send it to the subsequent module.

In this container image, curlftpfs is used to mount ftp-server as a local file system to manipulate metadata.

list-dir supports four types of directories:
- local: Local Directory
- rsync-over-ssh: Server directory backed by the ssh-based rsync protocol
- native rsync: Server directory that supports the standard rsync protocol
- ftp: Directory on ftp service that supports ssl encryption

## Environment variables

| parameter name   | Description  |
|  ----  | ----  |
| SOURCE_URL  | See table below. |
| REGEX_FILTER | File filtering rules represented by regular expressions |
| RSYNC_PASSWORD | Non-anonymous rsync user password |

### SOURCE_URL

| type | description |
| --- | ---- |
| local | represented by an absolute path ```</absolute-path> ```|
| rsync | anonymous access: ```rsync://<rsync-host><rsync-base-dir>```<br/> non-anonymous access: ```rsync://<rsync-user>@<rsync-host><rsync-base-dir>```|
| rsync-over-ssh | The ssh public key is stored in the ssh-server account to support password-free access <br/> ``` <ssh-user>@<ssh-host><ssh-base-dir>``` <br/>OR<br/> ``` <ssh-host><ssh-base-dir>```, default ssh-user is root |
| ftp | anonymous access: ```ftp://<ftp-host>/<ftp-base-dir>```<br/> non-anonymous access: ```ftp://<ftp-user>:<ftp-pass>@<ftp-host>/<ftp-base-dir>``` |

## Input Message

DIR_NAME: Relative to the subdirectory of SOURCE_URL, use "." to represent the current directory of SOURCE_URL

