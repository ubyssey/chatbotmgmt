#Deploying to Heroku

First make sure you have Node.js (which comes with npm as well) installed. 
Access to the heroku and gmail acc can be found in the Slack channel.

Then install the Heroku CLI which also comes with a Git installer. 

The app dependencies have already been installed. It is under dependencies in the package.json file. 
The version of Node.js that we are currently using is **6.9.5** (subject to change). It will be used to run the app on Heroku.

###Start Script
Heroku first looks for a Procfile. If that's unavailable, then Heroku attempts to start a web process via the
start script in **package.json**. I have currently set it up under start script because I was having some troubles
with the Procfile. This can always change later but it is easy to see right now. 

###Building your app and running it locally
1. Install your dependencies
'''
npm install
'''

2. Start your app locally. You should be able to see the processes with your shell. The app should be listening on **localhost:3000**.

'''
heroku local web
'''

###Deploying to Heroku
'''
git add -A 				// or whatever file you're adding
git commit -m "commit message"

heroku login				// if you're not logged in already

heroku create

git push heroku master

heroku open			
'''

**Make sure you're pushing the right commit. You can check this under Personal Apps in Heroku**

## Some deployment debugging
'''
heroku logs				// shows a log of your app
''' 

- check if dependencies are installed correctly or if the start scripts are set up correctly.
- set up npm-debug.log


