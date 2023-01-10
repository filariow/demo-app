# Install SQS Controller on OpenShift

1. Create the file `config/olm-aws/user-secret.yaml`
    ```
    cat << EOF > $AWS_SECRET_FILEPATH
    apiVersion: v1
    kind: Secret
    metadata:
      name: ack-user-secrets
      namespace: ack-system
    stringData:
      AWS_ACCESS_KEY_ID: $KEY_ID
      AWS_SECRET_ACCESS_KEY: $ACCESS_KEY
    EOF
    ```

2. Apply the configurations for the sqs controller
    ```
    kubectl apply -k config/olm-aws
    ```

3. Use OpenShift Console to install `AWS Controllers for Kubernetes - Amazon SQS`

4. Create the SQS Queue
    ```
    kubectl apply -k config/olm-aws-sqs
    ```

