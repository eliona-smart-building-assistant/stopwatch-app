go install golang.org/x/tools/cmd/goimports@latest

rm -rf ./apiserver/*

docker run --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli:v6.0.1 generate \
    -g go-server \
    --git-user-id eliona-smart-building-assistant \
    --git-repo-id python-eliona-api-client \
    -i /local/openapi.yaml \
    -o /local/apiserver \
    --additional-properties="packageName=apiserver,sourceFolder=,outputAsLibrary=true"

sudo chown $(id -u ${USER}):$(id -g ${USER}) -R apiserver
goimports -w ./apiserver

#rm -rf apiserver/.*
#rm -r apiserver/api
