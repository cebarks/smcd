default:
    @just --list

test:
    go test -v -cover ./... 

run-server:
    SMCD_DIR=test/ go run ./cmd/smcd

run-client:
    SMCD_DIR=test/ go run ./cmd/smcd-ctl
