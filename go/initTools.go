// Copyright 2012 Gary Burd
// reboot and shutdown
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

const DATETIME_TEMPLATE = "2006-01-02 15:04:05"

func body() string {
    return  `
<html>
<head>
<title>Hi</title>
<style type="text/css">
a {
    font-size:4em;                                                                                                                                                                            
    padding:5px;
    margin:5px;
    line-height:1.5em;
    font-family:"Courier New";
}

a:link { text-decoration: none;color: blue}
a:active { text-decoration:blink}
a:hover { text-decoration:underline;color: red} 
a:visited { text-decoration: none;color: green}

div#content {
    padding:10em 20em;
}
</style>
</head>
<body>
    <script type="text/javascript">
        function goto(action) {
            msg = "Are you sure " + action + "?"
            if(confirm(msg)) {
               location.href = "/" + action; 
            }
        }
    </script>
    <div id="content">
        <a href="javascript:void(0);" onclick="goto('shutdown');">shutdown</a>
        <a href="javascript:void(0);" onclick="goto('reboot');">reboot</a>
    </div>
</body>
</html>
`
}

func showMenu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start:\t", time.Now().Format(DATETIME_TEMPLATE))
	w.Write([]byte(body()))
	fmt.Println("end:\t", time.Now().Format(DATETIME_TEMPLATE))
}

func reboot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>rebooting... </h1>"))
	fmt.Println("reboot time:\t", time.Now().Format(DATETIME_TEMPLATE))
	cmd := exec.Command("/usr/bin/sudo", "/sbin/init", "6")
	cmd.Run()
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>shutdown .... </h1>"))
	fmt.Println("shutdown time:\t", time.Now().Format(DATETIME_TEMPLATE))
	cmd := exec.Command("/usr/bin/sudo", "/sbin/init", "0")
	cmd.Run()
}

func testExample() {
	fmt.Println("now:\t", time.Now().Format(DATETIME_TEMPLATE))
	cmd := exec.Command("/bin/ls", "-l", "/", "/home/neagle")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("failed.")
	}
	fmt.Println(out.String())
}

func main() {
	http.HandleFunc("/shutdown", shutdown)
	http.HandleFunc("/reboot", reboot)
	http.HandleFunc("/showMenu", showMenu)
	http.ListenAndServe(":8888", nil)
}
