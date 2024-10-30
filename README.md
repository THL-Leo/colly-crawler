# CSCI 572 HW 2 by Leo Lee

>Please refer to this document if you have any questions regarding the report. Thank you!

## FAQ for the assignment

### Why are there 625 GoRoutines instead of threads?
I did this assignment in Golang so I put what's the equivalent counterpart in Golang in the report. Technically there is a difference between threading and concurrency but it will work for this assignment. I have 625 GoRoutines which means there can have 625 concurrent requests at the same time.

### Why are there around 4.5 million URLs?
My guess is that 1. LA Times is a very large website with lots of Outlinks and 2. since I am using 625 GoRoutines, the queue is being populated by 625 workers which can dramatically increase the number of links we insert into our queue. It is like filling up a bucket with a faucet but we are removing the water using a spoon.

### How did you handle the redirects?
There are two paths that a redirect can go: 1) still within the domain so the redirect is handled by the RedirectHandler which returns the last response and 2) the link is not within the domain anymore so the HTTP status code is 0. If the status code is 0, it can mean a lot of things. However, most of the time it will be ERROR: The redirect link is not within the AllowedDomain list, so it won't be redirected. What I did is that when these status code of 0 happens, I just decrement the PageCount by 1 and not log it in Fetch.csv. Since Visit.csv will only be logged if the page is traveled to, this allows me to skip the incorrect redirects and only log pages that are within the domain.

## Concluding thoughts
This assignment is actually really fun as I learned a new language. I thoroughly enjoyed the assignment tasks and I am really glad I chose Go as the language to do this assignment. The speed of Go is insane when compared to Scrapy from Python and Crawler4j from Java. Thank you for the hardwork grading these assignments.