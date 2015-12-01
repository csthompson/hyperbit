# hyperbit
### tl;dr
A fast data analytics tool written in Golang that uses the Redis bitwise operations. Similar to bitmapist or minuteman.

### How it works:
A configuration file is used that holds the data analysis queries that are to be run on a Redis database. In addition to the queries, the configuration file would define the Redis indices that were used by Hyperbit.


### What is done so far:
1. Redis connection works
2. Data input into Redis works
3. Query parsing works

### What needs to be done:
1. Reading queries from configuration file correctly (.hbit file)
2. Running the queries quickly against a Redis database
3. Outputting meaningful data and creating a simple interface for users.
