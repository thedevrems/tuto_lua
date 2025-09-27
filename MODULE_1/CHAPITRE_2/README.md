# Chapitre 2 : Structures de Contrôle

## 🎯 Objectifs du Chapitre
À la fin de ce chapitre, vous serez capable de :
- Utiliser les conditions pour contrôler le flux d'exécution
- Implémenter différents types de boucles
- Maîtriser les instructions de contrôle (break, return)
- Créer des programmes avec une logique conditionnelle complexe
- Optimiser vos structures de contrôle

## 📚 1. Introduction aux Structures de Contrôle

Les structures de contrôle permettent de diriger le flux d'exécution de votre programme. Elles sont essentielles pour créer des scripts dynamiques et interactifs dans FiveM.

### Types de structures :
- **Conditionnelles** : if, elseif, else
- **Itératives** : for, while, repeat
- **Contrôle** : break, return, goto (Lua 5.2+)

## 🔀 2. Structures Conditionnelles

### 2.1 Instruction `if` Simple

```lua
local playerHealth = 50

if playerHealth <= 0 then
    print("Player is dead!")
end
```

### 2.2 Structure `if-else`

```lua
local playerAge = 20

if playerAge >= 18 then
    print("Player is an adult")
else
    print("Player is a minor")
end
```

### 2.3 Structure `if-elseif-else`

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

### 2.4 Conditions Multiples

```lua
local playerLevel = 25
local playerGold = 1500
local hasPermission = true

-- Opérateur ET (and)
if playerLevel >= 20 and playerGold >= 1000 then
    print("Can buy premium weapon")
end

-- Opérateur OU (or)
if playerLevel >= 50 or hasPermission then
    print("Can access VIP area")
end

-- Opérateur NON (not)
if not hasPermission then
    print("Access denied")
end

-- Conditions complexes avec parenthèses
if (playerLevel >= 10 and playerGold >= 500) or hasPermission then
    print("Can enter special zone")
end
```

### 2.5 Valeurs "Falsy" et "Truthy"

En Lua, seules `false` et `nil` sont considérées comme "fausses". Tout le reste est "vrai" !

```lua
-- Toutes ces conditions sont VRAIES :
if 0 then print("0 is true!") end           -- VRAI (différent de C/JavaScript)
if "" then print("Empty string is true!") end -- VRAI
if {} then print("Empty table is true!") end  -- VRAI

-- Seules ces conditions sont FAUSSES :
local testValue = nil
if testValue then
    print("This won't print")
end

testValue = false
if testValue then
    print("This won't print either")
end
```

## 🔄 3. Boucles

### 3.1 Boucle `for` Numérique

```lua
-- Syntaxe : for variable = début, fin, pas do
-- Le pas est optionnel (défaut : 1)

-- Boucle simple de 1 à 5
for i = 1, 5 do
    print("Iteration:", i)
end

-- Boucle avec pas personnalisé
for i = 0, 10, 2 do
    print("Even number:", i)  -- 0, 2, 4, 6, 8, 10
end

-- Boucle décroissante
for i = 10, 1, -1 do
    print("Countdown:", i)
end

-- Exemple pratique : Afficher la santé des joueurs
for playerId = 1, 32 do
    local playerHealth = math.random(0, 100)  -- Simulation
    print("Player " .. playerId .. " health: " .. playerHealth)
end
```

### 3.2 Boucle `for` Générique (avec Itérateurs)

```lua
-- Pour les tables (nous verrons les tables en détail au chapitre 3)
local playerNames = {"Alice", "Bob", "Charlie", "Diana"}

-- ipairs pour les tableaux indexés numériquement
for index, name in ipairs(playerNames) do
    print(index, name)
end

-- Résultat :
-- 1    Alice
-- 2    Bob
-- 3    Charlie
-- 4    Diana
```

### 3.3 Boucle `while`

