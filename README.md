# Golang Challenge (Transparent Cache)

This repo contains the solution to the [Golang-Challenge](https://github.com/deviget/Golang-Challenge). Feel free to head over it to get a better context of the problem :book:
  
## Approach :nerd_face:
The approach is aligned with the following considerations:
1. Fix the existing broken tests.
    1. Improve the existing transparent cache code with no breaking changes
    2. Don't touch the existing tests
    3. Add missing functionalities
2. Add more tests

Generally speaking, this is a kind of resumed approach. In order to fix the tests, we need to add new features (___TDD___). It's so common when someone wants to fix something, it breaks something else (especially when we have productive and tested code), so I took the problem carefully. In this case, the failing tests are because of the lack of certain functions that are missing, and I had to implement it. So far so now, we are good. In the code, there were a couple of things that could be improved, but the less code we touch, the fewer issues we'll get (this rule may be omitted if the change is minimal and don't break the tests or if it's a scalable change).


## Diving into details :diving_mask:
### Expiring Results :hourglass:
The requirement talks about fix the test by adding the missing feature. Caching results with the existing __TransparentCache__ implementation can't be done because there isn't any struct, variable, or something that we can use to map a price with the expired time. So, as a ___cachedAt___ variable is needed in order to know when each price has been added, I decided to create a Price struct and put that property over there. This __Price__ struct now will give us more flexibility because we can store more properties like ___cachedAt__, ___modifiedAt___, ___type___, etc... This will scale better. 
<br />

To avoid change all tests by modifying the ___float64___ price references to the Price structure because only the price value is being tested, I decide to keep that as it is and only uses the __Price__ reference inside the [cache.go](https://github.com/morarick/transparent-cache/blob/main/cache.go) file. <br />
If within the near future there are validations that match with the __Price__ properties in the tests, we work with the __Price__ structure in the tests instead of the ___float64___ price value.
<br />

In short, this feature is really easy. I just had to create the __Price__ struct, make the validation for ___cahedAt___ < ___maxAge___ (for cached values) and set the current time to ___cachedAt___ (for non-cached yet values)
No further magic was needed and the test [TestGetPriceFor_DoesNotReturnOldResults](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache_test.go#L130) works like a charm :fire:

### Parallelize calls :dancing_men:
This feature is linked to this failing test [TestGetPricesFor_ParallelizeCalls](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache_test.go#L162) and to fix it, I basically used the power of goroutines. The approach simple and scales pretty well.<br />
Moreover, there are several considerations that we have to keep in mind before start parallelizing things. Working with maps is one of them because they are not thread-safe.<br />
>Maps are not safe for concurrent use: it's not defined what happens when you read and write to them simultaneously. If you need to read from and write to a map from concurrently executing goroutines, the accesses must be mediated by some kind of synchronization mechanism. One common way to protect maps is with sync.RWMutex.<br />

So, I've created the [loadPriceSync](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache.go#L65) and [storePriceSync](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache.go#L73) that will handle safely the lock/unlock mechanism for the __TransparentCache__.___prices___ map.<br />
I decided to use ___sync.RWMutex___ because is preferable for data that is mostly read, and the resource that is saves compared to a ___sync.Mutex___ is less time.<br />

The main functionality of getting prices at once is straightforward: publish each price to the queue ([publishPrice](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache.go#L90)), and consume them ([consumePrices](https://github.com/morarick/transparent-cache/blob/5215acbf7d366538f43aca3954913b87b3ef99f2/cache.go#L96)) to be delivered.<br />

>Note: I went for the publish/consume approach using a channel because is the easiest, idiomatic, and best fits our problem (instead of using other complicated solutions).

## Development cycle :bicyclist:
* TDD
  * In the beginning, the challenge enforce you to apply __TDD__. To fix the tests, you have to add missing features (this is a really good practice :tada:)
* Conventional Commits
  * In order to follow some good practices, I made use of [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) to standardize the commit messages in a clearer way.<br />

## Testing :test_tube:
There are several tests that I added. Also, below I'll leave some interesting topics about the tests in the ___TransparentCanche___ implementation.<br />

### Code Coverage :rainbow:
I've covered 100% of test coverage for [cache.go](https://github.com/morarick/transparent-cache/blob/main/cache.go).<br />
To check for this, we can run the following command in a terminal:
```go test -cover```

### Data Races :runner:
As we are working with goroutines is really important to ensure that there aren't data races or race conditions in our code.<br />
To check for this, we can run the following command in a terminal:
```go test -race```

Thanks for reading :v:
