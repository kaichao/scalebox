# rsync-copy

## Introduction

## Environment variables

| parameter name   | Description  |
|  ----  | ----  |
| SOURCE_URL<br>TARGET_URL  | These two URL parameters represent the root directory on the remote ssh/rsync-server and the local root directory. The format of the remote URL is , and the format of the local URL is {local-path}. |

## Input Message

FILE_NAME: Relative file path to SOURCE_URL/TARGET_URL

## Error Code
| Code   | Description  |
|  ----  | ----  |
|  10  |  Connection timed out |
|  11  |  Input/output error |
|  23  |  Permission denied |
|  255  |   |

