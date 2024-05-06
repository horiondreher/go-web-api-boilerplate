#!/bin/bash

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
RENEW_TOKEN_JSON_FILE="json/renew_token.json"

send_post_request() {
  local url=$1
  local json_file=$2
  local json_data

  json_data=$(jq '.' "$json_file")
  
  curl -s -X POST "$url" -H "$JSON_CONTENT_TYPE" -d "$json_data" -w "%{http_code}\n" | {
    read -r body
    read -r code
    printf "Response code: %s\n\n" "$code"
    jq <<< "$body"
  }
}

case $ACTION in
  login)
    send_post_request "$BASE_URL/login" "$LOGIN_JSON_FILE"
    ;;
  create_user)
    send_post_request "$BASE_URL/users" "$USER_JSON_FILE"
    ;;
  renew_token)
    send_post_request "$BASE_URL/renew-token" "$RENEW_TOKEN_JSON_FILE"
    ;;
  *)
    echo "Invalid action: $ACTION"
    exit 1
    ;;
esac
