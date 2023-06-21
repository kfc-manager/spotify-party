# spotify-party

Spotify Party is a web application to control the queue of my personal spotify account. Initial motivation behind the project was to be able to share
a queue of someone, so people at a party have the possibility to add a song to the queue on their device and not have to pass around the smartphone of the
account owner. I started by developing the API first to supply the endpoints needed to control the queue. After that I decided to also implement a User
Interface, but to keep them coupled as losely as possible, so I don't have to make big changes to the API. The [React App](https://github.com/kfc-manager/spotify-party/tree/main/frontend) (frontend) is hosted on GitHub Pages via the [gh-pages branch](https://github.com/kfc-manager/spotify-party/tree/gh-pages)
and the [dockerized API](https://github.com/kfc-manager/spotify-party/tree/main/backend) (backend) is hosted as Google Cloud Run Service [here](https://spotify-party-zty7jo4vkq-ey.a.run.app). It is also possible to host the app by yourself and hook up you own Spotify Account.

## Getting started

Since my API is using the REST API provided by Spotify we have to follow their rules of authorization. To be able to share the access of your account you have
to login to your Spotify Account and register an application [here](https://developer.spotify.com/dashboard/create). The name and description of the application
are not important but required fields. The important part is the Redirect URI, if you want to host the app localy in your network you can put "localhost",
the Port on which you plan to host the API then followed by "/callback". Make sure that the Redirect URI is correct otherwise the API won't be able to request
an Access Token in order to be able to access you account.

**Example:**

![alt text](https://github.com/kfc-manager/spotify-party/blob/main/register-app.png?raw=true)

Next up clone this repository and create a ".env" file in [spotify-party/backend/](https://github.com/kfc-manager/spotify-party/tree/main/backend). You have
to assign values to the vriables PORT, SPOTIFY_ID, SPOTIFY_SECRET and REDIRECT_URL inside the ".env" file. The PORT must be set to the Port you already used in the Redirect URI while registering the app. For SPOTIFY_ID put in the Client ID and for SPOTIFY_SECRET the Client secret. Both of them you can find under "SETTINGS"
in your registered application. Use the Redirect URI you just registered for REDIRECT_URL. For my example the ".env" file would be:

```
PORT=8080
SPOTIFY_ID=329c0d9b1ce64024a9a6dfc6447d9e0a
SPOTIFY_SECRET=56a77c529c8c4ed6b019214ee7736ba9
REDIRECT_URL=http://localhost:8080/callback
```

Now you can compile and run the Go Code by using the command:

```
make run
```

inside the working directory. If you want to use Docker and run the API through a Container use:

```
make image
```

to create a Docker Image. Before you can use the endpoints to control the queue, you have to hit the endpoint "/login".
You will be redirected to Spotify's Authorization Website and have to login with your Spotify Account. After that it should say "Authorization successful". 
The API now has an Access Token and you can use the other endpoints. The Token expires after 60 minutes and you have to request a new one with the method
explained. This is due to the Spotify Authorization Flow (more to that below).
If you also want use my User Interface, navigate into [spotify-party/frontend](https://github.com/kfc-manager/spotify-party/tree/main/frontend) and then use:

```
npm install
```

to install all Node Modules and then run:

```
npm run dev
```

to start the React App. The app is now accessible with every device which is in your network (for example smartphones connected to your WLAN), with the IP Address 
of the machine you host the app on.
