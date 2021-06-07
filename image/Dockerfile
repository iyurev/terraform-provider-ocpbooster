FROM golang:1.16 as builder

WORKDIR /workspace
COPY download-installer.sh .

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy the go source
COPY ./main.go main.go
COPY pkg/  pkg/
RUN chmod +x ./download-installer.sh && \
    ./download-installer.sh
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o terraform-provider-ocpbooster ./main.go

##Terraform
FROM  docker.io/hashicorp/terraform:0.15.4 as terraform
COPY image/providers.tf /tmp
RUN mkdir /providers && \
    cd /tmp && \
    terraform providers  mirror /providers

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4-200
ENV HOME="/root"
ENV TF_CLI_CONFIG_FILE="/root/.terraform.rc"
RUN microdnf -y install git-core jq
COPY --from=terraform  /bin/terraform /usr/local/bin/terraform
COPY --from=terraform  /providers /providers
COPY image/.terraform.rc  ${HOME}/.terraform.rc

RUN mkdir -p /providers/registry.terraform.io/hashicorp/ocpbooster/0.1/linux_amd64/
COPY --from=builder /workspace/terraform-provider-ocpbooster /providers/registry.terraform.io/hashicorp/ocpbooster/0.1/linux_amd64/

ENTRYPOINT "/bin/bash"
