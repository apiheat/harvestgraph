FROM amd64/alpine:latest

ARG harvestgraph_version

ENV AKAMAI_CLI_HOME=/cli

RUN apk add --no-cache --update openssl \
        ca-certificates \
        libc6-compat \
        libstdc++ \
        wget \
        jq \
        bash \
        nodejs \
        npm \
        graphviz \
        font-bitstream-type1 && \
        rm -rf /var/cache/apk/* && \
    wget --quiet -O /usr/local/bin/akamai https://github.com/akamai/cli/releases/download/1.1.4/akamai-1.1.4-linuxamd64 && \
    chmod +x /usr/local/bin/akamai && \
    echo '[ ! -z "$TERM" -a -r /etc/motd ] && cat /etc/motd' >> /root/.bashrc


RUN mkdir -p /cli/.akamai-cli && \
    echo "[cli]" > /cli/.akamai-cli/config && \
    echo "cache-path            = /cli/.akamai-cli/cache" >> /cli/.akamai-cli/config && \
    echo "config-version        = 1" >> /cli/.akamai-cli/config && \
    echo "enable-cli-statistics = false" >> /cli/.akamai-cli/config && \
    echo "last-ping             = 2018-04-27T18:16:12Z" >> /cli/.akamai-cli/config && \
    echo "client-id             =" >> /cli/.akamai-cli/config && \
    echo "install-in-path       =" >> /cli/.akamai-cli/config && \
    echo "last-upgrade-check    = ignore" >> /cli/.akamai-cli/config

RUN akamai install https://github.com/apiheat/akamai-cli-netlist --force && \
    rm -rf /cli/.akamai-cli/src/akamai-cli-netlist/.git
RUN akamai install appsec --force && \
    rm -rf /cli/.akamai-cli/src/cli-appsec/.git
RUN wget --quiet -O /usr/local/bin/harvestgraph https://github.com/apiheat/harvestgraph/releases/download/v$harvestgraph_version/harvestgraph_linux_amd64 && \
    chmod +x /usr/local/bin/harvestgraph

RUN echo '                                                                        ' >  /etc/motd && \
    echo '   _                                                            _       ' >> /etc/motd && \
    echo '  | |                                _                         | |      ' >> /etc/motd && \
    echo '  | | _   ____  ____ _   _ ____  ___| |_  ____  ____ ____ ____ | | _    ' >> /etc/motd && \
    echo '  | || \ / _  |/ ___) | | / _  )/___)  _)/ _  |/ ___) _  |  _ \| || \   ' >> /etc/motd && \
    echo '  | | | ( ( | | |    \ V ( (/ /|___ | |_( ( | | |  ( ( | | | | | | | |  ' >> /etc/motd && \
    echo '  |_| |_|\_||_|_|     \_/ \____|___/ \___)_|| |_|   \_||_| ||_/|_| |_|  ' >> /etc/motd && \
    echo '                                        (_____|          |_|            ' >> /etc/motd && \
    echo '  ====================================================================' >> /etc/motd && \
    echo '  =  Welcome to the harvestgraph cli for Akamai                      =' >> /etc/motd && \
    echo '  ====================================================================' >> /etc/motd && \
    echo '  =  Project page:                                                   =' >> /etc/motd && \
    echo '  =  https://github.com/apiheat/harvestgraph                         =' >> /etc/motd && \
    echo '  ====================================================================' >> /etc/motd && \
    echo '  =  CLI versions:                                                   =' >> /etc/motd && \
    echo "     * $(harvestgraph --version)" >> /etc/motd && \
    echo "     * akamai appsec version v$(akamai appsec --version)" >> /etc/motd && \
    echo "     * $(akamai netlist --version)" >> /etc/motd && \
    echo "     * $(akamai --version)" >> /etc/motd && \
    echo '  ====================================================================' >> /etc/motd



ENV AKAMAI_CLI_HOME=/cli
VOLUME /cli
VOLUME /root/.edgerc

CMD ["/bin/bash"]