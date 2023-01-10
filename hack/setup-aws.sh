#!/bin/env sh

# WARNING #
# This script has not been tested due to lack of permissions
#

USERNAME=demo-soa-ack-controller-dev
REGION=us-west-2
POLICIES=( arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess arn:aws:rds:$REGION:*:* )
AWS_SECRET_FILEPATH=config/olm-aws/user-secret.yaml
TARGET_NAMESPACE=ack-system

aws iam create-user --user-name "$USERNAME"
ACCESS_KEY_RES=$(aws iam create-access-key --user-name "$USERNAME")

for p in ${POLICIES[@]}; do
	aws iam attach-user-policy \
		--user-name "$USERNAME" \
		--policy-arn "arn:aws:iam::aws:policy/$p"
done

KEY_ID=$(echo $ACCESS_KEY_RES | jq -r '.AccessKey.AccessKeyId' )
ACCESS_KEY=$(echo $ACCESS_KEY_RES | jq -r '.AccessKey.SecretAccessKey' )

cat << EOF > $AWS_SECRET_FILEPATH
apiVersion: v1
kind: Secret
metadata:
  name: ack-user-secrets
  namespace: $TARGET_NAMESPACE
stringData:
  AWS_ACCESS_KEY_ID: $KEY_ID
  AWS_SECRET_ACCESS_KEY: $ACCESS_KEY
EOF

