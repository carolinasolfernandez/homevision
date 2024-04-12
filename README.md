# HomeVision Backend Take Home

---

This is a take home interview for HomeVision that focuses primarily on writing clean code that accomplishes a very practical task. We have a simple paginated API that returns a list of houses along with some metadata. Your challenge is to write a script that meets the requirements.

**Note:** this is an unstable API! That means that it will likely fail with a non-200 response code. Your code *must* handle these errors correctly so that all photos are downloaded.

!https://media.giphy.com/media/QMHoU66sBXqqLqYvGO/giphy.gif


## Run it

```
    make deps
    make run
```

## Test it

```
    make test
```


## Environment Variables
You can set them as env variables or in the `config/app.env` file


## API Endpoint

You can request the data using the following endpoint:

```
http://app-homevision-staging.herokuapp.com/api_project/houses
```

This route by itself will respond with a default list of houses (or a server error!). You can use the following URL parameters:

- `page`: the page number you want to retrieve (default is 1)
- `per_page`: the number of houses per page (default is 10)

## Requirements

- Requests the first 10 pages of results from the API
- Parses the JSON returned by the API
- Downloads the photo for each house and saves it in a file with the name formatted as:

  `[id]-[address].[ext]`

- Downloading photos is slow so please optimize them and make use of concurrency

## Bonus Points

- Write tests
- Write your code in in a strongly typed language
- Structure your code as if you were planning to evolve it to production quality

## Managing your time

Include “TODO:” comments if there are things that you might want to address but that would take too much time for this exercise. That will help us understand items that you are considering but aren’t putting into this implementation. We can talk about what the improvements might look like during the interview that would get the code to final production quality.

## Submitting

- Create a GitHub repo with clear readme instructions for running your code on MacOS or Linux.
- Send us link to the public repo containing your submission (**preferred**) or a zip of the files.

Please let us know if you have any questions!

---

Thanks,

HomeVision Engineering