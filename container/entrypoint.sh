#!/bin/bash
# Copyright © 2016 Samsung CNCT
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Determine uid gid of current directory
pwd
FILESYSTEM_UID=$(stat -c "%u" .)
FILESYSTEM_GID=$(stat -c "%g" .)
echo "Filesystem UID:GID is $FILESYSTEM_UID : $FILESYSTEM_GID"

# Allow UID and GID to be over ridden by environment variables
USERNAME=${LOCAL_USERNAME:-user}
UID=${LOCAL_UID:-$FILESYSTEM_UID}
GID=${LOCAL_GID:-$FILESYSTEM_GID}

# Add user with specified UID and GID
echo "Adding $USERNAME with UID:GID $UID:$GID"
useradd --shell /bin/bash -u $UID -o -c "" -m $USERNAME -g $GID
export HOME=/home/$USERNAME

exec /usr/local/bin/gosu $USERNAME "$@"
