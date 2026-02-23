# Start client 
```bash 
  cd client
  npm i
  npm run dev # or npm run build

  # client just started on http://localhost:5173
```

# Architecture 

CLIENT
├───.storybook

├───public

└───src
    ├───app
    │   ├───router
    │   └───styles
    ├───assets
    ├───features
    │   └───auth
    │       └───ui
    │           └───LoginForm
    ├───pages
    │   ├───auth
    │   │   └───ui
    │   └───home
    │       └───ui
    ├───shared
    │   └───ui
    │       └───Button
    ├───store
    ├───stories
    │   └───assets
    └───widgets
