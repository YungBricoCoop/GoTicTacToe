// SPDX-FileCopyrightText: 2025 Haute école d'ingénierie et d'architecture de Fribourg
//
// SPDX-License-Identifier: CC0-1.0

#import "lib/template.typ": report
#import "@preview/big-todo:0.2.0": *
#import "@preview/glossarium:0.5.4": gls, glspl, make-glossary, print-glossary, register-glossary

#set text(lang: "fr")
#let theme_color = rgb(0, 124, 183, 255)

#show: make-glossary
/* #let entry-list = ()
#register-glossary(entry-list) */

#let authors = (
  (
    firstname: "Elwan",
    lastname: "Mayencourt",
    email: "elwan.mayencourt@edu.hefr.ch",
    gender: "m",
  ),
  (
    firstname: "Masami",
    lastname: "Morimura",
    email: "masami.schoellkopfmorimura@edu.hefr.ch",
    gender: "f",
  ),
)

#let versions = (
  (
    version: "1.0",
    date: datetime(year: 2026, month: 01, day: 11),
    author: [EM],
    changes: [Version finale],
  ),
)


#show: doc => report(
  title: "Gopher Dungeon",
  type: "Rapport de projet",
  theme_color: theme_color,
  authors: authors,
  versions: versions,
  supervisors: ("Supcik Jacques",),
  cover_image_path: "../assets/images/ai_fact_check.png",
  doc,
)

= Introduction

Le projet choisi transforme le jeu classique de Tic-Tac-Toe en remplaçant entièrement sa couche visuelle 2D par une interface 3D basée sur le raycasting. Au lieu d'afficher une simple grille en 2D, le joueur se déplace dans un environnement 3D où chaque case du jeu est représentée par une salle dans un labyrinthe. Cette approche permet d'explorer une technique de rendu historique tout en créant une expérience de jeu unique et immersive.

L'objectif principal est d'implémenter un moteur de raycasting fonctionnel capable de transformer une carte 2D en rendu 3D temps réel, tout en conservant la logique de jeu du Tic-Tac-Toe. Le projet est développé en Go et compilé en WebAssembly pour être exécuté directement dans le navigateur.

= Style de code

Notre code suit des principes clairs pour garantir la lisibilité et la maintenabilité :

- *Guarding clauses* : Les conditions de sortie sont placées en début de fonction pour réduire l'indentation et améliorer la clarté
- *Variables explicites* : Les noms de variables sont descriptifs et explique leur intention simplement
- *Fonctions à responsabilité unique* : Chaque fonction effectue une seule tâche bien définie
- *Commentaires pertinents* : Les commentaires expliquent le code qui n'est pas immédiatement évident

Pour assurer la qualité du code, nous utilisons `golangci-lint` avec une configuration stricte basée sur un Gist public et connu (475 stars le 11.01.26). Voici le lien du Gist #link("https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322")[Golden config for golangci-lint
]. Les règles principales activées sont :

- Les erreurs non gérées (`errcheck`)
- La complexité (`gocognit`)
- La complexité cyclomatique (`cyclop`, `gocyclo`)
- Les nombres magiques (`mnd`)
- Les variables globales (`gochecknoglobals`)
- Le respect des conventions Go (`revive`)
- La documentation (`godoclint`)

Cette configuration est très stricte mais nous aide à maintenir un code propre et cohérent tout au long du projet.

= Raycasting : Principe et Fonctionnement

Le raycasting est une technique qui crée l'illusion de la 3D en projetant des rayons depuis la position du joueur. Pour chaque colonne verticale de l'écran, un rayon est lancé dans la direction de vue du joueur. La hauteur de chaque colonne affichée dépend de la distance calculée, plus l'objet est proche, plus la colonne est haute.
Cette technique était très populaire pour d'ancien jeux comme Wolfenstein 3D et Doom, car elle permettait de simuler un environnement 3D avec des ressources limitées. Actuellement cette méthode n'est plus du tout utilisé car des nouvelles techniques sont apparus qui sont beaucoup plus efficace mais extrement plus complexe à implémenter. Cette méthode est très interessante car la complexité pour la comprendre et assez faible et donc l'implémentation est relativement simple même si le code est assez conséquent.

