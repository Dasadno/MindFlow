# MindFlow

## Client part flow 

```bash 
    #TODO
```



## Start client 
```bash 
  cd client
  npm i
  npm run dev
  # client just started on http://localhost:5173

  # for view production build
  npm run build
  npm run preview

  # build preview just started on http://localhost:4173
```


## Architecture 

```
CLIENT
├───.storybook
├───public            #static files  
└───src               #source code
    ├───app               #app configuration 
    │   ├───router    
    │   └───styles
    ├───assets
    ├───features          #functionality (auth, home, etc.)
    │   └───auth
    │       └───ui
    │           └───LoginForm
    ├───pages             #pages (home, about, etc.)  
    │   ├───auth
    │   │   └───ui
    │   └───home
    │       └───ui
    ├───shared             #shared code (components, hooks, utils, etc.)
    │   └───ui
    │       └───Button
    ├───store                      
    └───widgets            # layer that composes features and pages layers into reusable widgets (Header, Footer, etc.)
```

Architecture is based on: 
  Feature-Sliced Design (FSD) - https://feature-sliced.design/ (read about it)
1. Shared layer. This layer which contains reusable components, hooks, utils, etc. This layer provides a common interface without any functionality.
2. Features layer. This layer contains the business logic and functionality of the application and is responsible for the application's core functionality.
3. Widgets layer. This layer composes components from the shared layer and the features layer into reusable components with businesss logic.
4. Pages layer. This layer is responsible for the application's pages, it's just the structure of the pages.
5. App layer. This layer is responsible for the application's configuration, its routes, global stules, providers etc.


## UI

### Login page

<img src="./src/assets/screenshots/login.png" alt="login page" width="500" height="300">

### Registration page

<img src="./src/assets/screenshots/register.png" alt="login page" width="500" height="300">


### Main page

<img src="./src/assets/screenshots/mainPage1.png" alt="main page1" width="500" height="300">

<img src="./src/assets/screenshots/mainPage2.png" alt="main page2" width="500" height="300">

#### Some information about us and our project

<img src="./src/assets/screenshots/mainPage3.png" alt="main page3" width="500" height="300">


### Error page

<img src="./src/assets/screenshots/error.png" alt="main page3" width="500" height="300">


### Chat Page

#### How chats are look like
<img src="./src/assets/screenshots/chat1.png" alt="main page3" width="500" height="300">

<img src="./src/assets/screenshots/chat2.png" alt="main page3" width="500" height="300">

#### Sidebar with information of your agemts

<img src="./src/assets/screenshots/sidebar.png" alt="main page3" width="500" height="300">

#### How you can create a new agent

<img src="./src/assets/screenshots/agentCreate.png" alt="main page3" width="500" height="300">