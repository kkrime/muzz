import json
import time
import requests

# REST helpers
def post(host, port, url, data, headers):
    response = requests.post("http://"+ host + ':' + str(port) + url, data=data.encode('utf-8'), headers=headers)
    print(response)
    jsonResponse = response.json()
    # printJson(jsonResponse)
    return jsonResponse

def get(host, port, url, headers):
    response = requests.get("http://"+ host + ':' + str(port) + url, headers=headers)
    print(response)
    jsonResponse = response.json()
    # printJson(jsonResponse)
    return jsonResponse

# Debugging
def printJson(data):
    formattedJson = json.dumps(data, indent=2)
    print(highlight(formattedJson, JsonLexer(), TerminalFormatter()))

users = {}
coordiants = [{'long': 51.5072,'lat': 0.1276},{'long': 71.5072,'lat': 0.5276},{'long': 31.5072,'lat': 0.7276},{'long': 91.5072,'lat': 0.3276},{'long': 81.5072,'lat': 0.1276},]
# coordiants = [{'long':51.491972 ,'lat':-0.222220 },{'long':51.508784 ,'lat':-0.182220 },{'long':51.491972 ,'lat':-0.222220 },{'long':51.508784 ,'lat':-0.182220 },{'long':51.491972 ,'lat':-0.222220 },{'long':51.508784 ,'lat':-0.182220 },]

userCount = 5
for x in range(userCount):
    # usr = {}

    # create user
    res = get("localhost", 8080, "/user/create", None)
    print(res)
    name = res['result']['name'] 
    print("created user: ")
    print("name: " + res['result']['name'])
    print("email: " + res['result']['email'])
    print("password: " + res['result']['password'])
    print("\n")
    users[name] = res['result']

    # login
    login = {
        "email": res['result']['email'],
        "password": res['result']['password'],
        'long': coordiants[x]['long'],
        'lat': coordiants[x]['lat'],
    }
    # print(login)

    res = post("localhost", 8080, "/login", json.dumps(login, indent=2), None)
    print(res['token'])
    users[name]['token'] = res['token']

    # users.append(usr)
    # print(usr)

match = 0

while match != userCount:
    for name in users:
        u = users[name]
        # print(u)
        headers = {"Authorization": "Bearer " + u['token']}
        res = get("localhost", 8080, "/discover", headers)
        print(res)

        if res['result'] == None:
            match = match + 1
        else:
            for p in res['result']:
                swipeRight = False
                if int(time.time()) & 1:
                    swipeRight = True
                print(p)
                print(swipeRight)
                swipe = {"userID": p['id'], 'swipeRight':swipeRight}
                res = post("localhost", 8080, "/swipe", json.dumps(swipe, indent=2), headers)
                print(res)

