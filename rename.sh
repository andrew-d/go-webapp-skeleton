#!/bin/sh

set -e
set -u

main() {
    while true; do
        printf "Do you want to run the rename script (y/n)? "
        read answer

        if echo "$answer" | grep -iq "^y$"; then
            break
        elif echo "$answer" | grep -iq "^n$"; then
            echo "Cancelled"
            exit
        else
            echo "Please answer 'y' or 'n'"
        fi
    done

    printf "Please enter the name of the project: "
    read new_name
    printf "Please enter the new import path of the project: "
    read new_path

    echo "Renaming..."

    # Rename the project name and import path in the Makefile
    sed -i .bak \
        -e "s|github.com/andrew-d/go-webapp-skeleton|$new_path|g" \
        -e "s|NAME.*=.*skeleton|NAME := $new_name|g" \
        Makefile

    # Rename all the import paths in Go files
    find . -type f -name '*.go' \
        | grep -v '/vendor/' \
        | xargs sed -i .bak -e "s|github.com/andrew-d/go-webapp-skeleton|$new_path|g"

    # All done!
    echo "Done!"
}

main
