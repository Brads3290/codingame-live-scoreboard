#!/bin/zsh


#aws lambda update-function-code --function-name codingame-live-scoreboard-getevent --zip-file fileb://bin/api/getevent.zip

# Deploy
for d in bin/api/*.zip ; do
  d2="${d:gs/bin\/api\//}"
  d2="${d2:gs/.zip/}"

  fn_name="codingame-live-scoreboard-""$d2"

  echo "Deploying $fn_name"
  aws lambda update-function-code --function-name "$fn_name" --zip-file "fileb://""$d" > /dev/null
done
