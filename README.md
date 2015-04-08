# cw-search
Searches a specified log stream in cloudwatch logs

# Usage
For now...

```sh
export AWS_ACCESS_KEY_ID=accessKeyId
export AWS_SECRET_ACCESS_KEY=secretAccessKey
cw-search --lg "LogGroupName" --ls "LogStreamName" -s "2015-04-07 00:00:00" -e "2015-04-07 23:59:59"
```

Need to filter? Until I implement something better, just pipe that fool to grep or similar:

```sh
cw-search --lg "LogGroupName" --ls "LogStreamName" | grep --color 'search str'
```
