# What is this

this project is to visualize current number of confirmed [SARS-CoV-2](https://en.wikipedia.org/wiki/Severe_acute_respiratory_syndrome_coronavirus_2) infections and deaths
in Germany as well as providing a raw data API

# How does it work

this project consist of a server which scrape [Robert Koch Institute](https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html) every 9 minute and export the data for Promethues. The scraped data can also be accessed via `/api/epidemic`. 

Prometheus service will scrape the exported data and Grafana is used to visualize the result

# Run the server

the server can be run via a single command (you need to have docker installed)

```sh
docker-compose up
```

the docker-compose will run the server as container and expose port `8404`. Also it will spin up
Prometheus(port 9090) and Grafana(port 3000). 

>if you see any permission error from prometheus or grafana service, make sure `data` folder is writable by the user.

# Hosted service

a hosted dashboard can be found [Here](http://bit.ly/corona-weg) and the data API can be found [Here](http://bit.ly/corona-weg-api). It's a AWS EC2 micro VM. please be gentle with it.

![image](https://user-images.githubusercontent.com/3006506/75823439-11075600-5da2-11ea-8d81-c4e8d13ebed0.png)