Dans tous les scène 3D à la première personne, il y'a un paramètre nommé le FOV (Field Of View) qui détermine l'angle de vision du joueur. Un FOV plus large permet de voir plus de l'environnement mais peut aussi déformer la perspective. En général un FOV entre 60 et 90 degrés est utilisé pour un rendu naturel. Dans notre projet nous avons choisi un FOV de 1.58 radians (environ 90 degrés), cela donne un bon rendu visuel. Le fov nous permet de calculer la direction de chaque rayon à lancer en fonction de la position du joueur et de l'angle de vue actuel.

Dans l'ordre du rendu, nous commencons simplemet par afficher le sol et le plafond en deux rectangles avec des couleurs unis mais différentess. Ensuite nous lançons un rayon pour chaque colonne de l'écran, calculons la distance au mur le plus proche, déterminons la colonne de texture à afficher, et dessinons la colonne texturée. Enfin nous gérons les sprites (objets 3D) en les triant par distance et en les dessinant correctement avec le Z-buffer. Le z-buffer est un tableau qui stocke la distance de chaque colonne de l'écran au mur le plus proche, cela permet de savoir si un sprite doit être dessiné devant ou derrière un mur.

== Algorithme DDA (Digital Differential Analyzer)

Le but du raycasting est simplement d'envoyer un rayon dans une direction et de savoir à quelle distance il va toucher un mur. Par logique on pourrait juste avancer le long du rayon par petits incréments et vérifier à chaque étape si on a touché un mur, cette méthode ce nomme le *Ray Marching*. Cependant cette méthode est inefficace car elle nécessite beaucoup d'itérations pour chaque rayon, surtout si les murs sont éloignés et choisir un incrément trop grand peut faire rater des murs.

L'algorithme *DDA* (Digital Differential Analyzer) corrige ce problème en parcourant la grille de manière efficace en sautant de cellule en cellule. Les étapes principales pour chaque rayon sont :

+ *Initialisation* : On calcule pour chaque direction (X et Y) la distance à parcourir pour atteindre la prochaine ligne de grille
+ *Itération* : À chaque étape, on avance vers la ligne de grille la plus proche (verticale ou horizontale)
+ *Détection* : Dès qu'on entre dans une cellule contenant un mur, on s'arrête et on retourne la distance


La figure @fig:raycasting-ray-marching-and-dda illustre la différence entre le ray marching et l'algorithme DDA optimisé. Visuellement l'on voit assez rapidement l'avantage du DDA qui effectue beaucoup moins d'itérations pour atteindre le mur, 8 iteration pour le DDA et 15 pour le Ray marching. L'effet peux sembler négligeable, cependant dans l'exemple il y'a seulement 1 rayon, mais si nous considérons un écran de 1280 pixels de large, cela fait 1280 rayons à lancer à chaque frame. Avec le Ray Marching cela ferait 19200 itérations contre seulement 10240 pour le DDA, soit presque le double d'itérations. Ici c'est un exemple simple, en général le dda est beaucoup plus efficace que le ray marching, surtout dans des environnements complexes avec beaucoup de murs éloignés.

#figure(
  image("./assets/images/ray_marching_and_dda.svg", width: 80%),
  caption: ["Ray Marching vs DDA"],
) <fig:raycasting-ray-marching-and-dda>

L'implémentation de l'algorithme DDA est guidée par les explications présentes sur le site #link("https://lodev.org/cgtutor/raycasting.html")[Lode's Computer Graphics Tutorial - Raycasting].

== Sélection de la texture

