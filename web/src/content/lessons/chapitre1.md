# Chapitre 1 — Introduction et premiers pas

## Objectifs du chapitre

À la fin de ce chapitre, vous serez capable de :

- Comprendre la syntaxe de base de Lua
- Utiliser les différents types de données
- Manipuler les variables et les opérateurs
- Écrire vos premiers programmes Lua

## 1. Introduction à Lua

Lua est un langage de programmation léger, rapide et puissant, particulièrement adapté à l'intégration dans d'autres applications. Dans le contexte de **FiveM**, Lua sert à créer des scripts côté client et serveur.

**Caractéristiques principales :**

- **Simplicité** — syntaxe claire et intuitive
- **Performance** — exécution rapide grâce à LuaJIT
- **Flexibilité** — typage dynamique
- **Intégration** — facilement embarquable dans d'autres programmes

> Sur cette plateforme, vous n'avez rien à installer : chaque exemple et chaque exercice s'exécute directement dans votre navigateur. Cliquez sur **Lancer** pour voir le résultat.

## 2. Syntaxe de base

### Premier programme

```lua
-- Ceci est un commentaire sur une ligne
print("Hello, world!")

--[[
  Ceci est un commentaire
  sur plusieurs lignes
--]]
```

### Règles importantes

- **Sensibilité à la casse** : `Variable` ≠ `variable`
- **Pas de point-virgule obligatoire** (mais possible pour la clarté)
- **Mots-clés réservés** : `and`, `break`, `do`, `else`, `elseif`, `end`, `false`, `for`, `function`, `if`, `in`, `local`, `nil`, `not`, `or`, `repeat`, `return`, `then`, `true`, `until`, `while`

## 3. Variables et conventions de nommage

**Toutes les variables doivent être nommées en anglais.** On distingue plusieurs conventions selon la portée :

| Élément | Convention | Exemple |
| --- | --- | --- |
| Variable locale | `camelCase` | `playerName`, `currentHealth` |
| Constante | `SNAKE_CASE_MAJUSCULE` | `MAX_PLAYERS`, `MAX_HEALTH` |
| Variable globale | `PascalCase` | `PlayerData`, `ServerSettings` |
| Fonction locale | `camelCase` | `calculateDistance` |
| Fonction globale | `PascalCase` | `GetPlayerMoney` |

### Déclaration

```lua
-- Variable globale (PascalCase)
PlayerName = "John"

-- Variable locale (camelCase) — RECOMMANDÉE
local playerAge = 25

-- Constante (SNAKE_CASE_MAJUSCULE)
local MAX_HEALTH = 100

-- Plusieurs déclarations à la fois
local posX, posY, posZ = 1, 2, 3
```

> **Bonne pratique :** préférez toujours `local` pour éviter de polluer l'espace global.

## 4. Les types de données

### nil

```lua
local emptyVariable = nil
print(type(emptyVariable)) -- nil
```

### boolean

```lua
local isActive = true
local isDisabled = false
```

### number

```lua
local playerLevel = 42        -- entier
local healthPercentage = 85.5 -- décimal
local largeNumber = 1.23e10   -- notation scientifique
local colorCode = 0xFF0000    -- hexadécimal
```

### string

```lua
local firstName = 'Alice'                 -- guillemets simples
local welcome = "Welcome to the server!"  -- guillemets doubles

local serverRules = [[
1. No cheating allowed
2. Respect other players
3. Have fun!
]]

local greeting = "Hello " .. firstName    -- concaténation
```

### La fonction `type()`

```lua
print(type(42))      -- number
print(type("hello")) -- string
print(type(true))    -- boolean
print(type(nil))     -- nil
```

## 5. Les opérateurs

### Arithmétiques

```lua
local a, b = 10, 3

print(a + b)  -- Addition : 13
print(a - b)  -- Soustraction : 7
print(a * b)  -- Multiplication : 30
print(a / b)  -- Division : 3.333...
print(a % b)  -- Modulo (reste) : 1
print(a ^ b)  -- Puissance : 1000
print(-a)     -- Négation : -10
```

### Comparaison

```lua
local x, y = 5, 5

print(x == y)  -- Égalité : true
print(x ~= y)  -- Différence : false   (le "différent de" s'écrit ~=)
print(x < y)   -- false
print(x <= y)  -- true
```

### Logiques

```lua
print(true and false) -- ET logique : false
print(true or false)  -- OU logique : true
print(not true)       -- NON logique : false
```

### Concaténation

```lua
local firstName = "Marie"
local lastName = "Johnson"
print(firstName .. " " .. lastName) -- Marie Johnson
```

### L'opérateur ternaire de Lua

Lua n'a pas de `? :` mais on l'imite avec `and` / `or` :

```lua
local playerAge = 20
local status = (playerAge >= 18) and "adult" or "minor"
print("Status:", status) -- Status: adult

-- Forme générale : condition and valeur_si_vrai or valeur_si_faux
```

## 6. Points clés à retenir

1. Préférez toujours `local`.
2. Nommez en anglais, en respectant les conventions (`camelCase`, `PascalCase`, `SNAKE_CASE_MAJUSCULE`).
3. Les 4 types de base : `nil`, `boolean`, `number`, `string`.
4. Le typage est dynamique : une variable peut changer de type.
5. La concaténation se fait avec `..`.
6. Commentaires : `--` pour une ligne, `--[[ ]]` pour un bloc.

---

À vous de jouer ! Passez aux exercices ci-contre pour mettre tout ça en pratique.
