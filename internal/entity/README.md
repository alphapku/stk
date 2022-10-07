# Sub-directory structures
Each source has its own sub-directory here. For example, `mock` is for the mock source.

Also, we(Stake) has a sub-directory here - `stake` for internal use.

All messages from different sources have to be converted into Stake's format. The task is done by corresponding adapters.

So, basically, the data conversion  happens in this way:
```
Adatper message1 -> Stake message1
                                   \
                                    Combined Stake Message for clients
                                   /
Adatper message2 -> Stake message2
```