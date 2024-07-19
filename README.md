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


                                     message               timestamp    level
------------------------------------------------------------------------------
                             App has started     2022/01/01 01:00:00     INFO
                          Processing records     2022/01/01 01:00:01     INFO


                                     message               timestamp    level   exception
------------------------------------------------------------------------------------------
                    Error processing records     2022/01/01 01:03:00    ERROR     timeout
                             terminating app                   <nil>    <nil>       <nil>
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
❯ flog -f json -d 100ms -n 10 | chop -i datetime,host,method,status,request


                      datetime                        request  status              host  method
------------------------------------------------------------------------------------------------
    19/Jul/2024:09:38:02 -0400     /dot-com/synergistic/seize     406     25.42.250.210     PUT


                      datetime                        request  status              host  method
------------------------------------------------------------------------------------------------
    19/Jul/2024:09:38:02 -0400              /leverage/dynamic     502     165.86.242.53    HEAD


                      datetime                        request  status              host  method
------------------------------------------------------------------------------------------------
    19/Jul/2024:09:38:02 -0400     /niches/strategize/infrastructures/next-generation     304   191.161.227.211    POST


                      datetime                                                request  status              host  method
------------------------------------------------------------------------------------------------------------------------
    19/Jul/2024:09:38:02 -0400                    /value-added/web+services/sexy/grow     400       92.97.30.95  DELETE
    19/Jul/2024:09:38:02 -0400                                           /cross-media     406    60.105.221.131     PUT
    19/Jul/2024:09:38:02 -0400           /convergence/e-tailers/e-business/compelling     500       46.7.208.26   PATCH
    19/Jul/2024:09:38:02 -0400                               /interfaces/leading-edge     405      25.159.14.66     PUT
    19/Jul/2024:09:38:03 -0400                   /cutting-edge/world-class/extensible     501     174.218.96.64  DELETE
    19/Jul/2024:09:38:03 -0400                     /deploy/next-generation/synthesize     200    40.158.231.177     GET
    19/Jul/2024:09:38:03 -0400       /24%2f365/innovate/bricks-and-clicks/communities     501       4.82.26.248    POST
```
