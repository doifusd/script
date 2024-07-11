<?php

$dataTime = [
    "2023.05.16",
	"2023.05.17",
	"2023.05.18",
	"2023.05.19",
	"2023.05.20"
];


$file_path = "./datasource/daichuli_6.txt";
if(file_exists($file_path)){
    $fp = fopen($file_path,"r");
    while(!feof($fp)){
        $tmp = fgets($fp);//逐行读取。如果fgets不写length参数，默认是读取1k。
        $arr = explode("\t",trim($tmp));
        foreach($dataTime as $value){
            getClient($value,$arr[0]);
        }
    }
}

function getClient($dataDate,$invoice_no){
    $dataDateTime = str_replace(".","-",$dataDate);
$curl = curl_init();

curl_setopt_array($curl, [
  CURLOPT_PORT => "9200",
  CURLOPT_URL => "http://10.1.10.185:9200/logstash-prod-nginx-".$dataDate."/_doc/_search",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "GET",
  //CURLOPT_POSTFIELDS => "{\n    \"_source\": [\n        \"request_body\"\n    ],\n    \"sort\": [\n        {\n            \"time_local\": {\n                \"order\": \"desc\",\n                \"unmapped_type\": \"boolean\"\n            }\n        }\n    ],\n    \"size\": 100,\n    \"from\": 1,\n    \"query\": {\n      \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"multi_match\": {\n            \"type\": \"phrase\",\n            \"query\": \"".$invoice_no."\",\n            \"lenient\": true\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"".$dataDateTime."T00:00:00.000Z\",\n              \"lte\": \"".$dataDateTime."T23:59:59.762Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": []\n    }\n    }\n}",
  //CURLOPT_POSTFIELDS => "{\n    \"_source\": [\n        \"messsage\"\n    ],\n    \"sort\": [\n        {\n            \"time_local\": {\n                \"order\": \"desc\",\n                \"unmapped_type\": \"boolean\"\n            }\n        }\n    ],\n    \"size\": 100,\n    \"from\": 1,\n    \"query\": {\n      \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"multi_match\": {\n            \"type\": \"phrase\",\n            \"query\": \"".$invoice_no."\",\n            \"lenient\": true\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"".$dataDateTime."T00:00:00.000Z\",\n              \"lte\": \"".$dataDateTime."T23:59:59.762Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": []\n    }\n    }\n}",
  CURLOPT_POSTFIELDS => "{\n  \"_source\": [\n    \"message\"\n  ],\n  \"sort\": [\n    {\n      \"time_local\": {\n        \"order\": \"desc\",\n        \"unmapped_type\": \"boolean\"\n      }\n    }\n  ],\n  \"size\": 100,\n  \"from\": 1,\n  \"query\": {\n    \"bool\": {\n      \"must\": [],\n      \"filter\": [\n        {\n          \"multi_match\": {\n            \"type\": \"phrase\",\n            \"query\": \"".$invoice_no."\",\n            \"lenient\": true\n          }\n        },\n        {\n          \"multi_match\": {\n            \"type\": \"phrase\",\n            \"query\": \"/api/v1/invoice/notify\",\n            \"lenient\": true\n          }\n        },\n        {\n          \"range\": {\n            \"time_local\": {\n              \"format\": \"strict_date_optional_time\",\n              \"gte\": \"".$dataDateTime."T00:00:00.000Z\",\n              \"lte\": \"".$dataDateTime."T23:59:59.762Z\"\n            }\n          }\n        }\n      ],\n      \"should\": [],\n      \"must_not\": []\n    }\n  }\n}",

  CURLOPT_HTTPHEADER => [
    "Authorization: Basic bGVjb286dEw4aEhRNkRCOHBU",
    "content-type: application/json"
  ],
]);

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if (!$err) {
        var_dump($response);
    $result = json_decode($response,true);
            
        if(!empty($result["hits"]) && count($result["hits"]) > 0){
            if(isset($result["hits"]["hits"][0]["_source"])){
            //echo $result["hits"]["hits"][0]["_source"]["request_body"].PHP_EOL;
                $message = $result["hits"]["hits"][0]["_source"]["message"];
                var_dump($message);
            }
        }
}

}
