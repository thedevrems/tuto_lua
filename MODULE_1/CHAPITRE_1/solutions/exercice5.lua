-- 1. Échange de variables
local valueX = 5
local valueY = 10
print("Before swap - X:", valueX, "Y:", valueY)
valueX, valueY = valueY, valueX
print("After swap - X:", valueX, "Y:", valueY)

-- 2. Rectangle
local rectangleLength = 8
local rectangleWidth = 5
local rectangleArea = rectangleLength * rectangleWidth
local rectanglePerimeter = 2 * (rectangleLength + rectangleWidth)
print("Rectangle - Area:", rectangleArea, "Perimeter:", rectanglePerimeter)

-- 3. Nombre pair (utilise l'opérateur ternaire de Lua)
local testNumber = 8
local isEven = (testNumber % 2 == 0)
print("The number " .. testNumber .. " is " .. (isEven and "even" or "odd"))

-- 4. Phrase complexe
local playerName = "Alice"
local playerAge = 25
local playerCity = "Los Santos"
local playerDescription = "Hello, my name is " .. playerName .. ", I am " .. playerAge .. " years old and I live in " .. playerCity .. "."
print(playerDescription)