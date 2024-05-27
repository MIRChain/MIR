FROM centos:7 as builder

WORKDIR /mir

ENV GO_VERSION 1.22.3
ENV PATH ${PATH}:/usr/local/go/bin

COPY ["linux-amd64_rpm.tgz", "cades_linux_amd64.tar.gz", "kis_1", "cacerts.p7b","."]

RUN echo "Installing GO v${GO_VERSION}..." \
    && rm -rf /usr/local/go \
    && curl -sSLf https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && go version

RUN echo "Installing system packages..." \
    && yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
    && yum update -y \
    && yum install -y  lsb-core-noarch jq vim htop rsyslog net-tools git gcc gcc-c++ redhat-lsb-core jemalloc-devel gmp-devel  

RUN echo "Installing CSP..." \
    && tar -xzf linux-amd64_rpm.tgz -C /tmp \
    && /tmp/linux*/install.sh kc1 cprocsp-stunnel-msspi lsb-cprocsp-devel \
    && mkdir /etc/opt/cprocsp/stunnel \
    && rm -rf /tmp/* linux-amd64_rpm.tgz

RUN echo "Installing cades..." \
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
    && /opt/cprocsp/bin/amd64/csptestf -enum -info

# Build the binary with caching
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN env GO111MODULE=on go run build/ci.go install ./cmd/mir


FROM centos:7 

# Copy the built binary and entry-point script from the previous stage/build context
COPY --from=builder /mir/build/bin/mir /bin/mir
COPY entry-point.sh /entry-point.sh

# Ensure the entry-point script is executable
RUN chmod +x /entry-point.sh

ENTRYPOINT ["/entry-point.sh"]