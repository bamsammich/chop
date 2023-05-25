# chop
CLI tool to make structured logging human-readable.

## Install

```bash
go install github.com/bamsammich/chop@latest
```

## Usage

```text
Write structured logs in a human-readable way.

Usage:
  chop [path] [flags]

Flags:
  -f, --format strings   tuples of field names to print and column width (default [message=60])
  -h, --help             help for chop
  -a, --print-all        print all fields; fields without format defined will be printed as JSON
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
> chop ./example/app.log --format message=30 -a
                              message                                            _fields
     0               initializing app                                                  -
     1                App has started {"level":"INFO","timestamp":"2022/01/01 01:00:00"}
     2             Processing records {"level":"INFO","timestamp":"2022/01/01 01:00:01"}
     3       Error processing records {"exception":"timeout","level":"ERROR","timestamp":"2022/01/01 01:03:00"}
     4                terminating app                                                -
```

To add more/custom fields to `chop`'s output, pass additional fields as columns:

```bash
> chop ./example/app.log --format timestamp=19,level=5,message=50,_fields=25 -a
                 timestamp level                                            message                   _fields
     0                   -     -                                   initializing app                         -
     2 2022/01/01 01:00:01  INFO                                 Processing records                        {}
     3 2022/01/01 01:03:00 ERROR                           Error processing records   {"exception":"timeout"}
     4                   -     -                                    terminating app                         -
```

`chop` also accepts input from `stdin`:

```bash
> flog -f json -d 1s -n 10 | chop --format datetime=26,host=15,user-identifier=20,method=6,status=3,request=50
                         datetime            host      user-identifier method status                                            request
     0 24/May/2023:21:00:13 -0400 211.105.219.117                    -    GET    401                                           /dynamic
     1 24/May/2023:21:00:14 -0400     89.4.252.66           beatty7336  PATCH    302                 /communities/iterate/seamless/rich
     2 24/May/2023:21:00:15 -0400    91.45.224.98                    -  PATCH    405                         /dot-com/enable/synthesize
     3 24/May/2023:21:00:16 -0400  113.159.187.73                    -   POST    500                               /revolutionize/viral
     4 24/May/2023:21:00:17 -0400  149.165.249.24         mckenzie2588   POST    201                                   /next-generation
     5 24/May/2023:21:00:18 -0400  108.189.36.207                    -    GET    204                                        /one-to-one
     6 24/May/2023:21:00:19 -0400  219.22.174.217            bauch6671   HEAD    203                                     /best-of-breed
     7 24/May/2023:21:00:20 -0400  239.198.210.73            berge6168 DELETE    501                              /platforms/compelling
     8 24/May/2023:21:00:21 -0400 151.206.144.149                    -    GET    403     /web-readiness/front-end/one-to-one/productize
     9 24/May/2023:21:00:22 -0400   79.183.125.13          hilpert2613    PUT    406      /mission-critical/proactive/networks/magnetic
```
