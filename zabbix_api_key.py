import requests
import json


payload = {
    'jsonrpc': '2.0',
    'method': 'user.login',
    'params': {
        'user': 'YOUR_ZABBIX_USER_NAME',
        'password': 'YOUR_ZABBIX_PASSWORD'
    },
    'id': 1,
    'auth': None
}

url = 'https://ZABBIX_URL/api_jsonrpc.php'
headers = {
    'content-type': 'application/json'
}

auth_token = requests.get(
    url, data=json.dumps(payload), headers=headers
).json().get('result')
