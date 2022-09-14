# goblue
Midas blue file header rest api micro service

Sept/2022

Sets up a service that when passed a filename will return a json of the midas header and extended header

Command line:
  -p <port>                 - IPv4 port to listen on (default: 9580)
  -d <directory>            - goblue directory that files must be located in (default: current dir) 
                                (Only files in this directory will be used.  Subdirs are not followed for security reasons)

Endpoints:
  /version                   - return the build version of the service
  /api/v1/dir                - return a list of files in the goblue directory
  /api/v1/file/<filename>    - parse the file called filename and return the headers

NOTE:
  1) We plan to only parse blue files and only headers, not data
  2) headers that are of float type are written to json without a trailing .0 if they dont have a fraction part as per the json spec
  3) the adjunct header is returned as an array of bytes and not interpreted
  4) the returned json will have two elements,  "filename" and "header"

example:
    using the sample.tmp from blueheaders_test in /data
 
 bin/goblue -d /data &
 curl localhost:9580/api/v1/file/sample.tmp

{
    "filename":"sample.tmp",
    "header": {
        "version":"BLUE",
        "head_rep":"EEEI",
        "data_rep":"EEEI",
        "detached":0,
        "protected":0,
        "pipe":0,
        "ext_start":2,
        "ext_size":48,
        "data_start":512,
        "data_size":128,
        "type":1000,
        "flagmask":0,
        "timecode":0,
        "outmask":0,
        "pipeloc":0,
        "pipesize":0,
        "in_byte":0,
        "out_byte":0,
        "out_bytes":[0,0,0,0,0,0,0,0],
        "keylength":19,
        "adjunct":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,240,63,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,240,63,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
        "format":"SF",
        "inlet":0,
        "outlets":0,
        "keywords":{"IO":"X-Midas","VER":"1.1"},
        "extendedHeaders":{"KEYWORD1":"ONE","KEYWORD2":"TWO"}
    }
}

