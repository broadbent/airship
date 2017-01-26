#Assumes that you have the provsioner already running

import requests
import json

#Point this to the provioners northbound REST interface.
base_url = 'http://localhost:60000/'
json_header = {'content-type': 'application/json'}
def getTopology():
    url = base_url+'nodes'
    response = requests.get(url)
    return response

def getNode(node_id):
    url = base_url+'nodes/' + str(node_id)
    data = ''
    response = requests.get(url, data=data)
    return response

def provisionDocker(node_id, image_name):
    url = base_url+'nodes/'+str(node_id)+'/provision_docker'
    data = {'docker_name': image_name}
    response = requests.post(url, data=json.dumps(data), headers=json_header)
    return response

def provisionTOSCA(node_id, tosca):
    return 0

def removeService(service_id):
    url = base_url + 'remove_service'
    data = {'service_id': service_id}
    response = requests.delete(url, data=json.dumps(data), headers=json_header)
    return response



#Example commands to get started with

print(getTopology().json())

print(provisionDocker(1,"Netflix").json())

print(removeService(1).json())