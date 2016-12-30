# Hello Work

[![Build Status](https://travis-ci.org/italolelis/hellowork.svg?branch=master)](https://travis-ci.org/italolelis/hellowork)

Have you ever ran into a problem of not finding or not knowing where one of our work colleagues were?
Sometimes they were on vacations or having a business travel and you didn't know when they would be back. 
Imagine if you can notify your colleagues that you are leaving for a business trip, or maybe you got sick
and can't go to the office, just by telling slack about it?

That's the whole idea behind HelloWork. Notify your colleagues when you are out of the office and if someone 
mention you they will know who is out of the office and for how long.

## Notifing you are off

You can tell hellowork that you are off by using one of the follow commands.

```
@hellowork I'm on vacations until next month

@hellowork I'm on vacations until 20/02/2017

@hellowork I'm on vacations until next friday
```

You can tell @hellowork that you are on `vacations`, `business trip`, `out of the office` or `sick`.
Every time you say that you would need to tell from and until when you are not going to be available.

## Ask about someone

If you are curious to know where someone is, just ask hellowork
```
@hellowork Where is @wally?

@hellowork is @wally around?

@hellowork is @wally available?
```

## Instalation

You can choose to deploy this app with heroku. THis obviously the simplest way of doing it.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Another easy way of installing hellowork is to run the docker image for it. Just check the 
[docker-compose.yml](docker-compose.yml) example and then run.

```
docker-compose up -d
```

This will spin up 3 containers:

- *hellowork-bot* - This is the hellowork slack bot binary
- *hellowork-api* - This is the hellowork api, here is where the storage of data happens. This code is stored in [this repo](https://github.com/italolelis/hellowork-api)
- *postgres* - Our datastore of choice

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation
* Go lang: https://golang.org/
