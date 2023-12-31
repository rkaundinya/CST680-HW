The concept of HATEOAS is fairly simple – it’s essentially to let the API tell you how to interact with it rather than you figure it out yourself. This is driven via a state machine type manner where the details of the interaction available to you are driven by the state of the API – or more accurately the state of the objects the API manages. 

The states you have in our API structure are – no voter created, no poll created, and no poll or voter with ID created. This are all the error conditions that can cause problems with API interaction. 

When a POST for a vote is created, we can do a check in the main code file whether a voter with the given id exists in the voter list and whether such a poll with ID exists. If both exist, then we can successfully create the vote and return: the URL links to see the voter’s vote history (/voters/ID/polls), a link with type “POST” to /voters/ID/polls in order to post the vote to the VoteHIstory, and a URL to view the poll with ID with type GET. Ideally, the creation of the vote automatically adds to the vote history. If that’s the case, then we just return /voters/ID/polls with type “GET” and /polls/ID with type “GET” to let the user know how to view the vote they just created and registered and the poll they voted on.

Now it’s possible you create a vote but are trying to do so with an invalid Voter ID. If this is the case, then the HATEOAS should recognize this and return to you a link to /voters/ID with type “POST” and the user should be notified via an HTTP error that they are making an invalid request. They can also be given a link to /voters with type “GET” to tell them how to view all the current voters. 

Similarly, if trying to make a vote for a poll which does not exist, then return an HTTP error, a link to /polls and type “POST” and /polls with type “GET” to tell the user how to post and view polls. 

If both voter ID and poll ID are wrong, then do both of the above. 

This takes care of the major failure cases. The rest of the design would all be convience/communication to the user on how to interact. For example, when deleting a poll, return links to view the voters to see their updated vote history. Or when deleting a voter, return a link to see votes for a poll which may be updated to no longer include the deleted voter. 

The votes for polls would probably be stored separately and this would need to be added to the http URL convention so that the user could view all the collected votes for a specific poll. The polls should not keep track of votes within themselves as this would violate the single responsibility principle and create an unwanted dependency. Any update on votes would have to update an internal state in the poll API. Rather than this, the poll api should continue to define simply what a poll is and what its options are and it is up to the user of the polls to separately manage votes for this poll. This also allows some flexibility on how the user wants to create, store, and manage votes. 

So overall, the HATEOAS convention is mainly to allow the user to use the API as if they were browsing a website – all they have to do is click around and the website will tell them and link them to places that are relevant for them to go to. This convention doesn’t really drive software architecting in the sense that you should already have good design built regardless of whether you choose to use HATEOAS or not, but it can certainly help you catch instances of poor design you did not catch otherwise and thereby enforce better design principles. 

I mainly got the gist from the Wikipedia article, but here are two other links which had the detail of returning the type of REST interaction in the json as well as the URI link. 

https://restfulapi.net/hateoas/
https://www.w3schools.in/restful-web-services/rest-apis-hateoas-concept 
