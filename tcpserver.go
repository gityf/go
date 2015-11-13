// tcpserver
package tcpserver

import (
    "os"
    "net"
    "fmt"
)

func checkErr(err interface{}) {
    if err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}

func tcpServer(service string) {
    tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
    checkErr(err)
    listen, err := net.ListenTCP("tcp", tcpAddr)
    checkErr(err)
    for {
        conn,err := listen.Accept()
        if err != nil {
            continue
        }
        go handleConn(conn)
    }
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    cnt,err := conn.Read(buf)
    if err != nil {
        return
    }
    fmt.Println("read size:", cnt)
    conn.Write(buf) 
}

func main() {
    fmt.Println("Hello World!")
    tcpServer(os.Args[1])
}
