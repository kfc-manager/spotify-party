# spotify-party

Spotify Party is a web application to control the queue of my personal spotify account. Initial motivation behind the project was to be able to share
a queue of someone, so people at a party have the possibility to add a song to the queue on their device and not have to pass around the phone of the
account owner. I started by developing the API first to supply the endpoints needed to control the queue. After that I decided to also implement a User
Interface, but to keep them coupled as losely as possible, so I don't have to make big changes to the API. The [React App](https://github.com/kfc-manager/spotify-party/tree/main/frontend) (frontend) is hosted on GitHub Pages via the [gh-pages branch](https://github.com/kfc-manager/spotify-party/tree/gh-pages)
and the [dockerized API](https://github.com/kfc-manager/spotify-party/tree/main/backend) (backend) is hosted as Google Cloud Run Service [here](https://spotify-party-zty7jo4vkq-ey.a.run.app). It is also possible to host the app by yourself and hook up you own spotify account.
