#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Please provide an action"
  exit 1
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$SCRIPT_DIR" || exit

ACTION=$1

# Curl variables
BASE_URL="http://localhost:8080/api/v1"
JSON_CONTENT_TYPE="Content-Type: application/json"

# JSON files
LOGIN_JSON_FILE="json/login.json"
USER_JSON_FILE="json/create_user.json"

# Token file
TOKEN_FILE="temp/tokens.json"


send_post_request() {
  local url=$1
  local json_file=$2
  local jq_filter=${3:-.}
  local json_data

  json_data=$(jq "$jq_filter" "$json_file") || exit 1
  
  curl -s -X POST "$url" -H "$JSON_CONTENT_TYPE" -d "$json_data" -w "%{http_code}\n" 
}

send_get_request() {
  local url=$1
  local bearer_token 

  if [ -f $TOKEN_FILE ]; then
    bearer_token=$(jq -r '.access_token' "$TOKEN_FILE")
  fi

  curl -s -X GET "$url" -H "$JSON_CONTENT_TYPE" -H "Authorization: Bearer $bearer_token" -w "%{http_code}\n" 
}

save_tokens() {
  echo "$1" | sed '$d' | jq '{access_token: .access_token, refresh_token: .refresh_token}' > "$TOKEN_FILE"
}

case $ACTION in
  login)
    response=$(send_post_request "$BASE_URL/login" "$LOGIN_JSON_FILE")
    save_tokens "$response"
    ;;
  create_user)
    response=$(send_post_request "$BASE_URL/users" "$USER_JSON_FILE")
    ;;
  renew_token)
    response=$(send_post_request "$BASE_URL/renew-token" "$TOKEN_FILE" '{refresh_token: .refresh_token}')
    save_tokens "$response"
    ;;
  get_user)
    response=$(send_get_request "$BASE_URL/user/1")
    ;;
  *)
    echo "Invalid action: $ACTION"
    exit 1
    ;;
esac

http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | sed '$d')

echo "$response_body" | jq
echo "Response Code: $http_code"