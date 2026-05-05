# curly — TUI HTTP Client

## Vision

Un client HTTP entièrement dans le terminal, piloté au clavier, sans souris.
L'alternative TUI propre à Insomnia/Postman pour les devs qui vivent dans leur terminal.

**Angle différenciateur** : chaînage de requêtes + replay diff + import Postman/Insomnia/Bruno.
Aucun outil TUI n'offre tout ça aujourd'hui.

## Stack

- **Langage** : Go
- **TUI** : [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Charm)
- **Styling** : [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Composants** : [Bubbles](https://github.com/charmbracelet/bubbles) (inputs, listes, viewport)
- **Config / collections** : [Viper](https://github.com/spf13/viper) + fichiers JSON locaux

## Architecture cible

```
curly/
├── main.go
├── ui/
│   ├── app.go            # modèle principal Bubble Tea
│   ├── panels/
│   │   ├── method.go     # sélecteur GET/POST/PUT/PATCH/DELETE
│   │   ├── url.go        # champ URL avec historique navigable
│   │   ├── headers.go    # édition des headers
│   │   ├── body.go       # corps de la requête + schema detection
│   │   └── response.go   # affichage réponse avec coloration JSON + diff
├── http/
│   ├── client.go         # exécution des requêtes
│   ├── schema.go         # inférence du schéma de body
│   ├── chain.go          # chaînage de requêtes (extraction + injection)
│   └── watch.go          # mode watch (replay périodique)
├── history/
│   └── store.go          # historique navigable de toutes les requêtes
├── collections/
│   ├── store.go          # sauvegarde/chargement des requêtes
│   └── import.go         # import Postman / Insomnia / Bruno
├── diff/
│   └── diff.go           # comparaison de deux réponses JSON
└── config/
    └── config.go         # configuration globale (env vars, thème...)
```

## Navigation clavier

| Touche | Action |
|--------|--------|
| `Tab` / `Shift+Tab` | Passer d'un panneau à l'autre |
| `↑` `↓` | Naviguer dans les listes |
| `Enter` | Sélectionner / confirmer |
| `Ctrl+R` | Envoyer la requête |
| `Ctrl+S` | Sauvegarder dans une collection |
| `Ctrl+E` | Gérer les variables d'environnement |
| `Ctrl+N` | Nouvelle requête |
| `Ctrl+P` / `Ctrl+N` | Naviguer dans l'historique (style shell) |
| `Ctrl+D` | Replay diff — comparer avec la dernière réponse |
| `Ctrl+W` | Mode watch — rejouer toutes les N secondes |
| `Ctrl+L` | Chaîner avec une autre requête |
| `?` | Aide |
| `q` / `Ctrl+C` | Quitter |

## Features différenciantes

### Schema detection (body auto-détecté)
Quand l'utilisateur sélectionne POST/PUT/PATCH et tape une URL :
1. Tentative de lecture d'une spec OpenAPI si disponible (best effort)
2. Sinon tentative d'un GET pour inférer le schéma depuis la réponse
3. Sinon mode manuel — l'utilisateur définit les champs et les sauvegarde
- Les champs s'affichent sous forme de formulaire navigable à la `Tab`
- Le JSON est généré automatiquement

### Replay & diff
- `Ctrl+D` rejoue la requête et affiche le diff ligne par ligne avec la dernière réponse
- Idéal pour détecter des changements silencieux d'une API

### Historique navigable façon shell
- `Ctrl+P` / `Ctrl+↑` remonte dans l'historique de toutes les requêtes envoyées
- Persisté entre les sessions

### Chaînage de requêtes
- Extraire un champ de la réponse avec une syntaxe JSONPath : `$.data.token`
- L'injecter dans la requête suivante via une variable : `{{chain.token}}`
- Parfait pour les flows login → appel authentifié

### Mode watch
- `Ctrl+W` rejoue la requête toutes les N secondes
- Affiche les changements en live — utile pendant le dev

### Import de collections
- Import direct depuis Postman (v2.1), Insomnia, et Bruno
- Les gens ne switchent pas sans leur historique — c'est le pont d'entrée

## Workflow de développement — étape par étape

> **Règle principale : une étape = une PR = un commit propre.**
> Ne pas passer à l'étape suivante sans que la précédente soit stable et commitée.
> À chaque étape, les concepts Go utilisés sont expliqués au moment où ils apparaissent.

### Étape 1 — Scaffolding & Hello World TUI
- Initialiser le module Go (`go mod init`)
- Installer Bubble Tea, Lip Gloss, Bubbles
- Afficher une fenêtre TUI vide avec un message de bienvenue
- **Concepts Go** : packages, `main()`, structs de base, `go run`
- **Livrable** : `go run .` ouvre le TUI sans planter

### Étape 2 — Layout de base
- Diviser l'écran en panneaux (méthode | URL | headers | response)
- Navigation Tab entre les panneaux avec highlight du panneau actif
- **Concepts Go** : interfaces, méthodes sur structs, `switch`
- **Livrable** : layout visible, navigation clavier fonctionnelle

### Étape 3 — Sélecteur de méthode HTTP
- Composant liste pour GET / POST / PUT / PATCH / DELETE
- Méthode sélectionnée colorée (GET=vert, POST=jaune, DELETE=rouge…)
- **Concepts Go** : constantes, `iota`, slices
- **Livrable** : on peut choisir une méthode au clavier

### Étape 4 — Champ URL + envoi de requête GET
- Input texte pour l'URL
- `Ctrl+R` envoie la requête GET
- La réponse JSON s'affiche avec coloration syntaxique
- **Concepts Go** : `net/http`, goroutines, channels, gestion d'erreurs
- **Livrable** : premier vrai appel HTTP fonctionnel

### Étape 5 — Headers
- Panneau d'édition des headers (clé / valeur)
- Ajout / suppression de headers au clavier
- **Concepts Go** : maps, boucles `range`
- **Livrable** : on peut envoyer des headers custom

### Étape 6 — Body + schema detection
- Panneau body activé uniquement sur POST/PUT/PATCH
- Schema detection en 3 niveaux (OpenAPI → inférence GET → manuel)
- Formulaire navigable générant le JSON automatiquement
- **Concepts Go** : `encoding/json`, interfaces vides, type assertions
- **Livrable** : la killer feature fonctionne

### Étape 7 — Historique navigable
- Toutes les requêtes envoyées sauvegardées localement
- `Ctrl+P` / `Ctrl+↑` pour remonter l'historique style shell
- **Concepts Go** : lecture/écriture de fichiers, `os`, `filepath`
- **Livrable** : historique persisté entre sessions

### Étape 8 — Replay & diff
- `Ctrl+D` rejoue et affiche le diff ligne par ligne avec la réponse précédente
- Coloration rouge/vert des changements
- **Concepts Go** : comparaison de structs, algorithme de diff simple
- **Livrable** : on voit immédiatement ce qui a changé

### Étape 9 — Collections
- Sauvegarder une requête (`Ctrl+S`) dans un fichier JSON local
- Panneau latéral listant les requêtes sauvegardées
- **Concepts Go** : sérialisation JSON, organisation en packages
- **Livrable** : persistance et organisation des requêtes

### Étape 10 — Variables d'environnement
- Définir des variables (`BASE_URL`, `TOKEN`…) par environnement (dev/prod)
- Syntaxe `{{VAR}}` dans l'URL et les headers
- **Concepts Go** : `strings.ReplaceAll`, Viper pour la config
- **Livrable** : workflow multi-environnement fonctionnel

### Étape 11 — Chaînage de requêtes
- Extraire un champ JSONPath de la réponse : `$.data.token`
- L'injecter dans la requête suivante via `{{chain.token}}`
- **Concepts Go** : parsing JSONPath, state management
- **Livrable** : flows multi-requêtes automatisés

### Étape 12 — Mode watch
- `Ctrl+W` rejoue toutes les N secondes
- Affichage live des changements
- **Concepts Go** : `time.Ticker`, goroutines avancées
- **Livrable** : monitoring d'endpoint en live

### Étape 13 — Import Postman / Insomnia / Bruno
- Parser les formats de collection des trois outils
- Importer en une commande : `curly import ./collection.json`
- **Concepts Go** : parsing JSON complexe, CLI args avec `cobra`
- **Livrable** : migration sans friction depuis les outils existants

### Étape 14 — Polish & release
- Écran d'aide `?` complet
- README avec GIF de démo (via [vhs](https://github.com/charmbracelet/vhs))
- Binaires pré-compilés via GitHub Actions (Linux, macOS, Windows)
- **Livrable** : première release publique sur GitHub

## Conventions de code

- Tout le code en anglais (noms de variables, commentaires, commits)
- Messages de commit en anglais, format conventionnel : `feat:`, `fix:`, `chore:`
- Une fonction = une responsabilité
- Pas de logique métier dans les fichiers UI
- Chaque étape = une branche Git + une PR avant de merger sur `main`

## Ressources utiles

- [Tour of Go](https://go.dev/tour/) — apprendre Go interactivement
- [Bubble Tea tutorial officiel](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Exemples Bubbles](https://github.com/charmbracelet/bubbles)
- [vhs — enregistrer des GIF de démo terminal](https://github.com/charmbracelet/vhs)
- [Awesome Bubble Tea](https://github.com/charmbracelet/awesome-bubbletea)
- [cobra — CLI args](https://github.com/spf13/cobra)
- [Format collection Postman v2.1](https://schema.postman.com/)
- [Format collection Bruno](https://docs.usebruno.com/collection/)