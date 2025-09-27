local isAdult = true
local hasLicense = true
local hasVehicle = false

local canDrive = isAdult and hasLicense and hasVehicle

print("Is adult:", isAdult)
print("Has license:", hasLicense)  
print("Has vehicle:", hasVehicle)
print("Can drive:", canDrive)

-- Avec explication
print("Result: The person " .. (canDrive and "can" or "cannot") .. " drive.")