```lua
local attempts = 0
local maxAttempts = 5
local success = false

while attempts < maxAttempts and not success do
    attempts = attempts + 1
    print("Attempt number:", attempts)
    
    -- Simulation d'un test qui peut réussir
    if math.random() > 0.7 then
        success = true
        print("Success achieved!")
    end
end

if not success then
    print("Max attempts reached, operation failed")
end
```

### 3.4 Boucle `repeat-until`

```lua
-- La boucle repeat-until s'exécute AU MOINS UNE FOIS
local userInput
local validInput = false

repeat
    print("Enter a number between 1 and 10:")
    userInput = math.random(1, 15)  -- Simulation d'entrée utilisateur
    print("You entered:", userInput)
    
    if userInput >= 1 and userInput <= 10 then
        validInput = true
        print("Valid input!")
    else
        print("Invalid input, try again!")
    end
until validInput
```

## ⏹️ 4. Instructions de Contrôle

### 4.1 Instruction `break`

```lua
-- Sortir d'une boucle prématurément
print("Searching for target...")

for i = 1, 100 do
    if i == 42 then
        print("Target found at position:", i)
        break  -- Sort de la boucle immédiatement
    end
end

print("Search completed")

-- Exemple avec while
local counter = 0
while true do  -- Boucle infinie
    counter = counter + 1
    print("Counter:", counter)
    
    if counter >= 5 then
        break  -- Sort de la boucle
    end
end
```

### 4.2 Instruction `return`

```lua
-- return dans une fonction (avant-goût du chapitre sur les fonctions)
function checkPlayerAccess(playerLevel, requiredLevel)
    if playerLevel < requiredLevel then
        print("Access denied - insufficient level")
        return false  -- Sort de la fonction immédiatement
    end
    
    print("Access granted")
    return true
end

-- return peut aussi être utilisé dans le script principal pour l'arrêter
local isServerOnline = false

if not isServerOnline then
    print("Server is offline, stopping script")
    return  -- Arrête l'exécution du script
end

print("This line won't execute if server is offline")
```

## 🎯 5. Exemples Pratiques FiveM

### 5.1 Système de Santé du Joueur

```lua
local function checkPlayerStatus(playerHealth, playerArmor)
    if playerHealth <= 0 then
        print("Player is dead - respawn needed")
        return "dead"
    elseif playerHealth <= 25 then
        print("Critical health - seek medical attention!")
        return "critical"
    elseif playerHealth <= 50 then
        print("Low health - be careful")
        return "low"
    else
        print("Player is in good health")
        return "healthy"
    end
end

-- Test du système
local health = 30
local status = checkPlayerStatus(health, 0)
print("Player status:", status)
```

### 5.2 Système de Niveau et Récompenses

```lua
local playerExperience = 2500

-- Calcul du niveau basé sur l'expérience
local playerLevel = 1
local requiredExp = 100

while playerExperience >= requiredExp do
    playerExperience = playerExperience - requiredExp
    playerLevel = playerLevel + 1
    requiredExp = requiredExp * 1.2  -- Augmentation progressive
    
    print("Level up! New level:", playerLevel)
    
    -- Récompenses spéciales à certains niveaux
    if playerLevel == 10 then
        print("Unlocked: Vehicle spawning")
    elseif playerLevel == 25 then
        print("Unlocked: VIP area access")
    elseif playerLevel == 50 then
        print("Unlocked: Admin privileges")
    end
end

print("Final level:", playerLevel)
print("Remaining experience:", playerExperience)
```

### 5.3 Validation d'Entrée Utilisateur

