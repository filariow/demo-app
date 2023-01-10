#!/bin/env sh

[ -z $AWS_ACCESS_KEY_ID ] && echo "please define the AWS_ACCESS_KEY_ID env variable" && exit 1
[ -z $AWS_SECRET_ACCESS_KEY ] && echo "please define the AWS_SECRET_ACCESS_KEY env variable" && exit 1

ff=($(grep -Rl --include='*.yaml.tmpl' '$AWS_ACCESS_KEY_ID' ./config))

for f in ${ff[@]}; do
    t="${f%.*}"
    envsubst '$AWS_ACCESS_KEY_ID $AWS_SECRET_ACCESS_KEY' < "$f" > "$t" && \
        echo "generated $t from $f"
done

