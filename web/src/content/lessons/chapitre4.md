# Chapitre 4 — Les fonctions

## Objectifs du chapitre

À la fin de ce chapitre, vous serez capable de :

- Définir et appeler des fonctions
- Passer des paramètres et renvoyer une ou plusieurs valeurs
- Gérer des arguments par défaut
- Écrire des fonctions anonymes et les passer en argument (callbacks)
- Comprendre et utiliser les **closures**
- Écrire des fonctions **variadiques** (`...`)

## 1. Définir et appeler une fonction

Une fonction regroupe un bloc de code réutilisable. Comme pour les variables, on respecte
les conventions : `camelCase` pour les fonctions **locales**, `PascalCase` pour les
fonctions **globales**.

```lua
-- Fonction locale (recommandée)
local function sayHello()
    print("Hello!")
end

sayHello() -- Hello!

-- Fonction globale
function GetServerName()
    return "Los Santos RP"
end

print(GetServerName()) -- Los Santos RP
```

On peut aussi stocker une fonction dans une variable (les fonctions sont des **valeurs**) :

```lua
local greet = function()
    print("Hi there!")
end

greet()
```

## 2. Paramètres et arguments

Les **paramètres** sont les variables déclarées entre parenthèses ; les **arguments**
sont les valeurs réellement passées lors de l'appel.

```lua
local function greet(name)
    print("Hello " .. name .. "!")
end

greet("Alex")  -- Hello Alex!
greet("Sam")   -- Hello Sam!
```

Plusieurs paramètres se séparent par des virgules :

```lua
local function add(a, b)
    print(a + b)
end

add(3, 5) -- 8
```

> Si vous appelez une fonction avec **moins** d'arguments que de paramètres, les
> paramètres manquants valent `nil`. S'il y en a **trop**, les excédentaires sont ignorés.

## 3. Valeurs de retour

Le mot-clé `return` renvoie une valeur à l'appelant :

```lua
local function square(n)
    return n * n
end

local result = square(6)
print(result) -- 36
```

### Retours multiples

Particularité de Lua : une fonction peut renvoyer **plusieurs valeurs** à la fois.

```lua
local function minMax(a, b)
    if a <= b then
        return a, b
    else
        return b, a
    end
end

local low, high = minMax(8, 3)
print(low, high) -- 3   8
```

## 4. Arguments par défaut

Lua n'a pas de syntaxe dédiée, mais l'idiome `valeur or défaut` fait le travail :

```lua
local function createPlayer(name, health)
    health = health or 100  -- si health est nil, on utilise 100
    return { name = name, health = health }
end

local p1 = createPlayer("Alex")       -- health = 100
local p2 = createPlayer("Sam", 50)    -- health = 50
print(p1.health, p2.health)           -- 100   50
```

> ⚠️ Attention : `x or défaut` remplace aussi `false`. Pour un booléen optionnel,
> testez explicitement `if x == nil then`.

## 5. Fonctions anonymes et callbacks

Une fonction sans nom est dite **anonyme**. Comme les fonctions sont des valeurs, on peut
les passer en argument d'une autre fonction : c'est le principe des **callbacks**, omniprésent
dans FiveM (gestion d'événements).

```lua
local function repeatTimes(count, action)
    for i = 1, count do
        action(i)
    end
end

-- On passe une fonction anonyme comme argument
repeatTimes(3, function(i)
    print("Tick " .. i)
end)
-- Tick 1
-- Tick 2
-- Tick 3
```

On peut aussi stocker des fonctions dans une table (utile pour un routeur d'événements) :

```lua
local handlers = {
    onJoin = function(name) print(name .. " joined") end,
    onLeave = function(name) print(name .. " left") end,
}

handlers.onJoin("Alex")  -- Alex joined
```

## 6. Les closures

Une **closure** est une fonction qui « capture » les variables locales de l'environnement
où elle a été créée — et garde accès à ces variables même après.

```lua
local function makeCounter()
    local count = 0          -- variable capturée
    return function()
        count = count + 1
        return count
    end
end

local next = makeCounter()
print(next()) -- 1
print(next()) -- 2
print(next()) -- 3

-- Chaque compteur a son propre état :
local other = makeCounter()
print(other()) -- 1
```

Les closures servent à créer de l'**état privé** sans variable globale.

## 7. Fonctions variadiques (`...`)

Le symbole `...` permet d'accepter un nombre **variable** d'arguments.

```lua
local function printAll(...)
    for i = 1, select('#', ...) do
        print(i, (select(i, ...)))
    end
end

printAll("a", "b", "c")
-- 1   a
-- 2   b
-- 3   c
```

- `select('#', ...)` renvoie le **nombre** d'arguments,
- `select(i, ...)` renvoie les arguments à partir de la position `i`,
- `{...}` regroupe les arguments dans une table.

```lua
local function sum(...)
    local total = 0
    for i = 1, select('#', ...) do
        total = total + select(i, ...)
    end
    return total
end

print(sum(1, 2, 3, 4)) -- 10
print(sum())           -- 0
```

## 8. Récursivité

Une fonction peut s'appeler elle-même. Pensez toujours à un **cas d'arrêt** !

```lua
local function factorial(n)
    if n <= 1 then
        return 1            -- cas d'arrêt
    end
    return n * factorial(n - 1)
end

print(factorial(5)) -- 120
```

## 9. Bonnes pratiques

- Une fonction = **une responsabilité** ; un nom clair (verbe d'action).
- Préférez `local function` ; ne passez en global que ce qui doit l'être.
- Renvoyez tôt (`return`) pour éviter les `if` profondément imbriqués.
- Documentez les paramètres attendus si la fonction est exposée.

## 10. Points clés à retenir

1. `local function` (camelCase) vs `function` globale (PascalCase).
2. Les fonctions sont des **valeurs** : stockables, passables en argument.
3. `return` peut renvoyer **plusieurs valeurs**.
4. Arguments par défaut via l'idiome `x = x or défaut`.
5. Les **callbacks** = fonctions passées en argument.
6. Une **closure** capture les variables locales de son environnement.
7. `...` + `select('#', ...)` pour les fonctions variadiques.
8. La récursivité exige un cas d'arrêt.

---

Passez aux exercices : définition, retours multiples, variadiques et closures.
