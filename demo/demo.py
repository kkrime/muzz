import json
import time
import requests
import random


# REST helpers
def post(host, port, url, data, headers):
    response = requests.post("http://"+ host + ':' + str(port) + url, data=data.encode('utf-8'), headers=headers)
    # print(response)
    jsonResponse = response.json()
    # printJson(jsonResponse)
    return jsonResponse

def get(host, port, url, headers):
    response = requests.get("http://"+ host + ':' + str(port) + url, headers=headers)
    # print(response)
    jsonResponse = response.json()
    # printJson(jsonResponse)
    return jsonResponse

# Debugging
def printJson(data):
    formattedJson = json.dumps(data, indent=2)
    print(highlight(formattedJson, JsonLexer(), TerminalFormatter()))

users = {}
coordiants = [{'long': 51.5072,'lat': 0.1276},{'long': 71.5072,'lat': 0.5276},{'long': 31.5072,'lat': 0.7276},{'long': 91.5072,'lat': 0.3276},{'long': 81.5072,'lat': 0.1276},{'long': 51.5072,'lat': 0.1276},{'long': 71.5072,'lat': 0.5276},{'long': 31.5072,'lat': 0.7276},{'long': 91.5072,'lat': 0.3276},{'long': 81.5072,'lat': 0.1276},]

userCount = 10
for x in range(userCount):

    # create user
    res = get("0.0.0.0", 8080, "/user/create", None)
    print("user created\n")
    name = res['result']['name'] 
    print("name: " + name)
    print("email: " + res['result']['email'])
    print("gender: " + res['result']['gender'])
    print("password: " + res['result']['password'])
    users[name] = res['result']

    # login
    login = {
        "email": res['result']['email'],
        "password": res['result']['password'],
        'long': coordiants[x]['long'],
        'lat': coordiants[x]['lat'],
    }

    print("\n")
    print("loggin in ",name)
    res = post("0.0.0.0", 8080, "/login", json.dumps(login, indent=2), None)
    print(res['token'])
    users[name]['token'] = res['token']

    print("\n")
    # print("\n")

match = 0

while match != userCount:
    for name in users:
        u = users[name]
        # print(u)
        print(u['name'])
        headers = {"Authorization": "Bearer " + u['token']}
        res = get("0.0.0.0", 8080, "/discover", headers)
        # print(res)
        print("\n")

        if res['result'] == None:
            match = match + 1
        else:
            for p in res['result']:
                swipeRight =  random.choice([True, False] )

                swipe = {"userID": p['id'], 'swipeRight': swipeRight}
                res = post("0.0.0.0", 8080, "/swipe", json.dumps(swipe, indent=2), headers)

                if swipeRight:
                    print(u['name'] + " swiped right on " + p['name'])
                else:
                    print(u['name'] + " swiped left on " + p['name'])

                print('distance ' + p["distanceFromMe"])

                # print(swipeRight)
                if res['result']['matched']:
                    print("                                                    There's a match!")

            print("\n")
            print("\n")
            print("\n")
