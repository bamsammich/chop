# chop
CLI tool to make structured logging human-readable.

## Usage

```text
Write structured log in a human-readable way.

Usage:
  chop [path] [flags]

Flags:
  -A, --after int                      print lines after this count
  -B, --before int                     print lines before this count [-1 unsets this] (default -1)
  -c, --columns strings                field names to extract to columns (default [message])
  -x, --exclude-extra-fields columns   exclude extra fields not defined by columns
  -h, --help                           help for chop
      --max-column-width int           set maximum column width (default 60)
  -m, --message-field string           field containing log message (default "message")
```

## How to use

`chop` takes logs that look like this:

```bash
> cat example/app.log
initializing app
{"timestamp":"2022/01/01 01:00:00", "message": "App has started", "level": "INFO"}
{"timestamp":"2022/01/01 01:00:01", "message": "Processing records", "level": "INFO"}
{"timestamp":"2022/01/01 01:03:00", "message": "Error processing records: timeout", "level": "ERROR"}
terminating app
```

And prints them like this:

```bash
> chop ./example/app.log
   #  message                   fields
   0  initializing app          map[]
   1  App has started           map[level:INFO timestamp:2022/01/01 01:00:00]
   2  Processing records        map[level:INFO timestamp:2022/01/01 01:00:01]
   3  Error processing records  map[exception:timeout level:ERROR timestamp:2022/01/01 01:03:00]
   4  terminating app           map[]
```

To add more/custom fields to `chop`'s output, pass additional fields as columns:

```bash
> chop ./example/app.log --columns timestamp,level,message
   #  timestamp            level  message                   fields
   0   ---                  ---   initializing app          map[]
   1  2022/01/01 01:00:00  INFO   App has started           map[]
   2  2022/01/01 01:00:01  INFO   Processing records        map[]
   3  2022/01/01 01:03:00  ERROR  Error processing records  map[exception:timeout]
   4   ---                  ---   terminating app           map[]
```

`chop` also accepts input from `stdin`:

```bash
> kubectl logs -n kube-system etcd-kind-control-plane | chop -m "msg" -c ts,level,caller,msg -B 5 -x
   #  ts                        level  caller                msg
   0  2022-11-23T16:34:44.841Z  info   etcdmain/etcd.go:73   Running:
   1  2022-11-23T16:34:44.849Z  info   etcdmain/etcd.go:116  server has been already initialized
   2  2022-11-23T16:34:44.849Z  info   embed/etcd.go:131     configuring peer listeners
   3  2022-11-23T16:34:44.851Z  info   embed/etcd.go:479     starting with peer TLS
   4  2022-11-23T16:34:44.858Z  info   embed/etcd.go:139     configuring client listeners
   5  2022-11-23T16:34:44.862Z  info   embed/etcd.go:308     starting an etcd server
```
