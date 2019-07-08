<?php
ini_set('display_errors', 1);
error_reporting(E_ALL);

const ERR_EXIT = 1;
const ERR_NON_EXIT = 0;

const JSON_DECODE_RESULT_TYPE_ARRAY = true;
const JSON_DECODE_RESULT_TYPE_OBJECT = false;

$sJsonFile = $argv[1] ?? '';


if (strlen($sJsonFile) == 0) {
	errMsg("Please input json file");
}

if (!file_exists($sJsonFile)) {
	errMsg("Not file json file: $sJsonFile");
}

if (!is_readable($sJsonFile)) {
	errMsg("Json file is not readable : $sJsonFile");
}

$sJson = trim(file_get_contents($sJsonFile));
if (empty($sJson)) {
	errMsg("Json file is empty !!!");
}

if ($argv[2] ?? 1) {
	$bJsonDecodeResultType = JSON_DECODE_RESULT_TYPE_ARRAY;
} else {
	$bJsonDecodeResultType = JSON_DECODE_RESULT_TYPE_OBJECT;
}

$mJson = json_decode($sJson, $bJsonDecodeResultType);

if (NULL === $mJson) {
	errMsg("Decode json error!");
}

//注: 一些老版本不支持 后面的参数列表
echo json_encode($mJson, JSON_PRETTY_PRINT|JSON_UNESCAPED_UNICODE|JSON_UNESCAPED_SLASHES);

function errMsg($sMsg, $iExit = ERR_EXIT) {
	echo $sMsg, "\n";	
	$iExit && exit();
}

