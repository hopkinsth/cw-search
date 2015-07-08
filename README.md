# cw-search
Pulls data from a specified set of log streams in cloudwatch logs

# Usage

To query CloudWatch Logs with this tool, you'll need to have...

1. AWS credentials set up either as environment variables (as in the example below) or in a [aws credential file](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-config-files)
2. A start and end time for your search formatted in YYYY-MM-DD HH:MM:SS format

## Example

```sh
$ export AWS_ACCESS_KEY_ID=accessKeyId
$ export AWS_SECRET_ACCESS_KEY=secretAccessKey
$ cw-search LogGroup1:LogStream1,[LogStreamN] -s "2015-04-07 00:00:00" -e "2015-04-07 23:59:59"
```

Can filter the fields in json:
```sh
$ cw-search --format json --fields name.first,address.country,address.state LogGroup:LogStream
{"name": {"first": "Thomas"}, "address": {"country": "United States", "state": "NC"}}
```

for other things, pipe to grep (for now):

```sh
cw-search LogGroup:LogStream | grep --color 'search str'
```
