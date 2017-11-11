#!/bin/bash

# query for parameters, 10 at a time
# convert the output to a cli-input-json
# run ssm delete-paramters

head='{ "Names": '
tail='}'
body=`aws ssm get-parameters-by-path --path /test --max-results 10 --recursive --query Parameters[].Name`
nextlot="$head$body$tail"

while [ `echo $nextlot | grep -c ,` != "0" ]; do
echo $nextlot

aws ssm delete-parameters --cli-input-json "$nextlot"
body=`aws ssm get-parameters-by-path --path /test --max-results 10 --recursive --query Parameters[].Name`
nextlot="$head$body$tail"
done


