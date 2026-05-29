# Chapitre 2 — Structures de contrôle

## Objectifs du chapitre

À la fin de ce chapitre, vous serez capable de :

- Utiliser les conditions pour contrôler le flux d'exécution
- Implémenter les différents types de boucles
- Maîtriser `break` et `return`
- Construire une logique conditionnelle lisible

## 1. Les conditions

### `if` simple

```lua
local playerHealth = 50

if playerHealth <= 0 then
    print("Player is dead!")
end
```

### `if / else`

```lua
local playerAge = 20

if playerAge >= 18 then
    print("Player is an adult")
else
    print("Player is a minor")
end
```

### `if / elseif / else`

```lua
local examScore = 85

if examScore >= 90 then
    print("Grade: Excellent")
elseif examScore >= 70 then
    print("Grade: Good")
elseif examScore >= 50 then
    print("Grade: Average")
else
    print("Grade: Failed")
end
```

### Conditions multiples

```lua
local playerLevel = 25
local playerGold = 1500
local hasPermission = true

if playerLevel >= 20 and playerGold >= 1000 then
    print("Can buy premium weapon")
end

if playerLevel >= 50 or hasPermission then
    print("Can access VIP area")
end

if not hasPermission then
    print("Access denied")
end
```

### Valeurs « falsy » et « truthy »

En Lua, **seules `false` et `nil` sont fausses**. Tout le reste est vrai !

```lua
if 0 then print("0 is true!") end            -- VRAI (différent du C/JS)
if "" then print("Empty string is true!") end -- VRAI
if {} then print("Empty table is true!") end  -- VRAI
```

## 2. Les boucles

### Boucle `for` numérique

```lua
-- for variable = début, fin, pas do   (le pas est optionnel, défaut 1)

for i = 1, 5 do
    print("Iteration:", i)
end

for i = 0, 10, 2 do
    print("Even number:", i) -- 0, 2, 4, 6, 8, 10
end

for i = 10, 1, -1 do
    print("Countdown:", i)
end
```

### Boucle `for` générique (itérateurs)

```lua
local playerNames = {"Alice", "Bob", "Charlie", "Diana"}

for index, name in ipairs(playerNames) do
    print(index, name)
end
```

### Boucle `while`

```lua
local counter = 0
while counter < 5 do
    counter = counter + 1
    print("Counter:", counter) -- pensez TOUJOURS à faire évoluer la condition
end
```

### Boucle `repeat / until`

```lua
-- S'exécute AU MOINS UNE FOIS
local i = 0
repeat
    i = i + 1
    print("i =", i)
until i >= 3
```

## 3. Instructions de contrôle

### `break`

```lua
for i = 1, 100 do
    if i == 42 then
        print("Target found at:", i)
        break -- sort de la boucle immédiatement
    end
end
```

### `return`

```lua
function checkAccess(level, required)
    if level < required then
        print("Access denied")
        return false
    end
    print("Access granted")
    return true
end
```

## 4. Bonnes pratiques

- **Évitez les boucles infinies** : assurez-vous que la condition peut devenir fausse.
- **Préférez `for` à `while`** quand le nombre d'itérations est connu.
- **Découpez les conditions complexes** en variables intermédiaires nommées :

```lua
local canAccessPremium = playerLevel >= 20 and playerGold >= 1000
local hasVipAccess = isVip and playerLevel >= 10

if canAccessPremium or hasVipAccess or isAdmin then
    -- ...
end
```

## 5. Points clés à retenir

1. `if / elseif / else` pour les conditions.
2. `and`, `or`, `not` pour combiner.
3. Seules `false` et `nil` sont « falsy ».
4. `for` numérique, `for` générique, `while`, `repeat/until`.
5. `break` sort d'une boucle, `return` sort d'une fonction.

---

Passez aux exercices pour pratiquer conditions et boucles.
