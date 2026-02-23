import time

def greet(name: str) -> None:
    """Greet someone nicely."""
    messages = [
        f"Hello, {name}!",
        "How are you today?",
        "Welcome to gif-my-code! 🎨"
    ]
    
    for msg in messages:
        print(msg)
        time.sleep(0.5)

if __name__ == "__main__":
    greet("Developer")
