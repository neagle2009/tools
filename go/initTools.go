// Copyright 2017 guchuan
// reboot and shutdown

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

const (
	DATETIME_TEMPLATE = "2006-01-02 15:04:05"
	AUTH_KEY          = ""

	INIT_SHUTDOWN_CODE = "0"
	INIT_REBOOT_CODE   = "6"

	LISTEN_PORT = "7777"
)

func menuListHtml() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title> Hi </title>
<link rel="stylesheet" href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css">
<link rel="stylesheet" href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css">
<script src="//cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
<script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
<script type="text/javascript">
    function goto(action) {
        msg = "Are you sure " + action + "?"
            if(confirm(msg)) {
                location.href = "/" + action; 
            }
    }
</script>
</head>
<body style="padding: 50% 3em;">
    <p class="text-center">
        <button type="button" class="btn btn-primary btn-lg btn-block glyphicon glyphicon-off btn-danger" onclick="goto('shutdown');"> shutdown </button>
    </p>
    <p class="text-center"> </p>
    <p class="text-center">
        <button type="button" class="btn btn-primary btn-lg btn-block glyphicon glyphicon-refresh btn-success" onclick="goto('reboot');"> reboot </button>
    </p>
</body>
</html>
`
}

func showMenu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start:\t", time.Now().Format(DATETIME_TEMPLATE))
	w.Write([]byte(menuListHtml()))
	fmt.Println("end:\t", time.Now().Format(DATETIME_TEMPLATE))
}

func reboot(w http.ResponseWriter, r *http.Request) {
	if AUTH_KEY == "" || authCheck(w, r) {
		w.Write([]byte("<h1>rebooting... </h1>"))
		fmt.Println("reboot time:\t", time.Now().Format(DATETIME_TEMPLATE))
		go initRun(INIT_REBOOT_CODE)
	}
}

func initRun(c string) {
	//time.Sleep(1 * time.Second)
	cmd := exec.Command("/usr/bin/sudo", "/sbin/init", c)
	cmd.Run()
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	if AUTH_KEY == "" || authCheck(w, r) {
		w.Write([]byte("<h1>shutdown .... </h1>"))
		fmt.Println("shutdown time:\t", time.Now().Format(DATETIME_TEMPLATE))
		go initRun(INIT_SHUTDOWN_CODE)
	}
}

func testExample(w http.ResponseWriter, r *http.Request) {
	if AUTH_KEY == "" || authCheck(w, r) {
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
}

func authCheck(w http.ResponseWriter, r *http.Request) bool {
	pwd := r.PostFormValue("password")
	if pwd == "" {
		w.Write([]byte(inputPwdFormHtml("")))
		return false
	} else if pwd != AUTH_KEY {
		w.Write([]byte(inputPwdFormHtml("")))
		return false
	}
	return true
}

func inputPwdFormHtml(title string) string {
	if title == "" {
		title = ""
	}
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title> Please input auth code </title>
<link rel="stylesheet" href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css">
<link rel="stylesheet" href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css">
<script src="//cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
<script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
<script type="text/javascript">
function formCheck() {
    password = document.getElementById("password").value;
    if(password == "") {
        alert("please input password");
        return false;
    }
    return true;
}
</script>
</head>
<body style="padding: 50% 3em;">
        <form action="" method="POST">
        <div class="form-group">
            <input name="password" class="form-control" type="text" placeholder="auth code" id="password" />
        </div>
        <button type="submit" class="btn btn-default" onClick="return formCheck();" >Submit</button>
        </form>
</body>
</html>
`
}

func main() {
	http.HandleFunc("/shutdown", shutdown)
	http.HandleFunc("/reboot", reboot)
	http.HandleFunc("/showMenu", showMenu)
	http.HandleFunc("/test", testExample)
	http.ListenAndServe(":"+LISTEN_PORT, nil)
}
