import configparser
import json
import boto3
import base64
import os
from botocore.exceptions import ClientError
from kubernetes import client, config

K8_HOST = os.environ["K8_HOST"]
K8_TOKEN_SECRET = os.environ["K8_TOKEN_SECRET"]

cf = boto3.client('cloudformation')
session = boto3.session.Session()
sm = session.client(service_name='secretsmanager')


def _get_stack_output(stack_name):
    result = {}
    response = cf.describe_stacks(StackName=stack_name)
    stack = response["Stacks"][0]
    outputs = stack['Outputs']
    for output in outputs:
        result[output["OutputKey"]] = output["OutputValue"]
    
    return result


def _get_secret(secret_id, secret_key):
    try:
        get_secret_value_response = sm.get_secret_value(
            SecretId=secret_id
        )
    except ClientError as e:
        raise e
    # Decrypts secret using the associated KMS key.
    secret = get_secret_value_response['SecretString']
    secret_o = json.loads(secret)
    return secret_o[secret_key]

def _get_k8_client_instance():
    configuration = client.Configuration() 
    configuration.host = K8_HOST
    configuration.verify_ssl=False

    token = _get_secret(K8_TOKEN_SECRET, "K8_TOKEN")

    configuration.api_key={"authorization": f"Bearer {token}"}
    client.Configuration.set_default(configuration)
    api_client = client.ApiClient()
    api_instance = client.CustomObjectsApi(api_client)

    return api_instance

def _create_primaza_registered_service(output, stack_name):
    api_instance = _get_k8_client_instance()
    core_api_instance = client.CoreV1Api()

    sed = []
    sci = []
    sci.append({"name": "type", "value": output.pop("type")})
    sci.append({"name": "provider", "value": "aws"})

    secret_data = {}
    for key, value in output.items():
        if "Ref" in key:
            secret_key = key.replace("Ref", "")
            secret_data[secret_key] = base64.b64encode(value.encode("utf-8")).decode("utf-8")
            sed.append({"name": key, "valueFromSecret": {"name": stack_name, "key": secret_key}})
        else:
            sed.append({"name": key, "value": value})

    metadata = {
        "name": stack_name,
        "namespace": "primaza-system"
    }

    spec = {
        "serviceClassIdentity": sci,
        "serviceEndpointDefinition": sed,
        "sla": "L3"
    }

    rs_body = { 
        "apiVersion": "primaza.io/v1alpha1", 
        "kind": "RegisteredService",
        "metadata": metadata,
        "spec": spec
    }


    secret = client.V1Secret()
    secret.api_version = "v1"
    secret.data = secret_data
    secret.kind = 'Secret'
    secret.metadata = {"name": stack_name, "namespace": "primaza-system"}
    secret.type = 'Opaque'

             
    api_instance.create_namespaced_custom_object(
        group="primaza.io",
        version="v1alpha1",
        namespace="primaza-system",
        plural="registeredservices",
        body=rs_body)
    
    core_api_instance.create_namespaced_secret("primaza-system", secret)

def _remove_primaza_registered_service(stack_name):
    api_instance = _get_k8_client_instance()
    core_api_instance = client.CoreV1Api()

    api_instance.delete_namespaced_custom_object(
        group="primaza.io",
        version="v1alpha1",
        namespace="primaza-system",
        plural="registeredservices",
        name=stack_name)
    core_api_instance.delete_namespaced_secret(stack_name, "primaza-system", body=client.V1Secret())

def _unpack_secrets(output):
    service_type = output["type"]
    if  service_type in ["aurora-postgresql", "aurora-mysql"]:
        secret_id = output["DBPasswordRef"]
        password = _get_secret(secret_id, "password")
        output["DBPasswordRef"] = password
    elif service_type == "dynamodb":
        access_key_ref = output["AccessKeyRef"]
        access_key = _get_secret(access_key_ref, access_key_ref)
        output["AccessKeyRef"] = access_key
        secret_key_ref = output["SecretKeyRef"]
        secret_key = _get_secret(secret_key_ref, secret_key_ref)
        output["SecretKeyRef"] = secret_key


def lambda_handler(event, context):
    message = event["Records"][0]["Sns"]["Message"]
    config = configparser.ConfigParser()
    config.read_string(f"[message]\n{message}")
    resource_type = config["message"]["ResourceType"].strip("'")
    status = config["message"]["ResourceStatus"].strip("'")
    if resource_type == "AWS::CloudFormation::Stack" and status == "CREATE_COMPLETE":
        stack_id = config["message"]["StackId"].strip("'")
        stack_name = config["message"]["StackName"].strip("'").lower()
        output = _get_stack_output(stack_id)
        _unpack_secrets(output)
        print(f"Registering New Service {stack_name}")
        _create_primaza_registered_service(output, stack_name)
        #print(json.dumps(output, indent=4))
    elif resource_type == "AWS::CloudFormation::Stack" and status == "DELETE_COMPLETE":
        stack_name = config["message"]["StackName"].strip("'").lower()
        print(f"Deleting Registered Service {stack_name}")
        _remove_primaza_registered_service(stack_name)

    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Service Catalag Discovery!')
    }
