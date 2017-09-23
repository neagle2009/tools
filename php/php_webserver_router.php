<?php
/**
 * use for php web server
 * author: neagle2009@gmail.com
 * Usage: php -S 127.0.0.1:9001 router.php
 * 
 * */
ini_set("display_errors", 1);
error_reporting(E_ALL);

$aMimeTypes = [
    //'.css' => 'text/css',
    //'.js'  => 'application/javascript',
    //'.jpg' => 'image/jpg',
    //'.gif' => 'image/gif',
    //'.png' => 'image/png',
    //'.map' => 'application/json',
    //'.html' => 'text/html',
    '.asp' => 'text/html',
];

$sUri = $_SERVER['REQUEST_URI'];
if ($sUri == '/') {
    //按照默认方式处理
    return false;
}

//$sPreg /\.css|\.js|\.jpg|\.png|\.map|\.asp|\.html$/
$sPreg = '/'.implode('|', array_map(function($sExt){
    return '\\'.$sExt;
}, array_keys($aMimeTypes))).'$/';

if (!preg_match($sPreg, $sUri, $aMatch)) {
    //按照默认方式处理
    return false;
}

$sFile = $_SERVER['DOCUMENT_ROOT'].$sUri;
if (is_file($sFile)) {
    header("Content-Type: {$aMimeTypes[$aMatch[0]]}");
    readfile($sFile);
}