```lua
local function validatePlayerName(name)
    -- Vérifier si le nom est valide
    if not name or name == "" then
        print("Error: Name cannot be empty")
        return false
    end
    
    -- Vérifier la longueur
    if string.len(name) < 3 then
        print("Error: Name too short (minimum 3 characters)")
        return false
    elseif string.len(name) > 20 then
        print("Error: Name too long (maximum 20 characters)")
        return false
    end
    
    -- Vérifier les caractères autorisés (simulation simple)
    for i = 1, string.len(name) do
        local char = string.sub(name, i, i)
        local ascii = string.byte(char)
        
        -- Lettres majuscules (A-Z) et minuscules (a-z)
        if not ((ascii >= 65 and ascii <= 90) or (ascii >= 97 and ascii <= 122)) then
            print("Error: Only letters are allowed")
            return false
        end
    end
    
    print("Name validation successful")
    return true
end

-- Tests
local testNames = {"", "Al", "Alice", "VeryLongPlayerNameThatExceedsLimit", "Player123", "ValidName"}

for i, name in ipairs(testNames) do
    print("Testing name: '" .. name .. "'")
    validatePlayerName(name)
    print("---")
end
```

## 🎯 6. Exercices Pratiques

### Exercice 1 : Calculateur de Notes (Facile)
```lua
-- Créez un programme qui :
-- 1. Définit une note d'examen (0-100)
-- 2. Utilise if-elseif-else pour attribuer une lettre (A, B, C, D, F)
-- 3. A : 90-100, B : 80-89, C : 70-79, D : 60-69, F : 0-59
-- 4. Affiche la note et la lettre correspondante

-- Votre code ici dans exercices/chapitre2/exercice1.lua:
local examScore = 87  -- Testez avec différentes valeurs

-- À compléter...
```

### Exercice 2 : Table de Multiplication (Facile)
```lua
-- Créez un programme qui :
-- 1. Affiche la table de multiplication d'un nombre (ex: 7)
-- 2. Utilise une boucle for de 1 à 10
-- 3. Format : "7 x 1 = 7"

-- Votre code ici dans exercices/chapitre2/exercice2.lua:
local number = 7

-- À compléter...
```

### Exercice 3 : Jeu de Devinette (Moyen)
```lua
-- Créez un programme qui :
-- 1. Génère un nombre aléatoire entre 1 et 100
-- 2. Simule des tentatives de devinette (utilisez math.random)
-- 3. Guide l'utilisateur avec "trop grand" ou "trop petit"
-- 4. Compte le nombre de tentatives
-- 5. S'arrête quand le nombre est trouvé ou après 10 tentatives

-- Votre code ici dans exercices/chapitre2/exercice3.lua:
local targetNumber = math.random(1, 100)
local maxAttempts = 10

-- À compléter...
```

### Exercice 4 : Validateur de Mot de Passe (Moyen)
```lua
-- Créez une fonction qui valide un mot de passe selon ces critères :
-- 1. Au moins 8 caractères
-- 2. Au moins une majuscule
-- 3. Au moins une minuscule  
-- 4. Au moins un chiffre
-- 5. Affiche des messages d'erreur spécifiques pour chaque critère non respecté

-- Votre code ici dans exercices/chapitre2/exercice4.lua:
local function validatePassword(password)
    -- À compléter...
end

-- Tests
local testPasswords = {"abc", "Password", "password123", "Password123", "Pass123"}
for i, pwd in ipairs(testPasswords) do
    print("Testing: " .. pwd)
    validatePassword(pwd)
    print("---")
end
```

### Exercice 5 : Simulateur de Combat (Difficile)
```lua
-- Créez un système de combat tour par tour :
-- 1. Deux joueurs avec vie (100) et attaque (15-25 aléatoire)
-- 2. Tour par tour jusqu'à ce qu'un joueur meure
-- 3. 10% de chance de coup critique (x2 dégâts)
-- 4. 5% de chance de manquer l'attaque
-- 5. Affichage détaillé de chaque tour

-- Votre code ici dans exercices/chapitre2/exercice5.lua:
local player1Health = 100
local player2Health = 100
local turn = 1

-- À compléter...
```