Par défault dans le raycasting nous ne sommes pas obligé d'afficher des textures sur les murs, on peut très bien afficher des murs unis. Cependant pour rendre le rendu plus intéressant visuellement nous avons décidé d'ajouter des textures.

Une fois le mur détecté, on détermine quelle colonne de texture afficher :

+ On récupère la position exacte où le rayon a touché le mur (`wallX`), une valeur entre 0 et 1
+ On multiplie `wallX` par la largeur de la texture pour obtenir la colonne
+ Cette colonne est étirée verticalement selon la distance pour créer la perspective


La figure @fig:raycasting-wallx-calculation affiche un rayon et les coordonées d'impact sur le mur. En fonction de l'orientation du mur (vertical ou horizontal) nous utilisons soit la coordonnée Y (pour les murs verticaux) soit la coordonnée X (pour les murs horizontaux) pour calculer la coordonnée `wallX`. Dans l'exemple le rayon touche un mur vertical, nous utilisons donc la coordonnée Y de l'impact (ici 6.3) et nous soustrayons la partie entière (6) pour obtenir la partie décimale (0.3). Cette valeur est ensuite multipliée par la largeur de la texture, par exemple 64 pixels, pour obtenir la colonne de texture à afficher ($0.3 * 64 = 19.2$) et donc la collone 19.

#figure(
  image("./assets/images/texture_wall_x.svg", width: 40%),
  caption: ["Calcul de la coordonnée wallX pour la sélection de la texture"],
) <fig:raycasting-wallx-calculation>


= Implémentation du moteur de Raycasting

Le moteur de raycasting est implémenté dans le fichier `raycaster.go` et est utilisé par le fichier `world.go` pour afficher le résultat à l'écran. Le but de ce moteur est de lancer un rayon pour chaque colonne de l'écran, calculer la distance au mur le plus proche, déterminer la colonne de texture à afficher. Le moteur n'affiche pas directement à l'écran, il retourne les informations nécessaires à `world.go` qui se charge du rendu.

Le code du raycasting se trouve dans `raycaster.go`. Voici les fonctions principales :


TODO: Changer pour afficher la nouvelle structure
```go
// RayHit represents the result of casting a ray in the raycasting engine.
// hit indicates whether a wall was hit.
// cellX and cellY are the grid coordinates of the hit cell.
// distance is the distance from the player to the wall.
// wallX is the exact position along the wall where the ray hit (between 0 and 1).
// side indicates whether a vertical (0) or horizontal (1) wall was hit.
type RayHit struct {
    hit      bool
    cellX    int
    cellY    int
    distance float64
    wallX    float64
    side     int
}
```


Cette fonction implémente l'algorithme DDA complet :

```go
func CastRay(
    playerPosition Vec2,
    rayDirection Vec2,
    grid Grid,
    maxIterations int,
) RayHit {
    // Cellule de départ
    mapX := int(playerPosition.X)
    mapY := int(playerPosition.Y)

    // Calcul des deltas : distance pour traverser une cellule
    deltaDistX := math.Abs(1 / rayDirection.X)
    deltaDistY := math.Abs(1 / rayDirection.Y)

    // Détermination de la direction et distance initiale
    var stepX, stepY int
    var sideDistX, sideDistY float64

    if rayDirection.X < 0 {
        stepX = -1
        sideDistX = (playerPosition.X - float64(mapX)) * deltaDistX
    } else {
        stepX = 1
        sideDistX = (float64(mapX) + 1.0 - playerPosition.X) * deltaDistX
    }

    // Boucle DDA : avancer jusqu'à toucher un mur
    for range maxIterations {
        if sideDistX < sideDistY {
            sideDistX += deltaDistX
            mapX += stepX
            side = 0
        } else {
            sideDistY += deltaDistY
            mapY += stepY
            side = 1
        }

        if isGridCellNotEmpty(grid, mapX, mapY) {
            hit = true
            break
        }
    }

    // Calcul de la distance perpendiculaire
    if side == 0 {
        distance = sideDistX - deltaDistX
    } else {
        distance = sideDistY - deltaDistY
    }

    return RayHit{...}
}
```

