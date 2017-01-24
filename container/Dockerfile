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

FROM golang:1.7.4

ENV DEBIAN_RELEASE jessie

ENV GLIDE_VERSION v0.12.3
ENV GLIDE_URL https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.tar.gz
ENV GLIDE_SHA256 0e2be5e863464610ebc420443ccfab15cdfdf1c4ab63b5eb25d1216900a75109

ENV GOSU_VERSION 1.9
ENV GOSU_URL https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-amd64
ENV GOSU_BINARY /usr/local/bin/gosu
ENV GOSU_SHA256 5ec5d23079e94aea5f7ed92ee8a1a34bbf64c2d4053dadf383992908a2f9dc8a

ENV LIBGIT2_VERSION 0.24.3
ENV LIBGIT2_URL https://github.com/libgit2/libgit2/archive/v$LIBGIT2_VERSION.tar.gz
ENV LIBGIT2_SHA256 0a24e6a51dbf3beecb0ebcd2cafb1e09b1212e910be6477b5de03c84a5586754

ENV LD_LIBRARY_PATH /usr/local/lib:$LD_LIBRARY_PATH

RUN apt-get update && apt-get install -y --no-install-recommends \
                cmake \
                apt-transport-https \
                ca-certificates \
                software-properties-common \
        && curl -fsSL https://yum.dockerproject.org/gpg | apt-key add - \
        && add-apt-repository \
                "deb https://apt.dockerproject.org/repo/ \
                 debian-$(lsb_release -cs) \
                 main" \
        && apt-get update && apt-get install -y --no-install-recommends \
                docker-engine \
        && rm -rf /var/lib/apt/lists/*


RUN curl -fsSL "$GLIDE_URL" -o glide.tgz \
        && echo "$GLIDE_SHA256 glide.tgz" | sha256sum -c - \
        && tar xzvf glide.tgz \
        && mv linux-amd64/glide /usr/local/bin \
        && rm -rf linux-amd64 \
        && rm glide.tgz

RUN curl -fsSL "$GOSU_URL" -o "$GOSU_BINARY" \
        && echo "$GOSU_SHA256 $GOSU_BINARY" | sha256sum -c - \
        && chmod +x "$GOSU_BINARY"

RUN curl -fsSL "$LIBGIT2_URL" -o libgit2.tgz \
        && echo "$LIBGIT2_SHA256" libgit2.tgz | sha256sum -c - \
        && tar xfz libgit2.tgz \
        && cd libgit2-$LIBGIT2_VERSION \
        && mkdir build \
        && cd build \
        && cmake .. \
        && cmake --build . \
        && cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local \
        && cmake --build . --target install \
        && cd ../../ \
        && rm -rf libgit2-$LIBGIT2_VERSION \
        && rm libgit2.tgz

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 755 /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
