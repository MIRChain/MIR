FROM centos:7 as builder

WORKDIR /mir

ENV GO_VERSION 1.22.3
ENV PATH ${PATH}:/usr/local/go/bin

COPY ["linux-amd64_rpm.tgz", "cades_linux_amd64.tar.gz", "kis_1", "cacerts.p7b","entry-point.sh","."]

RUN echo "Installing GO v${GO_VERSION}..." \
    && rm -rf /usr/local/go \
    && curl -sSLf https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && go version

RUN echo "Installing system packages..." \
    && sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-* \
    && sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-* \
    && yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
    && yum update -y \
    && yum install -y  lsb-core-noarch jq vim htop rsyslog net-tools git gcc gcc-c++ redhat-lsb-core jemalloc-devel gmp-devel  

RUN echo "Installing CryptoPro..." \
    && tar -xzf linux-amd64_rpm.tgz -C /tmp \
    && /tmp/linux*/install.sh kc1 cprocsp-stunnel-msspi lsb-cprocsp-devel \
    && mkdir /etc/opt/cprocsp/stunnel \
    && rm -rf /tmp/* linux-amd64_rpm.tgz

RUN echo "Installing CryptoPro Cades..." \
    && tar -xzf cades_linux_amd64.tar.gz -C /tmp \
    && cd /tmp/cades* && rpm -Uvh cprocsp-pki-cades*.rpm \
    && cd - && rm -rf /tmp/* cades_linux_amd64.tar.gz

COPY go.mod go.sum ./

RUN echo "Building MIR..." \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,mode=0755,target=/go/pkg \
    go mod download

# Copy the rest of the source code
COPY . .

RUN echo "Postintall & validation..." \
    && rpm -qa | grep csp \
    && /opt/cprocsp/bin/amd64/csptestf -enum -info \
    && chmod +x entry-point.sh

# Build the binary with caching
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN env GO111MODULE=on go run build/ci.go install ./cmd/mir

ENTRYPOINT ["./entry-point.sh"]