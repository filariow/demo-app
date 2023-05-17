import json
import os
import boto3


TARGET_QUEUE = os.environ['TARGET_QUEUE']
sqs = boto3.resource('sqs')
queue = sqs.Queue(url=TARGET_QUEUE)

def handler(event, context):
    print(event)
    response = queue.send_message(MessageBody=json.dumps(event, indent=2))
    print(response.get('Failed'))
    return {
        'statusCode': 200,
        'headers': {
            'Content-Type': 'application/json'
        }
    }
