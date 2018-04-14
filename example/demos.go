// timer
package main

import (
	"path"
	"crypto/sha512"
	"encoding/base64"
	"crypto/sha256"
	"crypto/sha1"
	"encoding/hex"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"fmt"
	"time"
	"syscall"
	"os/signal"
)

type A struct {
	x,y float64
}

type B struct {
	A
	z float64
}


func  readFile(name string) {

	file,err := os.Open(name)
	if err != nil {
		fmt.Println("open file errorcode", err)
		return
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		size, err := file.Read(buf)
		if size == 0 || err != nil {
			fmt.Println("read size is null")
			break
		}
		os.Stdout.Write(buf[:size])
	}
}

func writeFile() {
	return
	file,err := os.Create("test.txt")
	if err != nil {
		fmt.Println("open file err", err)
		return
	}
	defer file.Close()

	file.Write([]byte("this is test for go write."))

}

func keyWordFunc() {
	for i := 1; i < 10;i++ {
		fmt.Println("keyword:", i)
		if i > 5 {
			break
		}
	}
	s := 90
	switch s {
		case 80:
		fmt.Print("switch case:80")
		case 90:
		fmt.Print("switch case:90")
		default:
		fmt.Println("switch default")
	}
}

func printt(arr []int) {
	for _, a:= range arr {
		fmt.Println("print.", a)
	}
}

func get() {
	resp, err := http.Get("http://baidu.com")
	if err != nil {
		fmt.Println("http get err.", err)
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	os.Stdout.Write(body)
}

func main() {
	fmt.Println("Hello World!")
	time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")

	fmt.Println(os.Getenv("PATH"))
	var ar = [14]int {1,2,3,4,5,6,7,8,9,10}
	var a = ar[5:9]
	fmt.Println(len(a))

	fmt.Println(cap(a))
	m := map[string]float32 {"var1":1.1, "pi":3.1415}
	fmt.Println(m["var1"])
	fmt.Println(m["var1"])
	for i,d := range m {
		fmt.Println(i, d)
	}
	for key := range m {
		fmt.Println(key)
	}

	var b *B = &B{A{1.1,2.3}, 3.4}
	fmt.Println(b.A.x)
	fmt.Println(b.y)
	fmt.Println(b.z)
	readFile("timer.go")
	writeFile()
	keyWordFunc()
	printt([]int {1,2,3})
	//get()
	go signal_proc()
	md5Test("123")
	sha1Test("1234567687890")
	sha256Test("abc")
	sha512Test("abc")
	base64Test("admin:123")
	pathTest()
	//http.HandleFunc("/hi", handler)
	//http.ListenAndServe(":9090", nil)
}

type result struct {
	ErrNo int `json:"No"`
	ErrMsg string `json:"Msg"`
	Desc []string `json:"Desc"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte("this hello page."))
	res := result{
		ErrNo:200,
		ErrMsg: "OK",
		Desc: []string{"desc1", "desc3"},
	}
	resp, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(resp)
}

func signal_proc() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGALRM, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c

	//logger.Warn("Signal received: %v", sig)
	fmt.Println("singal is:", sig)

}

func md5Test(plain string) {
	h := md5.New()
	h.Write([]byte(plain))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}

func sha1Test(plain string) {
	h := sha1.New()
	h.Write([]byte(plain))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}

func sha256Test(plain string) {
	h := sha256.New()
	h.Write([]byte(plain))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}

func sha512Test(plain string) {
	h := sha512.New()
	h.Write([]byte(plain))
	fmt.Println(hex.EncodeToString([]byte(h.Sum(nil))))
}

func base64Test(plain string) {
	h := base64.StdEncoding
	encodeString := h.EncodeToString([]byte(plain))
	fmt.Println(encodeString)
	input := []byte("foo\x00bar")
	fmt.Println(base64.StdEncoding.EncodeToString(input))
}

func pathTest() {
	f := "d:/a/b/c.txt"
	fmt.Println(path.Base(f))
	fmt.Println(path.Match("^[d]t$",f))
}