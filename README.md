# **WHATFLIX**

Search and Movie Recommendation Engine.                    Build a horizontal scaling Go application using decentralised and distributed microservices archictecture.                 All the services are build using native golang except gorilla mux package for HTTP routing.

Following are the details of the services:-

* Webservice: Main web application.
* LoadBalancer
* CacheService
* LogService

## Build - Docker Images

* Create a network
docker network create --subnet=127.0.0.1/16 whatflixnet
docker network ls

### Build and run images of services

* #### LogService

docker build -t whatflix/logservice -f Dockerfile-logservice .
docker run --name logservice --ip=127.0.0.10 --net=whatflixnet -P --rm -it -v /Users/syedqrizvi/sendinblue/goworkspace/src/github.com/whatflix/logservice:/log whatflix/logservice  

* #### CacheService

docker build -t whatflix/cacheservice -f Dockerfile-cacheservice .
docker run --name cacheservice --ip=127.0.0.11 --net=whatflixnet -P --rm -it whatflix/cacheservice

* #### LoadBalancer

docker build -t whatflix/loadbalancer -f Dockerfile-loadbalancer .
docker run --name loadbalancer --ip=127.0.0.12 --net=whatflixnet -P --rm -it -- whatflix/loadbalancer --logservice=<http://127.0.0.10:6000>

* #### WebMain

docker build -t whatflix/web -f Dockerfile-web .
docker run --name web --ip=127.0.0.13 --net=whatflixnet -P --rm -it -- whatflix/web --loadbalancer=<http://127.0.0.12:2001> --cacheservice=<http://127.0.0.11:5000> --logservice=<http://127.0.0.10:6000>

## Build - Using native go build and run the application using postman

Build all the four services using `go build` command and run as follows:

* Ping: GET Request <http://localhost:2000/ping>
  
* SignIn: POST Request with body username and password <http://localhost:2000/movies/signin>

* Search Engine: GET Request <http://localhost:2000/movies/user/101/search?text=Tom Hanks>

* Recommendation Engine: GET Request <http://localhost:2000/movies/users>

## Data Sets

I have used the Kaggle TMDB data set for this problem. Complete movie data set can be downloaded from here [LINK](https://www.kaggle.com/tmdb/tmdb-movie-metadata)

The user preferences data file can be downloaded from here
[LINK](https://github.com/qasimhbti/whatflix/blob/master/user_preferences.json)

## 1. Search Engine

    Restful web service API which will accept a search string and userID and return unique movies in the order of preference for that user.

    URL : http://<url>/movies/user/$userId/search?text=<text>
    Request Type : GET

    Where the $userId is the id of the user from the user preferences JSON file (see file) and <text> is the search text. The search text could be multiple words separated by comma, in which case it will search for all of those. Or it could be a single urlencoded entry. The search text is matched against actor, director and title fields (director_name, actor_1_name, actor_2_name, actor_3_name and movie_title in the movies data) of the movies. All matches are included in the results.
    
    <url> is the domain name in the URL, hosting the solution.
    For example - http://<url>/movies/user/$userId/search?text=Tom%20Hanks -> This will return all the movies matching “Tom Hanks” considering the preferences of the user.

    An array of movies name, found based on the preferences of users, sorted in this order:

      * First show the movies matching the user’s preferences and search term. This should be further sorted on the alphabetic order of the titles. There is a chance this set of movies could be empty if there is no search result matching the user’s preferences.

      * Next show the movies matching the search term in the alphabetic order of titles (even if it does not match the user’s preferences). If #1 is empty then only these set of movies will be there.
      [“Movie 1”,”Movie 2”,”Movie 3”]

## 2. Recommendation Engine

    Restful web service API which will list out all the users and the top 3 recommended movies for each user based on their preferences.

    URL : http://<url>/movies/users
    Request Type : GET

    A json array of userids and top 3 movies recommended. The movie names are sorted in the alphabetic order of titles. If
    there is no recommendations then it can be an empty array.

    ```
    [
        {
            "user": "100",
            "movies": [
                "A",
                "B",
                "C",
            ]
        },
        {
            "user": "101",
            "movies": [
                "A",
                "D",
                "C",
            ]
        },
    ]
    ```