Le fichier `world.go` orchestre le rendu complet en appelant `CastRay` pour chaque colonne de l'écran :

```go
func (w *World) raycastColumnsAndDrawWalls(screen *ebiten.Image, g *Game, p *Player) {
    for x := 0; x < WindowSizeX; x++ {
        // Lancer un rayon pour cette colonne
        hit, ok := w.castRayForScreenColumn(g, p, x)
        if !ok {
            continue
        }

        // Stocker la distance dans le Z-buffer
        w.zBuffer[x] = hit.distance

        // Récupérer la colonne de texture correspondante
        strip, ok := w.resolveTextureStripFromHit(g, hit)

        // Calculer la hauteur du mur à l'écran
        lineH := w.wallSliceHeightOnScreen(hit.distance)

        // Dessiner la colonne texturée
        w.drawTexturedWallSlice(screen, strip, x, drawStart, lineH, hit.distance)
    }
}
```

= Gestion des Textures et Design

Toutes les textures du projet ont été générées avec ChatGPT, ce qui s'est révélé étonnamment puissant et pratique. ChatGPT ne permet de générer que des images en haute résolution, mais il suffit de les redimensionner à 64x64 pixels dans Figma pour obtenir le format requis et les couleurs voulues. Nous avons décider d'adopter un style pixel art et plutot sombre pour coller à l'ambiance d'un donjon. Les textures ont une taille fixe de 64x64 pixels, ce qui créé cet effet pixelisé et permet aussi d'améliorer les performances en réduisant la quantité de données à traiter.

Le design system complet est disponible sur #link("https://www.figma.com/design/I9ar6LeUqDPe7zEKGo3DGS/Ray-Casting?node-id=0-1&t=Cq7NpTd7JF853whj-1")[Figma - Ray Casting Design System].

Le raycasting necessite toujours de dessiner des rectangles de 1 pixel de large pour chaque colonne de l'écran. Pour optimiser le rendu, lors du chargement initiales du jeu , toutes les textures sont prédécoupées en bandes verticales d'un pixel de large et stockées dans une structure `TextureMap` pour un accès rapide lors du rendu. Cela nous permet de pouvoir directement accéder à la bande verticale nécessaire sans avoir à découper l'image à chaque frame, ce qui serait très coûteux en performance.

TODO: Affiche le type Texture et TextureMap avant ce code
```go
//go:embed assets/textures/*.png
var texturesFS embed.FS

func LoadTextures() (TextureMap, error) {
    for id, filename := range imageManifest {
        img, _, err := ebitenutil.NewImageFromReader(f)

        strips, err := sliceIntoVerticalStrips(img)

        out[id] = Texture{
            Source: img,
            Strips: strips,
        }
    }
}
```

Chaque texture est prédécoupée en 64 bandes verticales d'un pixel de large, permettant un rendu extrêmement rapide lors du raycasting.

= Tests et Validation

TODO: Décrire les tests unitaires et d'intégration mis en place pour valider le moteur de raycasting et la logique du jeu.


= Pipeline CI/CD

Le projet utilise des GitHub Actions pour automatiser le workflow de développement et de déploiement. Deux pipelines sont configurées :

- *CI* (`.github/workflows/ci.yml`) : Linting du code avec `golangci-lint` à chaque push et pull request
- *CD* (`.github/workflows/cd.yml`) : Compilation en WebAssembly et déploiement sur GitHub Pages après une CI réussie


= Résultat Final

Le résultat correspond exactement à notre vision, un jeu de Tic-Tac-Toe normal, mais présenté dans un environnement 3D. Le joueur peut se déplacer dans une sorte labyrinthe où chaque salle représente une case du jeu. Les joueurs peuvent se déplacer avec les touches *WASD* et placer leurs symboles (X ou O) dans les cases en appuyant sur la touche *E*. Nous sommes vraiment fiers du résultat final, qui combine une technique de rendu classique avec une logique de jeu simple pour créer une expérience unique et immersive.

