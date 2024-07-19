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
  -d, --default-field string   default field for unstructured logs (default "message")
  -e, --exclude strings        fields to exclude
  -h, --help                   help for chop
  -i, --include strings        fields to print (excludes others)
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
❯ chop ./example/app.log
                                     message
---------------------------------------------
                            initializing app


                                     message    level               timestamp
------------------------------------------------------------------------------
                             App has started     INFO     2022/01/01 01:00:00
                          Processing records     INFO     2022/01/01 01:00:01


                                     message    level               timestamp   exception
------------------------------------------------------------------------------------------
                    Error processing records    ERROR     2022/01/01 01:03:00     timeout
                             terminating app    <nil>                   <nil>       <nil>
```

You can select which fields should be printed with `-i/--include:
```bash
❯ chop ./example/app.log -i message,level,exception
                                     message
---------------------------------------------
                            initializing app


                                     message    level
------------------------------------------------------
                             App has started     INFO
                          Processing records     INFO


                                     message    level   exception
------------------------------------------------------------------
                    Error processing records    ERROR     timeout
                             terminating app    <nil>       <nil>
```

Alternatively, you can exclude certain fields from printing with `-e/--exclude`:

```bash
❯ chop ./example/app.log -e level
                                     message
---------------------------------------------
                            initializing app


                                     message               timestamp
---------------------------------------------------------------------
                             App has started     2022/01/01 01:00:00
                          Processing records     2022/01/01 01:00:01


                                     message               timestamp   exception
---------------------------------------------------------------------------------
                    Error processing records     2022/01/01 01:03:00     timeout
                             terminating app                   <nil>       <nil>
```

`chop` also accepts input from `stdin`:

```bash
❯ ❯ flog -f json -d 100ms -n 10 | chop -i datetime,host,method,status,request


 status             host                       datetime                        request     method
--------------------------------------------------------------------------------------------------
    205     63.175.94.92     19/Jul/2024:09:54:10 -0400     /out-of-the-box/synthesize     DELETE


 status             host                       datetime                        request     method
--------------------------------------------------------------------------------------------------
    502   111.19.159.247     19/Jul/2024:09:54:11 -0400                      /wireless     DELETE


 status             host                       datetime                                       request     method
-----------------------------------------------------------------------------------------------------------------
    401     89.221.29.45     19/Jul/2024:09:54:11 -0400     /web+services/morph/cutting-edge/innovate       POST


 status             host                       datetime                                       request     method
-----------------------------------------------------------------------------------------------------------------
    302     37.19.69.141     19/Jul/2024:09:54:11 -0400                                  /open-source      PATCH
    301     190.75.31.62     19/Jul/2024:09:54:11 -0400               /reintermediate/deploy/redefine        PUT
    200     157.60.37.28     19/Jul/2024:09:54:11 -0400            /killer/architect/evolve/paradigms       HEAD
    100   80.155.247.118     19/Jul/2024:09:54:11 -0400                             /utilize/seamless        PUT
    400       6.9.54.144     19/Jul/2024:09:54:11 -0400                       /cross-media/e-business        PUT
    501     237.57.56.63     19/Jul/2024:09:54:11 -0400               /paradigms/e-business/platforms     DELETE
    200     134.227.7.62     19/Jul/2024:09:54:11 -0400                     /cultivate/maximize/morph        GET
```
