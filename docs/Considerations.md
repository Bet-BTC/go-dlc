# Considersations


## !!!Do not use this library without considering the following points!!!

There are critical considersations if you use this library in the mainnet.

## Oracle's private key management

Oracle's private key must be kept safe when running oracle server. [The current implementation](./internal/oracle/oracle.go) doesn't support mainnet.

## Wallet key management

Currently, this library's private key generation is not safe for production/mainnet environments.