Le jeu est jouable directement dans le navigateur à l'adresse : #link("https://yungbricocoop.github.io/GoTicTacToe/")[yungbricocoop.github.io/GoTicTacToe]


= Apprentissages et améliorations

Ce projet nous a beaucoup appris sur plusieurs aspects du développement logiciel.

La collaboration sur un projet d'une certaine ampleur nécessite une bonne entente et une communication claire. Nous avons dû synchroniser régulièrement nos efforts et nous assurer que les modifications de chacun s'intègrent correctement. L'utilisation de Git avec des branches séparées et des pull requests a été essentielle pour maintenir un code stable.

Le raycasting est conceptuellement simple : lancer des rayons et calculer des distances. Cependant, l'implémentation révèle de nombreux détails subtils auxquels il faut faire attention : gérer correctement les divisions par zéro, éviter l'effet fish-eye en utilisant la distance perpendiculaire, gérer les bords de la carte. Écrire un code clair et compréhensible pour cet algorithme s'est révélé difficile, nécessitant plusieurs refactorings pour améliorer la structure.

Nous avons notamment dû extraire des fonctions comme `isGridCellNotEmpty`, `isRayOutOfBounds` pour rendre la boucle DDA plus lisible. Le découpage du rendu en étapes distinctes (raycasting des murs, puis rendu des sprites) a également clarifié le code.

La gestion du Z-buffer pour afficher correctement les sprites derrière les murs a demandé une réflexion particulière sur l'ordre de rendu et le stockage des distances.

Le raycasting est une technique simple en théorie, mais son implémentation est assez complexe, surtout la partie affichage au final. Le code require beaucoup de petit calcule et le code vient très rapidement illisible, nous avons donc dû faire beaucoup d'efforts pour garder un code propre et compréhensible mais nous pensons qu'il est encore perfectible. Dans les améliorations futures, nous pourrions envisager d'ajouter des fonctionnalités comme des portes ouvrables, des ennemis basiques qui se déplacent dans le labyrinthe, ou même un mode multijoueur local. Avec l'état actuel du code ces ajouts ne sont pas très compliqués à implémenter.

= Conclusion sur Go

Go est un langage génial pour ce type de projet. Sa syntaxe est simple et lisible, ce qui permet d'avancer rapidement dans le développement sans se perdre dans des abstractions complexes. La courbe d'apprentissage est plutot douce et donc permet à presque n'importe quel developpeur de se lancer rapidement.

La cross-compilation est vraiment ultra pratique, surtout pour WebAssembly. Avec une simple commande `GOOS=js GOARCH=wasm go build`, le projet est compilé pour le web. Go fournit directement `wasm_exec.js` dans son installation, simplifiant encore plus le processus.

Le système de modules et la gestion des dépendances et vraiment très simple. Nous n'avons eu aucun problème à intégrer et mettre à jour des bibliothèques externes comme Ebiten.

Pour conclure, c'est vraiment une très bonne découverte et c'est un language qui fait envie de travailler avec. Peut-être qu'il y'a des petits points négatifs que l'on trouve quand l'on développe des projets plus complexes, mais dans notre cas tout s'est très bien passé.


= Références

- #link("https://lodev.org/cgtutor/raycasting.html")[Lode's Computer Graphics Tutorial - Raycasting]
- #link("https://lodev.org/cgtutor/raycasting3.html")[Lode's Computer Graphics Tutorial - Raycasting 3]
- #link("https://yungbricocoop.github.io/GoTicTacToe/")[Projet en ligne]
- #link("https://www.figma.com/design/I9ar6LeUqDPe7zEKGo3DGS/Ray-Casting")[Design System Figma]
