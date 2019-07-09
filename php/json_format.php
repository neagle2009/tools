<?php
ini_set('display_errors', 1);
error_reporting(E_ALL);

const ERR_EXIT = 1;
const ERR_NON_EXIT = 0;

$aOptions = [
	'f:',
	'a::',
	'h',
];
$aLongOptions = [
	'file:',
	'toarray::',
	'help',
];

$aInputParam = getopt(implode('', $aOptions), $aLongOptions);
if (empty($aInputParam) || key_exists('h', $aInputParam) || key_exists('help', $aInputParam)) {
	showHelp();
}


$sJsonFile = $aInputParam['f'] ?? ($aInputParam['file'] ?? '');

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

$bJsonDecodeResultTypeArray = false;
if ((isset($aInputParam['a']) && $aInputParam['a']) || (isset($aInputParam['toarray']) && $aInputParam['toarray'])) {
	$bJsonDecodeResultTypeArray = true;
}

$mJson = json_decode($sJson, $bJsonDecodeResultTypeArray);

if (NULL === $mJson) {
	errMsg("Decode json error!");
}

//注: 一些老版本不支持 后面的参数列表
echo json_encode($mJson, JSON_PRETTY_PRINT|JSON_UNESCAPED_UNICODE|JSON_UNESCAPED_SLASHES);

function errMsg($sMsg, $iExit = ERR_EXIT) {
	echo $sMsg, "\n";	
	$iExit && exit();
}

function showHelp() {
	$sFname = __FILE__;
	echo <<<HELP

	php $sFname -f jsonFile -a resultTypeArray

-f, --file			json file
-a, --toarray		out put json convert to array
-h, --help			show this message


HELP;
	exit;
}
