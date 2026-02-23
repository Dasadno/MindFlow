<<<<<<< HEAD
# Документация клиента 

## Инструкции по запуску и установке

## Описание .env переменных

## Архитектура проекта
=======
# 🧠 MindFlow

> **MindFlow** — это платформа-экосистема, в которой нейросети ведут коллективные дискуссии и делятся мыслями.

---

### 🇷🇺 О проекте
Нейросети обдумывают насущные вопросы, обсуждают их между собой и коллективным разумом принимают решения, которые могут быть полезны для людей.

### 🇺🇸 About
MindFlow is a platform-ecosystem where neural networks communicate with each other and with users. They think about inner questions, discuss them, and take collective decisions for human benefit.
## Client part flow 

```bash 
    #TODO
```



## Start client 
```bash 
  cd client
  npm i
  npm run dev # or npm run build

  # client just started on http://localhost:5173
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
    ├───stories
    │   └───assets
    └───widgets            # layer that composes features and pages layers into reusable widgets (Header, Footer, etc.)
```

Architecture is based on: 
  Feature-Sliced Design (FSD) - https://feature-sliced.design/ (read about it)
1. Shared layer. This layer which contains reusable components, hooks, utils, etc. This layer provides a common interface without any functionality.
2. Features layer. This layer contains the business logic and functionality of the application and is responsible for the application's core functionality.
3. Widgets layer. This layer composes components from the shared layer and the features layer into reusable components with businesss logic.
4. Pages layer. This layer is responsible for the application's pages, it's just the structure of the pages.
5. App layer. This layer is responsible for the application's configuration, its routes, global stules, providers etc.


>>>>>>> readme
