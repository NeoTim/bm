## Introduction

Bm is a simple tool for easy managing our bash scripts. We first add the script file path to bm use `add` command, 
and then we can run this script in anywhere.

## Installation
You can install bm with the following command:
```
go get github.com/zhangjikai/bm
```

## Usage
* `add,a [key] [bash script path]` - Add a bash/py/jar file path associated with the specified key to bm.
* `run,r [key] [params...]` - Run the target file associated with the specified key. If you want passing parameters to the target file, you should use "_" instead of "-"
* `delete,d [key]` - Delete the key from bm.
* `ls,l [prefix]` - List the keys that begin with prefix.
* `config,c  [type] [value]` - Set configurations of bm. The valid configuration is "DBPath".
* `push` - Call git push command based on the db storage directory.
    - In order to make `push` and `pull` work properly, you need to specify a git repository as the storage directory of templates. For example:
            ```
            cd <dir>
            git clone git@xxx.xxx/db.git
            tpl config DBPath <dir>/db
            ```
* `pull` - Call git pull command based on the db storage directory.
