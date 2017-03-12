<?php
function getClientIP() 
{
    if (isset($_SERVER)) {
        if (isset($_SERVER['HTTP_X_FORWARDED_FOR'])) {
            $realip = $_SERVER['HTTP_X_FORWARDED_FOR'];
        } else if (isset($_SERVER['HTTP_CLIENT_IP'])) {
            $realip = $_SERVER['HTTP_CLIENT_IP'];
        } else {
            $realip = $_SERVER['REMOTE_ADDR'];
        }
    } else {
        if (getenv('HTTP_X_FORWARDED_FOR')) {
            $realip = getenv('HTTP_X_FORWARDED_FOR');
        } else if (getenv('HTTP_CLIENT_IP')) {
            $realip = getenv('HTTP_CLIENT_IP');
        } else {
            $realip = getenv('REMOTE_ADDR');
        }
    }

	if (strpos($realip,',') !== false) {
		$ips = explode(',',$realip);
        return end($ips);
	}

    return $realip;
}

function isEmail($email) 
{
    return preg_match('/^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/', $email);
}

/**
function vCode($num = 4, $size = 20, $width = 0, $height = 0) 
{
    if(empty($width)) {
        $width = $num * $size * 4 / 5 + 10;
    }

    if(empty($height)) {
        $height = $size + 10;
    }
    // 去掉了 0 1 O l 等
    $str = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVW";
    $code = '';
    for ($i = 0; $i < $num; $i++) {
        $code .= $str[mt_rand(0, strlen($str) - 1)];
    }
    // 画图像
    $im = imagecreatetruecolor($width, $height);
    // 定义要用到的颜色
    $back_color = imagecolorallocate($im, 235, 236, 237);
    $boer_color = imagecolorallocate($im, 118, 151, 199);
    $text_color = imagecolorallocate($im, mt_rand(0, 200), mt_rand(0, 120), mt_rand(0, 120));
    // 画背景
    imagefilledrectangle($im, 0, 0, $width, $height, $back_color);
    // 画边框
    imagerectangle($im, 0, 0, $width - 1, $height - 1, $boer_color);
    // 画干扰线
    for ($i = 0; $i < 5; $i++) {
        $font_color = imagecolorallocate($im, mt_rand(0, 255), mt_rand(0, 255), mt_rand(0, 255));
        imagearc($im, mt_rand(- $width, $width), mt_rand(- $height, $height), mt_rand(30, $width * 2), mt_rand(20, $height * 2), mt_rand(0, 360), mt_rand(0, 360), $font_color);
    }
    // 画干扰点
    for ($i = 0; $i < 50; $i++) {
        $font_color = imagecolorallocate($im, mt_rand(0, 255), mt_rand(0, 255), mt_rand(0, 255));
        imagesetpixel($im, mt_rand(0, $width), mt_rand(0, $height), $font_color);
    }
    // 画验证码
    @imagefttext($im, $size, 0, 5, $size + 3, $text_color, ROOT_PATH . '/web/static/font/tahoma.ttf', $code);
    $_SESSION["VerifyCode"] = $code;
    header("Cache-Control: max-age=1, s-maxage=1, no-cache, must-revalidate");
    header("Content-type: image/png;charset=gb2312");
    imagepng($im);
    imagedestroy($im);
}
**/

function httpPost($url, $post) {
    // Initialize a cURL session:
    $c = curl_init();
    // Set the URL that we are going to talk to:
    curl_setopt($c, CURLOPT_URL, $url);
    // Now tell cURL that we are doing a POST, and give it the data:
    curl_setopt($c, CURLOPT_POST, true);
    curl_setopt($c, CURLOPT_POSTFIELDS, $post);
    // Tell cURL to return the output of the page, instead of echo'ing it:
    curl_setopt ($c, CURLOPT_RETURNTRANSFER, true);
    // Now, execute the request, and return the data that we receive:
    $output = curl_exec($c);
	// Check if any error occured
	$info = curl_getinfo($c);
	
	if ($output === false || $info['http_code'] != 200) {
	  $output = "No cURL data returned for $url [{$info['http_code']}]";
	  if (curl_error($c))
		$output .= "\n". curl_error($c);
	}
	curl_close($c);
	return $output;
}

function getCookie($name, $default=null, $trim=true) 
{
    return getP($_COOKIE, $name, $default, $trim);
}

function getSession($name, $default=null, $trim=true) 
{
    return getP($_SESSION, $name, $default, $trim);
}

function getServer($name, $default=null, $trim=true) 
{
    return getP($_SERVER, $name, $default, $trim);
}

function getPost($name, $default=null, $trim=true) 
{
    return getP($_POST, $name, $default, $trim);
}

function getGet($name, $default=null, $trim=true) 
{
    return getP($_GET, $name, $default, $trim);
}

function getParam($name, $default=null, $trim=true) 
{
    return getP(getParams(), $name, $default, $trim);
}

function getParams() 
{
    return array_merge($_GET, $_POST);
}

function hasPost($name) 
{
    return isset($_POST[$name]);
}

function hasGet($name) 
{
    return isset($_GET[$name]);
}

function hasParam($name) 
{
    $params = getParams();
    return isset($params[$name]);
}

function getP($data, $name, $default=null, $trim=true) 
{
    if(isset($data[$name])) {
        return $trim ? trimParam($data[$name]) : $data[$name];
    } else {
        return $default;
    }
}

function trimParam($param) 
{
    if(!is_array($param)) {
        return trim($param);
    }

    foreach($param as $k=>&$v) {
       if(!is_array($v)) {
            $v = trim($v);
        } else {
            $v = trimParam($v);
        } 
    }
    return $param;
}
