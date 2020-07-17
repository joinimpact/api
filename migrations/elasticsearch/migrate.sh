#! /bin/bash

# Environment variables
HOSTNAME=${IMPACT_ELASTIC_HOST:-localhost}
PORT=${IMPACT_ELASTIC_PORT:-9200}

echo "Checking ./indices folder..."

# Iterate through all files in the indices folder.
for FILE in $(find ./indicies -name '*.json'); do
    BASE=$(basename "$FILE")
    INDEX_NAME="${BASE%%.*}"

    echo "Updating index $INDEX_NAME..."
    # Send request to create index if it doesn't already exist.
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT -H "Content-type: application/json" http://$HOSTNAME:$PORT/$INDEX_NAME)
    # On a 400 Bad Request response, the index already exists.
    if [ $RESPONSE == "400" ]
    then
        echo "Index $INDEX_NAME already exists, updating mapping..."
    fi

    # Send request to update mapping on the index.
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT -H "Content-type: application/json" -d "$(cat "$FILE")" http://$HOSTNAME:$PORT/$INDEX_NAME/_mapping)
    # On 400 Bad Request, there was an error updating the index mapping.
    if [ $RESPONSE == "400" ]
    then
        echo "Failed to map index $INDEX_NAME"
    fi
done