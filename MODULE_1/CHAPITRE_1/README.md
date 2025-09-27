# Chapitre 1 : Introduction et Premiers Pas en Lua

## 🎯 Objectifs du Chapitre
À la fin de ce chapitre, vous serez capable de :
- Installer et configurer un environnement de développement Lua
- Comprendre la syntaxe de base de Lua
- Utiliser les différents types de données
- Manipuler les variables et opérateurs
- Écrire vos premiers programmes Lua

## 📚 1. Introduction à Lua

Lua est un langage de programmation léger, rapide et puissant, particulièrement adapté pour l'intégration dans d'autres applications. Dans le contexte de FiveM, Lua est utilisé pour créer des scripts côté client et serveur.

### Caractéristiques principales :
- **Simplicité** : Syntaxe claire et intuitive
- **Performance** : Exécution rapide grâce à LuaJIT
- **Flexibilité** : Typage dynamique
- **Intégration** : Facilement intégrable dans d'autres programmes

## 🛠️ 2. Installation et Environnement

### Option 1 : Lua Standalone
Téléchargez Lua depuis [lua.org](https://www.lua.org/download.html)

### Option 2 : Éditeurs recommandés
- **Visual Studio Code** avec l'extension "Lua"
- **Sublime Text** avec le package Lua
- **Atom** avec lua-language

### Option 3 : Environnement en ligne (pour débuter)
- [Repl.it](https://repl.it) - Section Lua
- [OnlineGDB](https://onlinegdb.com) - Support Lua

## 📝 3. Syntaxe de Base

### Premier Programme
```lua
-- This is a single line comment
print("Hello, world!")

--[[
This is a multi-line
comment block
--]]
```

### Règles de syntaxe importantes :
- **Sensibilité à la casse** : `Variable` ≠ `variable`
- **Pas de point-virgule obligatoire** (mais recommandé pour la clarté)
- **Mots-clés réservés** : `and`, `break`, `do`, `else`, `elseif`, `end`, `false`, `for`, `function`, `if`, `in`, `local`, `nil`, `not`, `or`, `repeat`, `return`, `then`, `true`, `until`, `while`

## 🏷️ 4. Variables et Conventions de Nommage

### Règles de Nommage (OBLIGATOIRES)

**Toutes les variables doivent être nommées en anglais uniquement !**

#### Variables Locales : camelCase
- Première lettre **minuscule**
- Mots suivants attachés avec première lettre **majuscule**
```lua
local playerName = "John"
local currentHealth = 100
local isPlayerAlive = true
local maxAmmoCount = 250
```

#### Constantes : SNAKE_CASE_MAJUSCULE
- Toutes les lettres en **MAJUSCULE**
- Mots séparés par des **underscores (_)**
```lua
local MAX_PLAYERS = 32
local DEFAULT_SPAWN_POSITION = vector3(0, 0, 0)
local WEAPON_DAMAGE_MULTIPLIER = 1.5
```

#### Variables Globales : PascalCase
- Première lettre **majuscule**
- Chaque mot attaché avec première lettre **majuscule**
```lua
PlayerData = {}
ServerSettings = {}
WeaponsList = {}
CurrentGameMode = "freeroam"
```

#### Fonctions Locales : camelCase
- Première lettre **minuscule**
- Première lettre des mots suivants en **majuscule**
```lua
local function calculateDistance(pos1, pos2)
    -- code here
end

local function spawnVehicle(model, coords)
    -- code here
end
```

#### Fonctions Globales : PascalCase
- Première lettre de chaque mot en **majuscule**
```lua
function GetPlayerMoney(playerId)
    -- code here
end

function SpawnWeapon(weaponHash, ammo)
    -- code here
end
```

### Déclaration de Variables
```lua
-- Variable globale (PascalCase)
PlayerName = "John"

-- Variable locale (camelCase) - RECOMMANDÉE
local playerAge = 25

-- Constante (SNAKE_CASE_MAJUSCULE)
local MAX_HEALTH = 100

-- Plusieurs déclarations
local posX, posY, posZ = 1, 2, 3
```

### Types de Données Fondamentaux

#### 4.1 Nil
```lua
local emptyVariable = nil
print(type(emptyVariable)) -- affiche : nil
```

#### 4.2 Boolean
```lua
local isActive = true
local isDisabled = false
```

#### 4.3 Number
```lua
-- Entiers
local playerLevel = 42

-- Décimaux
local healthPercentage = 85.5

-- Notation scientifique
local largeNumber = 1.23e10

-- Hexadécimal
local colorCode = 0xFF0000
```

#### 4.4 String
```lua
-- Guillemets simples
local firstName = 'Alice'

-- Guillemets doubles
local welcomeMessage = "Welcome to the server!"

-- Chaînes multiligne
local serverRules = [[
1. No cheating allowed
2. Respect other players
3. Have fun!
]]

-- Concaténation
local fullGreeting = "Hello " .. firstName
```

### Fonction `type()`
```lua
print(type(42))           -- number
print(type("hello"))      -- string
print(type(true))         -- boolean
print(type(nil))          -- nil
```

## ⚙️ 5. Opérateurs

### Opérateurs Arithmétiques
```lua
local firstNumber, secondNumber = 10, 3

print(firstNumber + secondNumber)    -- Addition : 13
print(firstNumber - secondNumber)    -- Soustraction : 7
print(firstNumber * secondNumber)    -- Multiplication : 30
print(firstNumber / secondNumber)    -- Division : 3.333...
print(firstNumber % secondNumber)    -- Modulo : 1
print(firstNumber ^ secondNumber)    -- Puissance : 1000
print(-firstNumber)                  -- Négation : -10
```

### Opérateurs de Comparaison
```lua
local valueA, valueB = 5, 5

print(valueA == valueB)   -- Égalité : true
print(valueA ~= valueB)   -- Différence : false
print(valueA < valueB)    -- Inférieur : false
print(valueA <= valueB)   -- Inférieur ou égal : true
print(valueA > valueB)    -- Supérieur : false
print(valueA >= valueB)   -- Supérieur ou égal : true
```

### Opérateurs Logiques
```lua
local isTrue, isFalse = true, false

print(isTrue and isFalse)    -- ET logique : false
print(isTrue or isFalse)     -- OU logique : true
print(not isTrue)            -- NON logique : false
```

### Opérateur de Concaténation
```lua
local firstName = "Marie"
local lastName = "Johnson"
local fullName = firstName .. " " .. lastName
print(fullName)  -- Marie Johnson
```

### Opérateur Ternaire (Bonus)
Lua possède une forme d'opérateur ternaire utilisant `and` et `or` :
```lua
local playerAge = 20
local playerStatus = (playerAge >= 18) and "adult" or "minor"
print("Status:", playerStatus)  -- Status: adult

-- Équivaut à : condition and valeur_si_vrai or valeur_si_faux
local examScore = 15
local examResult = (examScore >= 10) and "Passed" or "Failed"
print(examResult)  -- Passed
```

## 🎯 6. Exercices Pratiques

### Exercice 1 : Variables et Types (Facile)
```lua
-- À compléter (respectez les conventions de nommage) :
-- 1. Créez une variable locale contenant votre prénom (en anglais)
-- 2. Créez une variable contenant votre âge
-- 3. Créez une variable booléenne indiquant si vous aimez programmer
-- 4. Créez une constante pour le nombre maximum de joueurs (32)
-- 5. Affichez le type de chaque variable

-- Votre code ici dans exercices/exercice1.lua:
```

### Exercice 2 : Calculs Arithmétiques (Facile)
```lua
-- Créez un petit programme qui :
-- 1. Définit deux nombres (utilisez les bonnes conventions)
-- 2. Calcule et affiche leur somme, différence, produit et division
-- 3. Calcule le reste de la division du premier par le second

-- Votre code ici dans exercices/exercice2.lua:
```

### Exercice 3 : Manipulation de Chaînes (Moyen)
```lua
-- Créez un programme qui :
-- 1. Définit votre prénom et nom dans des variables séparées (en anglais)
-- 2. Crée votre nom complet par concaténation
-- 3. Affiche "Welcome [nom complet] to Lua programming!"

-- Votre code ici dans exercices/exercice3.lua:
```

### Exercice 4 : Opérateurs Logiques (Moyen)
```lua
-- Créez un programme qui :
-- 1. Définit trois variables booléennes : isAdult, hasLicense, hasVehicle
-- 2. Utilise les opérateurs logiques pour déterminer si une personne peut conduire
-- 3. Affiche le résultat avec une phrase explicative

-- Une personne peut conduire si elle est majeure ET qu'elle a le permis ET une voiture

-- Votre code ici dans exercices/exercice4.lua:
```

### Exercice 5 : Défis (Difficile)
```lua
-- 1. Échangez les valeurs de deux variables sans utiliser de variable temporaire
-- 2. Calculez l'aire et le périmètre d'un rectangle
-- 3. Déterminez si un nombre est pair en utilisant l'opérateur modulo
-- 4. Créez une phrase en utilisant plusieurs concaténations

-- Votre code ici dans exercices/exercice5.lua:
```

## 📋 7. Points Clés à Retenir

1. **Variables locales** : Toujours préférer `local` pour éviter la pollution globale
2. **Nommage en anglais** : Obligatoire pour toutes les variables et fonctions
3. **Conventions de nommage** :
   - Variables locales : `camelCase` (ex: `playerHealth`)
   - Constantes : `SNAKE_CASE_MAJUSCULE` (ex: `MAX_PLAYERS`)
   - Variables globales : `PascalCase` (ex: `PlayerData`)
   - Fonctions locales : `camelCase` (ex: `calculateDistance`)
   - Fonctions globales : `PascalCase` (ex: `GetPlayerMoney`)
4. **Types dynamiques** : Une variable peut changer de type durant l'exécution
5. **Nil** : Valeur par défaut des variables non initialisées
6. **Concaténation** : Utilisez `..` pour joindre des chaînes
7. **Commentaires** : `--` pour une ligne, `--[[ ]]` pour plusieurs lignes

## 🏁 8. Validation des Acquis

Avant de passer au chapitre suivant, assurez-vous de maîtriser :
- ✅ Les conventions de nommage en anglais (camelCase, PascalCase, SNAKE_CASE_MAJUSCULE)
- ✅ La déclaration de variables locales et globales
- ✅ Les 4 types de données de base (nil, boolean, number, string)
- ✅ L'utilisation des opérateurs arithmétiques, logiques et de comparaison
- ✅ La concaténation de chaînes
- ✅ L'utilisation de la fonction `print()` et `type()`
- ✅ L'opérateur ternaire de Lua avec `and` et `or`

## 📖 Prochaine Étape

Dans le **Chapitre 2**, nous découvrirons les structures de contrôle (conditions et boucles) qui vous permettront de créer des programmes plus dynamiques et interactifs.

---

*💡 Conseil : Pratiquez chaque exercice plusieurs fois avec des valeurs différentes pour bien assimiler les concepts !*