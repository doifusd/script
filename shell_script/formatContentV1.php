#!/usr/local/bin/php -q
<?php
ini_set ('memory_limit',  '512M');
//$file_path = "/Users/sky/Documents/log/result/2023-08-07-111.log";
$file_path = $argv[1];
// $file_path = "/Users/sky/Documents/log/result/2023-08-07-111.txt";
if(file_exists($file_path)){
    $fp = fopen($file_path,"r");
    while(!feof($fp)){
        $tmp = fgets($fp);//逐行读取。如果fgets不写length参数，默认是读取1k。
        $startPot= strrpos($tmp,"notice_type")-3;

        $pot= strrpos($tmp,"}}} |");
        //$len = $pot-129+3;
        $len = $pot-$startPot+3;
        $subStr = substr($tmp,$startPot,$len);
        echo $subStr.PHP_EOL;
    }

}

