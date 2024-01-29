# data-grouping-2d

Grouping of data on a 2D plane

## dataset introduction

- x direction: horizontal direction
- y direction: vertical direction

## data-grouping

- Horizontal grouping: all data items are divided into 1 group;
- Vertical grouping: given the grouping length, group data according to this length

- command: ``` docker run -d -P {{ENVS}} {{VOLUMES}} {{IMAGE}} ```


## Environment variables

| variable name   | Description  |
|  ----  | ----  |
| DATASET_PREFIX | Optional parameter, usually set to project name |
| COORD_TYPE | 'string'/'integer', 'string' -> horizon-only |

## Input Message

## Output Message

## App Exit Code

