# rtail
# send log to remote server

### data flow as follow
```
                              rtail(server) -> save_file
                             /
                         tcp:5254
                           /
watch_file -> rtail(client)

```

### init
```
git clone https://github.com/tongpengfei/rtail.git

go mod vendor
```

### build
```
make install
```

### usage
```
Usage:
  rtail [OPTIONS]

Application Options:
  -l, --listen_port= listen on port (default: 5254)
  -f, --save_file=   save to file (default: /var/log/rtail_log.txt)
  -H, --host=        connect to host eg. 127.0.0.1:5254
  -w, --watch_file=  watch file to tail (default: /var/log/client.log)

Help Options:
  -h, --help         Show this help message

watch_file => client => server => save_file
server: ./rtail -l 5254 -f /var/log/rtail/rtail_log.txt
client: ./rtail -H 127.0.0.1:5254 -w /var/log/client.log
```


### sample
```
#start server
$ ./rtail -l 5254 -f /var/log/all_client.log

#start client, watch /var/log/client.log and send it to server
$ ./rtail -H 127.0.0.1:5254 -w /var/log/client.log

#append log
echo hello >> /var/log/client.log
```