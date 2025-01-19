#!/bin/bash

# Версия, которую вы хотите установить
VERSION="$1"

MODULES=$(grep -oE "\./[^)]+" go.work)

for MODULE in $MODULES; do
    MODULE_NAME=$(echo $MODULE | sed 's|^\./||')
    TAG="$MODULE_NAME/$VERSION"

    echo "Creating tag $TAG"
    git tag "$TAG"
done