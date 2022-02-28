# Ditto
Ditto is a cli tool to mock an API or a kafka producer.

## Schema
Supported fields:
- `INTEGER`:
  - `VALUE`: hardcoded value
  - `INCREMENTAL`: from a `START` value, increase by a fixed `STEP` 
    - NOTE: An optional `SIGMA` can be passed to make the delta follow a normal distribution around the `STEP`
  - `DISTRIBUTION`:
    - `NORMAL`: returns a value for a given `MU` and `SIGMA`
    - `UNIFORM`: returns a value between a `MIX` and `MAX`
    - `RANDOM`: returns a random integer, optionally limited by a `MAX` value
- `FLOAT`:
  - `VALUE`: hardcoded value
  - `INCREMENTAL`: from a `START` value, increase by a fixed `STEP`
    - NOTE: An optional `SIGMA` can be passed to make the delta follow a normal distribution around the `STEP`
  - `DISTRIBUTION`:
    - `NORMAL`: returns a value for a given `MU` and `SIGMA`
    - `UNIFORM`: returns a value between a `MIX` and `MAX`
    - `RANDOM`: returns a random float between (0,1], optionally multiplied by a `SCALE` value
  - Additional Params:
    - `ROUND`: Number of decimal values. Defaults to 2
- `STRING`:
  - `VALUE`: hardcoded value
  - `REGEX`: regex template to generate a random string
  - Additional Params:
    - `EMPTY`: Probability (`(0, 1]`) of the string being empty. Defaults to 0
- `TIMESTAMP`:
  - `VALUE`: hardcoded value
  - `NOW`: current timestamp
  - `INCREMENTAL`: from a `START` value, increase by a fixed `STEP`. The start value defaults to `NOW`
    - NOTE: An optional `SIGMA` can be passed to make the delta follow a normal distribution around the `STEP`
    - NOTE: The step should follow the format "\d+[hms]"
  - Additional Params:
    - `PRECISION`: Units of the timestamp. Supported values are "s" (default), "ms", "us", "ns"
- `DATETIME`: 
  - `VALUE`: hardcoded value
  - `NOW`: current timestamp
  - `INCREMENTAL`: from a `START` value, increase by a fixed `STEP`. The start value defaults to `NOW`
    - NOTE: `STEP` should follow the format "\d+[hms]"
    - NOTE: `START` should be a unix timestamp
    - NOTE: An optional `SIGMA` can be passed to make the delta follow a normal distribution around the `STEP`
  - Additional Params:
    - `FORMAT`: Format of the datetime string (using the go reference date)
- `OBJECT`:
  - `FIELDS`: Array with the nested fields
  
Additionally, there are certain parameters common between all types:
  - `NULLABLE`: Probability (`(0, 1]`) of the value being null. Defaults to 0
  - `REPEATED`: Whether the value is an array. It will create an array between 0 and 10 elements.

## Usage

### Rest
To serve message as an rest api:

``` sh
ditto http server -s schema.json -p 8080 
```

If the schema is a folder, it will create a unique endpoint for each schema present using the name of the file as the url.

### Kafka
To send messages to kafka (requires a running Kafka cluster):

``` sh
ditto kafka producer -s schema.json -b localhost:9092 -S 3 [-c 2] [-t sample]
```

If the schema is a folder, it will send messages to a topic that matches the name of the file. If the `-t` flag is passed, it will send all messages to that topic.


## Roadmap
- [ ] Cli command to generate schema from json
- [ ] Consume messages and verify that they follow the defined schema
- [ ] Add error support (random error codes, wrong types,...)
- [ ] Support other file formats
- [ ] Return arrays
- [ ] Tests
- [ ] Support other outputs (GraphQL,...)