### Exercice 6 : Système de Niveau Avancé (Difficile)
```lua
-- Créez un système complet :
-- 1. Calcul de niveau basé sur l'expérience (formule exponentielle)
-- 2. Récompenses à des niveaux spécifiques (5, 10, 15, 20, etc.)
-- 3. Calcul du progrès vers le niveau suivant (%)
-- 4. Affichage formaté avec barres de progression simulées

-- Votre code ici dans exercices/chapitre2/exercice6.lua:
local playerExp = 5500

-- À compléter...
```

## 🔧 7. Optimisations et Bonnes Pratiques

### 7.1 Éviter les Boucles Infinies

```lua
-- MAUVAIS : Risque de boucle infinie
local counter = 0
while counter < 10 do
    print(counter)
    -- Oubli d'incrémenter counter !
end

-- BON : Toujours s'assurer que la condition peut devenir fausse
local counter = 0
while counter < 10 do
    print(counter)
    counter = counter + 1  -- Incrémentation cruciale
end
```

### 7.2 Utiliser `break` pour Sortir Tôt

```lua
-- Recherche d'un élément dans une liste
local playerNames = {"Alice", "Bob", "Charlie", "Diana", "Eve"}
local searchName = "Charlie"
local found = false

for i, name in ipairs(playerNames) do
    if name == searchName then
        print("Found " .. searchName .. " at position " .. i)
        found = true
        break  -- Sort dès qu'on trouve l'élément
    end
end

if not found then
    print(searchName .. " not found")
end
```

### 7.3 Préférer `for` à `while` Quand Possible

```lua
-- MOINS BON : while pour un nombre fixe d'itérations
local i = 1
while i <= 10 do
    print(i)
    i = i + 1
end

-- MEILLEUR : for pour un nombre fixe d'itérations
for i = 1, 10 do
    print(i)
end
```

### 7.4 Éviter les Conditions Trop Complexes

```lua
-- DIFFICILE À LIRE :
if (playerLevel >= 20 and playerGold >= 1000 and hasLicense and not isBanned) or (isVip and playerLevel >= 10) or isAdmin then
    -- code...
end

-- PLUS LISIBLE :
local canAccessPremium = (playerLevel >= 20 and playerGold >= 1000 and hasLicense and not isBanned)
local hasVipAccess = (isVip and playerLevel >= 10)

if canAccessPremium or hasVipAccess or isAdmin then
    -- code...
end
```

## 📋 8. Points Clés à Retenir

1. **if-then-end** : Structure conditionnelle de base
2. **elseif** : Pour plusieurs conditions mutuellement exclusives
3. **and, or, not** : Opérateurs logiques pour conditions complexes
4. **Valeurs falsy** : Seules `false` et `nil` sont fausses en Lua
5. **for numérique** : `for i = début, fin, pas do`
6. **for générique** : `for clé, valeur in pairs(table) do`
7. **while** : Boucle avec condition au début
8. **repeat-until** : Boucle avec condition à la fin (au moins une exécution)
9. **break** : Sort d'une boucle immédiatement
10. **return** : Sort d'une fonction ou arrête un script

## ✅ 9. Validation des Acquis

Avant de passer au chapitre suivant, assurez-vous de maîtriser :
- ✅ Les structures if-elseif-else
- ✅ Les opérateurs logiques (and, or, not)
- ✅ Les différents types de boucles (for, while, repeat-until)
- ✅ L'utilisation de break et return
- ✅ La création de conditions complexes lisibles
- ✅ L'optimisation des structures de contrôle
- ✅ La prévention des boucles infinies

## 📖 Prochaine Étape

Dans le **Chapitre 3**, nous découvrirons les **tables**, la structure de données la plus importante et polyvalente de Lua. Les tables sont au cœur de tout script FiveM et permettent de créer des tableaux, des objets, et bien plus encore !

---

*💡 Conseil : Les structures de contrôle sont la base de tout programme complexe. Prenez le temps de maîtriser chaque type avant de continuer !*