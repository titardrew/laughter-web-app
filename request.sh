echo "Expected user list (2 users)"
curl --location --request GET 'localhost:5000'

echo "\nExpected empty user"
curl --location --request POST 'localhost:5000/auth' \
    --data-raw '{}'

echo "\nExpected new user"
curl --location --request POST 'localhost:5000/auth' \
    --data-raw '{"DeviceId": "0x000AF"}'

echo "\nExpected just registered user but with Id 2 (not 223)"
curl --location --request POST 'localhost:5000/auth' \
    --data-raw '{"Id":"223", "DeviceId": "0x000AF"}'
echo "\nExpected rec picture"
curl --location --request POST 'localhost:5000/meme-rec' \
    --data-raw '{"UserId":"0"}'

echo "\nExpected rec picture and non-empty History in logs"
curl --location --request POST 'localhost:5000/meme-rec' \
    --data-raw '{"UserId":"2"}'

echo "\nExpected random rec picture and unk user in logs"
curl --location --request POST 'localhost:5000/meme-rec' \
    --data-raw '{"UserId":"11"}'

echo "\nExpected user list (3 users)"
curl --location --request GET 'localhost:5000'

echo "\nSave DB"
curl --location --request GET 'localhost:5000/save-user-db'
