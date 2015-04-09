# cw-search
Searches a specified log stream in cloudwatch logs

# Usage
For now...

```sh
$ export AWS_ACCESS_KEY_ID=accessKeyId
$ export AWS_SECRET_ACCESS_KEY=secretAccessKey
$ cw-search LogGroup LogStream -s "2015-04-07 00:00:00" -e "2015-04-07 23:59:59"
```

Can filter the fields in json:
```sh
$ cw-search --format json --fields name.first,address.country,address.state LogGroup LogStream
{"name": {"first": "Thomas"}, "address": {"country": "United States", "state": "NC"}}
```

for other things, pipe to grep (for now):

```sh
cw-search LogGroup LogStream | grep --color 'search str'
```
