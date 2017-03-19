# Dropbox Ignore

Automatically ignore directory names in Dropbox based on their name. The application will run as lightweight service
and will exclude files with selective sync on creation.

Installing node_modules won't cause a cpu stress test anymore!

* [Dropbox cli is required](https://www.dropbox.com/install-linux)
* The .dbignore does not support globbing or mask patterns. It is just a list with directory names. [example](./.dbignore)


## Installation
Build from source.

## Usage
Put the .dbignore file in the root of the directory that you want watched. All directories will be watched recursively.

You can run the application by one of the following:

1. Saving the binary in the same location as the .dbignore file
2. Save the binary anywhere you like and call it with the path to .dbignore file as second argument.

