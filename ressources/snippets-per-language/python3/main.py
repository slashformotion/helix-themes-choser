name = "Alice"  # String
age = 25  # Integer
height = 5.5  # Float
is_student = True  # Boolean

# 2. Control Flow (if-else statement)
def check_age(age):
    if age < 18:
        return "You are a minor."
    elif age >= 18 and age < 65:
        return "You are an adult."
    else:
        return "You are a senior citizen."

print(check_age(age))

# 3. Loops (for and while)
print("\nFor loop example:")
for i in range(5):
    print(f"Counting: {i}")

# 4. Functions
def greet(person_name):
    return f"Hello, {person_name}!"
print(greet(name))

# 6. Classes and Objects
class Person:
    def __init__(self, name, age):
        self.name = name
        self.age = age

    def introduce(self):
        return f"Hi, I'm {self.name} and I'm {self.age} years old."

# Creating an object of the class
person1 = Person("Charlie", 28)
print("\nClass and Object Example:")
print(person1.introduce())

# 7. Error Handling (try-except)
print("\nError Handling Example:")
try:
    result = 10 / 0  # This will cause a division by zero error
except ZeroDivisionError:
    print("Oops! You can't divide by zero!")
finally:
    print("This block runs no matter what.")

# 8. List Comprehensions (Compact iteration)
squared_numbers = [x ** 2 for x in range(5)]
print("\nSquared Numbers:", squared_numbers)

# 9. Lambda Functions (Anonymous functions)
add = lambda x, y: x + y
print("\nLambda function result:", add(3, 4))


