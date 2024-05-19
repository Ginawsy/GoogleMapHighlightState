# MongoDB

## Setup

**Setup docker image**
* https://www.mongodb.com/docs/manual/tutorial/install-mongodb-community-with-docker/

**Import dataset**
* https://github.com/jakevdp/data-USstates

**Validate dataset is inside the database**
```
./mongosh --port 27017
use project_dataset

db.us_states.find({'state':'California'})
```


# Vue.js
* Watcher
* Form Bindings
```
npm install -g @vue/cli
npm create vue@latest
cd min_project_fronend
npm install
npm run dev
```









