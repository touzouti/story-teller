# Story Teller

Story Teller est une application web de création d'histoires du soir. L'utilisateur choisit un univers, un ton, un décor, un compagnon et une durée. Le frontend transmet ensuite ces choix au backend, qui peut appeler un webhook externe pour générer l'histoire et conserver les dernières histoires dans un fichier JSON.



## Fonctionnalités

- page d'accueil animée ;
- six catégories d'histoires ;
- choix de la durée de lecture ;
- envoi des choix au backend ;
- affichage de l'histoire générée ;
- réglage de la taille du texte ;
- sauvegarde des 50 dernières histoires ;
- consultation des histoires sauvegardées ;
- exécution locale avec Docker Compose.

## Technologies

| Partie | Technologies |
| --- | --- |
| Frontend | Vue 3, TypeScript, Vite |
| Interface | Tailwind CSS, DaisyUI, Lucide |
| Animations | CSS, Vue3 Lottie |
| Backend | Go, API HTTP |
| Authentification du webhook | JWT facultatif |
| Stockage actuel | `backend/data/stories.json` |
| Conteneurs | Docker, Docker Compose, Nginx |

## Architecture

```text
Navigateur
   |
   | /api/story et /api/stories
   v
Frontend Vue servi par Nginx
   |
   v
Backend Go
   |-- appelle WEBHOOK_URL si cette variable existe
   |-- utilise fallbackStory() dans le cas contraire
   `-- sauvegarde les histoires dans data/stories.json
```

Le code du backend accepte plusieurs formes de réponse provenant d'un webhook, notamment :

```json
{
  "title": "Le titre de l'histoire",
  "story": "Le texte de l'histoire"
}
```

Le code contient également un traitement de réponses au format Gemini/n8n. Le workflow n8n et sa configuration ne sont cependant pas inclus dans ce dépôt.

## Structure du projet

```text
story-teller/
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── docker-compose.yml
└── README.md
```

## Prérequis

Pour la méthode recommandée :

- Windows 10 ou 11, macOS ou Linux ;
- Docker Desktop avec Docker Compose ;
- Git ;
- VS Code, facultatif mais conseillé.

Pour lancer les deux parties sans Docker, il faut également Node.js et Go.



