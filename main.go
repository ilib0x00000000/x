package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	remoteProxy "github.com/ilib0x00000000/x/remote"
	"github.com/mholt/certmagic"
)

var (
	remote bool
	local  bool
	listen string
	domain string
	useTLS bool
)

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1") // tls v1.3

	flag.BoolVar(&remote, "remote", false, "start remote proxy")
	flag.BoolVar(&local, "local", false, "start local proxy")
	flag.StringVar(&listen, "listen", ":8888", "local proxy listen address")
	flag.StringVar(&domain, "domain", "", "remote proxy domain")
	flag.BoolVar(&useTLS, "tls", true, "enable tls for tunnel")
}

func main() {
	flag.Parse()

	if !remote && !local {
		usage()
		os.Exit(1)
	}

	if remote {
		startRemote()
	} else {
		startLocal()
	}

	return
}

// startRemote 远程代理启动入口
// 需要指定 域名  www.xxxx.com/www.xxxx.tk... 格式
// 默认开启 TLS
// 启动方式 go run main.go --remote --domain=xxxx --tls=true
func startRemote() {
	handler := remoteProxy.NewRemoteProxy()

	if useTLS {
		certmagic.Default.Email = "ilib0x00000001@gmail.com"
		certmagic.Default.CA = certmagic.LetsEncryptProductionCA

		ln, err := certmagic.Listen([]string{domain})
		if err != nil {
			panic(err)
		}

		err = http.Serve(ln, handler)
		if err != nil {
			panic(err)
		}
	}
}

// startLocal 本地代理启动入口
// 本地代理需要指定监听的地址和端口
// 启动方式 go run main.go --local --listen=:8800 --tls=true
func startLocal() {

}

// usage 启动代理方式
func usage() {
	fmt.Println("Usage: ")
	fmt.Println("\t in local")
	fmt.Println("\t\t go run main.go --local --listem=:8080 --tls=true")
	fmt.Println("\t\t set env http_proxy=http://127.0.0.1:8080 https_proxy=https://127.0.0.1:8080")
	fmt.Println("\t\t Then will make you feel happy cross The Great Fire Wall (゜-゜)つロ~")
	fmt.Println("\n")
	fmt.Println("\t in remote")
	fmt.Println("\t\t go run main.go --remote --domain=www.your-domain.com --tls=true")
	fmt.Println("\t\t I'm outside the wall")
	fmt.Println("\t\t" + `  _   _          _   _         _______ _       _ _______  _______  _______  _______  _______  _______  _______  _______ `)
	fmt.Println("\t\t" + ` (_) | |        (_) | |       / _____ \\\     /// _____ \/ _____ \/ _____ \/ _____ \/ _____ \/ _____ \/ _____ \/ _____ \`)
	fmt.Println("\t\t" + ` | | | |        | | | |       |/     \| \\   // |/     \||/     \||/     \||/     \||/     \||/     \||/     \||/     \|`)
	fmt.Println("\t\t" + ` | | | |        | | | |______ ||     ||  \\ //  ||     ||||     ||||     ||||     ||||     ||||     ||||     ||||     ||`)
	fmt.Println("\t\t" + ` | | | |        | | | |_____ \||     ||   ||    ||     ||||     ||||     ||||     ||||     ||||     ||||     ||||     ||`)
	fmt.Println("\t\t" + ` | | | |        | | | |     \|||     ||  //\\   ||     ||||     ||||     ||||     ||||     ||||     ||||     ||||     ||`)
	fmt.Println("\t\t" + ` | | | |_______ | | | |____ /||\_____/| //  \\  |\_____/||\_____/||\_____/||\_____/||\_____/||\_____/||\_____/||\_____/|`)
	fmt.Println("\t\t" + ` |_| |_________||_| |_|______/\_______///    \\ \_______/\_______/\_______/\_______/\_______/\_______/\_______/\_______/`)
}
