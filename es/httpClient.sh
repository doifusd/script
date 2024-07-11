#!/bin/bash
logPath="/Users/sky/Documents/script/"
# strArray=("2023.05.16" "2023.05.17" "2023.05.18" "2023.05.19");
strArray=("2023.05.19");
for line in $(cat /Users/sky/Desktop/invoiceNo.txt)
do
    for value in ${strArray[@]};
    do
timeWhere=${value//./-}
urlStr="http://10.1.10.185:9200/logstash-prod-nginx-${value}/_doc/_search"
curl --request GET \
 --url $urlStr \
 --header 'Authorization: Basic bGVjb286dEw4aEhRNkRCOHBU' \
 --header 'content-type: application/json' \
 --data '{
 "_source": [
   "message"
 ],
 "query": {
   "bool": {
     "must": [],
     "filter": [
       {
         "bool": {
           "filter": [
             {
               "multi_match": {
                 "type": "phrase",
                 "query": "${line}",
                 "lenient": true
               }
             },
             {
               "multi_match": {
                 "type": "phrase",
                 "query": "/api/v1/invoice/notify",
                 "lenient": true
               }
             }
           ]
         }
       },
       {
         "range": {
           "time_local": {
             "format": "strict_date_optional_time",
             "gte": "2023-05-15T00:00:00.000Z",
             "lte": "2023-05-19T23:59:59.000Z"
           }
         }
       }
     ],
     "should": [],
     "must_not": [
       {
         "match_phrase": {
           "fields.tag": "ngw"
         }
       }
     ]
   }
 }
}'

    done;
done
