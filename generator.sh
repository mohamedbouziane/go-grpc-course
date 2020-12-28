#!/bin/bash

protoc  --go_out=plugins=grpc:.   greet/greetpb/greet.proto

protoc  --go_out=plugins=grpc:.   blog/blogpb/blog.proto


# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=fr/ST=occitanie/L=toulouse/O=Moh/OU=mohtech/CN=mac/emailAddress=djmohamed1@gmail.com"

# Certificat info
openssl x509 -in ca-cert.pem -noout -text

# Generate server private key and certificate signing request (CSR)
openssl req  -newkey rsa:4096 -days 365 -keyout server-key.pem -out server-req.pem -subj "/C=fr/ST=Ile de france/L=Paris/O=Moh/OU=mohtech/CN=PC/emailAddress=dzdjul@gmail.com"

# Use CA'a private key to sign web server's CSR and get back the signed certificate

openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pm