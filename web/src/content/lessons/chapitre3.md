# Chapitre 3 — Les tables, le cœur de Lua

## Objectifs du chapitre

À la fin de ce chapitre, vous serez capable de :

- Créer et manipuler des tables
- Utiliser les tables comme **tableaux** (listes indexées)
- Utiliser les tables comme **dictionnaires / objets** (clé → valeur)
- Parcourir une table avec `ipairs` et `pairs`
- Utiliser la bibliothèque `table` (`insert`, `remove`, `sort`, `concat`)
- Comprendre les bases des **métatables** et de `__index`

## 1. Pourquoi les tables ?

En Lua, la **table** est la **seule** structure de données composée. Tout ce qui contient
« plusieurs choses » est une table : une liste de joueurs, l'inventaire d'un personnage,
la configuration d'une ressource FiveM… Une table peut servir à la fois de **tableau**
(éléments numérotés) et de **dictionnaire** (paires clé → valeur).

```lua
local emptyTable = {}
print(type(emptyTable)) -- table
```

## 2. Les tables comme tableaux (arrays)

### Création et accès

> ⚠️ En Lua, les tableaux commencent à l'indice **1**, pas à 0 !

```lua
local weapons = {"pistol", "rifle", "knife"}

print(weapons[1]) -- pistol
print(weapons[2]) -- rifle
print(weapons[3]) -- knife
```

### Longueur d'un tableau

L'opérateur `#` donne le nombre d'éléments :

```lua
local weapons = {"pistol", "rifle", "knife"}
print(#weapons) -- 3
```

### Ajouter et retirer des éléments

```lua
local players = {"Alice", "Bob"}

table.insert(players, "Charlie")      -- ajoute à la fin
print(#players)                       -- 3

table.insert(players, 1, "Zoe")       -- insère en position 1
print(players[1])                     -- Zoe

table.remove(players)                 -- retire le dernier
table.remove(players, 1)              -- retire le premier
```

### Parcourir un tableau avec `ipairs`

`ipairs` parcourt les éléments **dans l'ordre**, de l'indice 1 jusqu'au premier trou :

```lua
local fruits = {"apple", "banana", "cherry"}

for index, fruit in ipairs(fruits) do
    print(index, fruit)
end
-- 1   apple
-- 2   banana
-- 3   cherry
```

### Trier un tableau

```lua
local scores = {42, 7, 100, 23}
table.sort(scores)                 -- ordre croissant
print(table.concat(scores, ", "))  -- 7, 23, 42, 100

table.sort(scores, function(a, b) return a > b end) -- décroissant
print(table.concat(scores, ", "))  -- 100, 42, 23, 7
```

## 3. Les tables comme dictionnaires (objets)

Une table peut associer des **clés** (souvent des chaînes) à des **valeurs**. C'est ainsi
qu'on représente un « objet » comme un joueur.

```lua
local player = {
    name = "Alex",
    level = 10,
    money = 500,
    isAdmin = false,
}

-- Deux notations équivalentes pour accéder à une clé :
print(player.name)      -- notation point : Alex
print(player["level"])  -- notation crochets : 10
```

### Modifier, ajouter, supprimer une clé

```lua
local player = { name = "Alex", money = 500 }

player.money = player.money + 250   -- modifier
player.job = "police"               -- ajouter une nouvelle clé
player.name = nil                   -- supprimer une clé (mettre à nil)

print(player.money) -- 750
print(player.job)   -- police
print(player.name)  -- nil
```

### Parcourir un dictionnaire avec `pairs`

`pairs` parcourt **toutes** les paires clé/valeur (l'ordre n'est pas garanti) :

```lua
local stats = { health = 100, armor = 50, stamina = 75 }

for key, value in pairs(stats) do
    print(key .. " = " .. value)
end
```

> **À retenir :** `ipairs` pour les tableaux (indices 1, 2, 3…), `pairs` pour les
> dictionnaires (n'importe quelles clés).

## 4. Tables imbriquées

Une valeur de table peut elle-même être une table. C'est la base de toute donnée
structurée dans FiveM.

```lua
local players = {
    {
        name = "Alex",
        position = { x = 100.5, y = 200.0, z = 30.0 },
        inventory = { "phone", "water", "bread" },
    },
    {
        name = "Sam",
        position = { x = -50.0, y = 10.0, z = 28.5 },
        inventory = { "radio" },
    },
}

print(players[1].name)            -- Alex
print(players[1].position.x)      -- 100.5
print(players[1].inventory[2])    -- water
print(#players[2].inventory)      -- 1
```

## 5. La bibliothèque `table`

| Fonction | Rôle |
| --- | --- |
| `table.insert(t, v)` | ajoute `v` à la fin |
| `table.insert(t, pos, v)` | insère `v` à la position `pos` |
| `table.remove(t)` | retire et renvoie le dernier élément |
| `table.remove(t, pos)` | retire l'élément à `pos` |
| `table.sort(t [, comp])` | trie sur place |
| `table.concat(t, sep)` | concatène les éléments en une chaîne |

```lua
local items = {"bread", "water", "phone"}
print(table.concat(items, " | ")) -- bread | water | phone
```

## 6. Introduction aux métatables

Une **métatable** permet de définir un comportement spécial pour une table. La
métaméthode la plus courante est `__index` : elle est consultée quand on accède à une
**clé absente**.

```lua
local defaults = { health = 100, armor = 0 }

local player = setmetatable({ health = 80 }, { __index = defaults })

print(player.health) -- 80  (présent dans player)
print(player.armor)  -- 0   (absent -> récupéré dans defaults via __index)
```

D'autres métaméthodes existent (`__newindex`, `__add`, `__call`, `__tostring`…). Elles
permettent par exemple de surcharger les opérateurs :

```lua
local Vector = {}
Vector.__add = function(a, b)
    return setmetatable({ x = a.x + b.x, y = a.y + b.y }, Vector)
end

local v1 = setmetatable({ x = 1, y = 2 }, Vector)
local v2 = setmetatable({ x = 3, y = 4 }, Vector)
local sum = v1 + v2
print(sum.x, sum.y) -- 4   6
```

> Les métatables sont la base de la **programmation orientée objet** en Lua, que nous
> approfondirons plus tard.

## 7. Bonnes pratiques

- Préférez `ipairs` pour les listes, `pairs` pour les dictionnaires.
- Souvenez-vous que `#` n'est fiable que sur des tableaux **sans trou** (indices continus).
- Pour « vider » une clé, affectez `nil` — n'utilisez pas de chaîne vide.
- Une table est passée **par référence** : la modifier dans une fonction modifie l'original.

```lua
local function addBonus(stats)
    stats.money = stats.money + 100 -- modifie la table d'origine
end

local player = { money = 500 }
addBonus(player)
print(player.money) -- 600
```

## 8. Points clés à retenir

1. La table est la seule structure de données composée de Lua.
2. Les indices de tableau commencent à **1**.
3. `#t` donne la longueur d'un tableau.
4. `table.insert` / `table.remove` pour ajouter / retirer.
5. `ipairs` pour les tableaux, `pairs` pour les dictionnaires.
6. On supprime une clé en lui affectant `nil`.
7. Les tables sont passées **par référence**.
8. `setmetatable` + `__index` fournit des valeurs par défaut et la base de la POO.

---

Passez aux exercices pour manipuler tableaux, dictionnaires et métatables.
