package main

import (
	"os"
	"fmt"
    "net"
    "log"
    "errors"
	"strings"
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
)

var (
    kServerPort int = 5254
)

type TcpClient struct{
    conn net.Conn
}

func (o *TcpClient) Connect(host string) error{

    conn, err := net.Dial("tcp", host)
    if nil != err {
        return err
    }
    o.conn = conn
    return nil
}

func (o *TcpClient) Send(str string) error{
    if nil == o.conn {
        return errors.New("client conn is nil")
    }

    n, err := o.conn.Write([]byte(str))
	_ = n
    return err
}


type TcpServer struct{
    listener *net.Listener
}

func (o *TcpServer) Start(port int) error{

    listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
    if nil != err {
		fmt.Printf("listen port %d fail, %s\n", port, err)
        return err
    }
    o.listener = &listener

    fmt.Printf("listen on port %d\n", port)

    for{
        c, e := listener.Accept()
        if nil != e {
            fmt.Printf("listener accept error, %v\n", e)
            continue
        }

        fmt.Printf("connected %v\n", c.RemoteAddr().String())

        go func(conn net.Conn){
            defer conn.Close()
            
            buf := make([]byte, 4096, 4096)
            for{
                nrecv, erecv := conn.Read(buf[:])
                if nil != erecv {
                    fmt.Printf("disconnected %v\n", conn.RemoteAddr().String())
                    return
                }

                log.Printf("%v", string(buf[:nrecv]))
            }
        }(c)
    }

    return nil
}

func rtail_svr( port int, save_file string) error {
    dir := dir_name(save_file)
    mkdir_p( dir )

    //log
    f, err := os.OpenFile(save_file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("error opening file: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()

    log.SetOutput(f)

    log.SetPrefix("")
    log.SetFlags(0)

    svr := &TcpServer{}
    err = svr.Start(port)
    if nil != err {
        fmt.Printf("server start error, port:%d, %v", port, err)
        return err
    }
    return nil
}

func rtail_cli( src_file string, host string ) error {

	client := &TcpClient{}
	if err := client.Connect(host); nil != err {
		fmt.Printf("connect %s fail, %v\n", host, err)
		return err
	}

	t, err := tail.TailFile(src_file, tail.Config{Follow: true, ReOpen:true}) //-f -F
    if nil != err {
        return err
    }

	for line := range t.Lines {
	    fmt.Println(line.Text)
		client.Send(line.Text + "\n")
	}
    return nil
}

type Option struct {
    //server 
    ListenPort int    `short:"l" long:"listen_port" description:"listen on port" default:"5254"`
    SaveFile   string `short:"f" long:"save_file" description:"save to file" default:"/var/log/rtail_log.txt"`

    //client
    Host       string `short:"H" long:"host" description:"connect to host eg. 127.0.0.1:5254"`
    WatchFile  string `short:"w" long:"watch_file" description:"watch file to tail" default:"/var/log/client.log"`
}

var gopt = Option{}

func usage() {
    fmt.Printf("watch_file => client => server => save_file\n")
	fmt.Printf("server: ./rtail -l %d -f /var/log/rtail/rtail_log.txt\n", kServerPort)
	fmt.Printf("client: ./rtail -H 127.0.0.1:%d -w /var/log/client.log\n", kServerPort)
}

func parse_opt(){
	parser := flags.NewParser(&gopt, flags.Default)
    _, err := parser.Parse()
    if nil != err {
        //fmt.Println(err)
        if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
            usage()
            os.Exit(0)
        } else {
            usage()
            os.Exit(1)
        }
    }
}

func dir_name(path string) string {
    if strings.HasSuffix(path,"/") {
        dir := path[:len(path)-1]
        return dir
    }

    arr := strings.Split(path,"/")

    n := len(arr)
    last := arr[n-1]

    arr2 := arr
    if strings.Index(last, ".") != -1 {
        arr2 = arr[:n-1]
    }

    dir := strings.Join(arr2,"/")
    return dir
}

func mkdir_p( path string ){
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.MkdirAll( path, 0777 )
        os.Chmod( path, 0777 )
    }
}

func main(){
    parse_opt()

    if "" != gopt.Host {
        if "" == gopt.WatchFile {
            fmt.Println("need arg -w to watch file")
            os.Exit(1)
            return
        }
	    rtail_cli( gopt.WatchFile, gopt.Host )
    }else{
        rtail_svr( gopt.ListenPort, gopt.SaveFile)
    }
}