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
    --data-raw '{"Id":"223", "DeviceId": "0x000AF", "history": ["999", "31"]}'

echo "\nExpected user list (3 users)"
curl --location --request GET 'localhost:5000'
