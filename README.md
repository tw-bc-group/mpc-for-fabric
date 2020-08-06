## Quick start
This repo now demonstrate how to use HE(homomorphic encryption) on the fabric for multiparty computation to do financial exchange without unnecessary data reveal.

Ask and Bid will exchange limit price of an order as an encrypted ciphertext, credit will be kept as secret in server during substraction computation, credit will be stored only in Server side, while the exchange only revealing distance of credits with thirdparty bids

```
                 {                         {
                   LimitPrice(ciphertext)    LimitPrice(ciphertext)
                 }                           Credit(ciphertext)
   +------------+           +------------+ }         +-------------+
   |            |           |            |           |             |
   |   client   +-- Ask --->+   server   +<-- Bid ---+ third_party |
   |            |           |            |           |             |
   +------------+           +------------+           +-------------+
          ^                     Credit(ciphertext)
          |                        +
          |                        |
          |                        |
          +---- EligibleBid -------+

               {
                LimitPriceDistance = Ask.LimitPrice - Bid.LimitPrice (invisible to server/third_party)
                CreditDistance = Credit - Bid.Credit (invisible to server/third_party)
               }
```


## Deploy local fabric test network

### Prerequisites

- Git

- cURL

- Docker and Docker Compose

- Golang

  

#### Note

If you are using Docker Toolbox or macOS, you will need to use a location under `/Users` (macOS) when installing and running the samples.

If you are using Docker for Mac, you will need to use a location under `/Users`, `/Volumes`, `/private`, or `/tmp`. To use a different location, please consult the Docker documentation for [file sharing](https://docs.docker.com/docker-for-mac/#file-sharing)


## Run Demo

```bash
./deploy.sh
```

```bash
./setAsk.sh
```

```
./setBids.sh
```

```bash
./getBids.sh
```