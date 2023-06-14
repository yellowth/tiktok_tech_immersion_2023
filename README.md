# TikTok Tech Immersion 2023

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This project is the backend assignment for the TikTok Tech Immersion Program 2023. It is a simple Instant Messaging platform backend built on Golang, Protobuf, Kitex, Redis and Docker. It is built using [this](https://github.com/TikTokTechImmersion/assignment_demo_2023) demo template.

## 1. Setup
Run: 
```bash
docker compose up -d
```
in the same directory as
`docker-compose.yml`:

## 2. Send messages
Send a `POST` request to `localhost:8080/api/send` (using Postman or curl) in this format:

eg. 
Send message "hello" from alex to ben:
```bash
curl -X POST 'localhost:8080/api/send?chat=ben:b&sender=alex&text=hello'
```

## 3. Pull messages
Send a `GET` request to `localhost:8080/api/pull` in this format:

eg. 
Pull all messages between alex and ben:
```bash
curl 'localhost:8080/api/pull?chat=alex:ben&cursor=0&limit=10&reverse=false'
```
** cursor is the earliest epoch to retrieve messages, 0 by default

** limit is the max number of messages to retrieve, 10 by default

** reverse decides whether messages are sorted in ascending order (earliest first)
